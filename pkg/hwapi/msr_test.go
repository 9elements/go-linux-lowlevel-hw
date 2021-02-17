package hwapi

import (
	"os"
	"testing"
)

const (
	MsrFsbFreq      = 0x000000cd
	MsrPlatformInfo = 0x000000ce
)

type msrs struct {
	msr int64
	val uint64
}

func TestReadMSR(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}
	tests := []msrs{
		{MsrFsbFreq, 3},
		{MsrPlatformInfo, 0x80000000},
	}

	for _, test := range tests {
		val, err := h.ReadMSR(test.msr)
		if err != nil {
			t.Errorf("ReadMSR failed with %v", err)
		}
		if val != test.val {
			t.Errorf("MsrFsbFreq unexpected value %v", val)
		}
	}
}

func TestSMRR(t *testing.T) {
	t.Skip()

	h := GetAPI()

	has, err := h.HasSMRR()
	if err != nil {
		t.Errorf("HasSMRR() failed: %v", err)
	}

	if has {
		t.Log("System has SMRR")

		got, err := h.GetSMRRInfo()

		if err != nil {
			t.Errorf("GetSMRRInfo() failed: %v", err)
		}

		if got.Active != (got.PhysBase != 0 && got.PhysMask != 0) {
			t.Error("Invalid SMRR config.")
		}

		if got.Active {
			t.Logf("SMRR is active. PHYS_BASE: %x, PHYS_MASK: %x", got.PhysBase, got.PhysMask)
		} else {
			t.Log("SMRR is not active")
		}
	} else {
		t.Log("No SMRR")
	}
}
