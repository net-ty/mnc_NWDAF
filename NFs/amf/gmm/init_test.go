package gmm_test

import (
	"fmt"
	"testing"

	"github.com/free5gc/amf/gmm"
	"github.com/free5gc/fsm"
)

func TestGmmFSM(t *testing.T) {
	if err := fsm.ExportDot(gmm.GmmFSM, "gmm"); err != nil {
		fmt.Printf("fsm export data return error: %+v", err)
	}
}
