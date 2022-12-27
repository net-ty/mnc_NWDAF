package service

import (
	"context"
	"errors"
	"net"

	"github.com/sirupsen/logrus"
	gtpv1 "github.com/wmnsk/go-gtp/v1"

	n3iwf_context "github.com/free5gc/n3iwf/context"
	"github.com/free5gc/n3iwf/logger"
)

var gtpLog *logrus.Entry

func init() {
	gtpLog = logger.GTPLog
}

// SetupGTPTunnelWithUPF set up GTP connection with UPF
// return *gtpv1.UPlaneConn, net.Addr and error
func SetupGTPTunnelWithUPF(upfIPAddr string) (*gtpv1.UPlaneConn, net.Addr, error) {
	n3iwfSelf := n3iwf_context.N3IWFSelf()

	// Set up GTP connection
	upfUDPAddr := upfIPAddr + ":2152"

	remoteUDPAddr, err := net.ResolveUDPAddr("udp", upfUDPAddr)
	if err != nil {
		gtpLog.Errorf("Resolve UDP address %s failed: %+v", upfUDPAddr, err)
		return nil, nil, errors.New("Resolve Address Failed")
	}

	n3iwfUDPAddr := n3iwfSelf.GTPBindAddress + ":2152"

	localUDPAddr, err := net.ResolveUDPAddr("udp", n3iwfUDPAddr)
	if err != nil {
		gtpLog.Errorf("Resolve UDP address %s failed: %+v", n3iwfUDPAddr, err)
		return nil, nil, errors.New("Resolve Address Failed")
	}

	context := context.TODO()

	// Dial to UPF
	userPlaneConnection, err := gtpv1.DialUPlane(context, localUDPAddr, remoteUDPAddr)
	if err != nil {
		gtpLog.Errorf("Dial to UPF failed: %+v", err)
		return nil, nil, errors.New("Dial failed")
	}

	return userPlaneConnection, remoteUDPAddr, nil
}

// ListenAndServe binds and listens user plane socket on N3IWF N3 interface,
// catching GTP packets and send it to NWu interface
func ListenAndServe(userPlaneConnection *gtpv1.UPlaneConn) error {
	go listenGTP(userPlaneConnection)
	return nil
}

// listenGTP handle the gtpv1 UPlane connection. It reads packets(without
// GTP header) from the connection and call forward() to forward user data
// to NWu interface.
func listenGTP(userPlaneConnection *gtpv1.UPlaneConn) {
	defer func() {
		err := userPlaneConnection.Close()
		if err != nil {
			gtpLog.Errorf("userPlaneConnection Close failed: %+v", err)
		}
	}()

	payload := make([]byte, 65535)

	for {
		n, _, teid, err := userPlaneConnection.ReadFromGTP(payload)
		gtpLog.Tracef("Read %d bytes", n)
		if err != nil {
			gtpLog.Errorf("Read from GTP failed: %+v", err)
			return
		}

		forwardData := make([]byte, n)
		copy(forwardData, payload[:n])

		go forward(teid, forwardData)
	}
}

// forward forwards user plane packets from N3 to UE,
// with GRE header and new IP header encapsulated
func forward(ueTEID uint32, packet []byte) {
	// N3IWF context
	self := n3iwf_context.N3IWFSelf()
	// IPv4 packet connection
	ipv4PacketConn := self.NWuIPv4PacketConn
	// Find UE information
	ue, ok := self.AllocatedUETEIDLoad(ueTEID)
	if !ok {
		gtpLog.Error("UE context not found")
		return
	}
	// UE IP
	ueInnerIPAddr := ue.IPSecInnerIPAddr

	// GRE header
	greHeader := []byte{0, 0, 8, 0}
	// IP payload
	greEncapsulatedPacket := append(greHeader, packet...)

	// Send to UE
	if n, err := ipv4PacketConn.WriteTo(greEncapsulatedPacket, nil, ueInnerIPAddr); err != nil {
		gtpLog.Errorf("Write to UE failed: %+v", err)
		return
	} else {
		gtpLog.Trace("Forward NWu <- N3")
		gtpLog.Tracef("Wrote %d bytes", n)
	}
}
