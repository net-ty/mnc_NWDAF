package context

import (
	"errors"
	"net"
)

type IPAllocator struct {
	ipNetwork *net.IPNet
	g         *_IDPool
}

func NewIPAllocator(cidr string) (*IPAllocator, error) {
	allocator := &IPAllocator{}

	if _, ipnet, err := net.ParseCIDR(cidr); err != nil {
		return nil, err
	} else {
		allocator.ipNetwork = ipnet
	}
	allocator.g = newIDPool(1, 1<<int64(32-maskBits(allocator.ipNetwork.Mask))-2)

	return allocator, nil
}

func maskBits(mask net.IPMask) int {
	var cnt int
	for _, b := range mask {
		for ; b != 0; b /= 2 {
			if b%2 != 0 {
				cnt++
			}
		}
	}
	return cnt
}

// IPAddrWithOffset add offset on base ip
func IPAddrWithOffset(ip net.IP, offset int) net.IP {
	retIP := make(net.IP, len(ip))
	copy(retIP, ip)

	var carry int
	for i := len(retIP) - 1; i >= 0; i-- {
		if offset == 0 {
			break
		}

		val := int(retIP[i]) + carry + offset%256
		retIP[i] = byte(val % 256)
		carry = val / 256

		offset /= 256
	}

	return retIP
}

// IPAddrOffset calculate the input ip with base ip offset
func IPAddrOffset(in, base net.IP) int {
	offset := 0
	exp := 1
	for i := len(base) - 1; i >= 0; i-- {
		offset += int(in[i]-base[i]) * exp
		exp *= 256
	}
	return offset
}

// Allocate will allocate the IP address and returns it
func (a *IPAllocator) Allocate() (net.IP, error) {
	if offset, err := a.g.allocate(); err != nil {
		return nil, errors.New("ip allocation failed" + err.Error())
	} else {
		return IPAddrWithOffset(a.ipNetwork.IP, int(offset)), nil
	}
}

func (a *IPAllocator) Release(ip net.IP) {
	offset := IPAddrOffset(ip, a.ipNetwork.IP)
	a.g.release(int64(offset))
}

type _IDPool struct {
	minValue int64
	maxValue int64
	isUsed   map[int64]bool
}

func newIDPool(minValue int64, maxValue int64) (idPool *_IDPool) {
	idPool = new(_IDPool)
	idPool.minValue = minValue
	idPool.maxValue = maxValue
	idPool.isUsed = make(map[int64]bool)
	return
}

func (i *_IDPool) allocate() (id int64, err error) {
	for id := i.minValue; id <= i.maxValue; id++ {
		if _, exist := i.isUsed[id]; !exist {
			i.isUsed[id] = true
			return id, nil
		}
	}
	return 0, errors.New("No available value range to allocate id")
}

func (i *_IDPool) release(id int64) {
	delete(i.isUsed, id)
}
