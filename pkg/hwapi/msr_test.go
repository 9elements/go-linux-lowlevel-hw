package hwapi

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/intel-go/cpuid"
)

const (
	TimeStampCounter = 0x10
	MsrPlatformID    = 0x17
	MsrFsbFreq       = 0x000000cd
	MsrPlatformInfo  = 0x000000ce
	Ia32Efer         = 0xC0000080
)

type msrs struct {
	name string
	msr  int64
	val  uint64
}

func TestReadMSR(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}
	tests := []msrs{
		{"IA32_TIME_STAMP_COUNTER", MsrFsbFreq, 3},
		{"MSR_PLATFORM_INFO", MsrPlatformInfo, 0x80000000},
		{"IA32_PLATFORM_ID", MsrPlatformID, 0},
	}

	for _, test := range tests {
		val, err := h.ReadMSRAllCores(test.msr)
		if err != nil {
			t.Errorf("ReadMSR for MSR %s failed with %v", test.name, err)
		}
		if val != test.val {
			t.Errorf("ReadMSR got unexpected value for MSR %s %v", test.name, val)
		}
	}
}

func TestReadMSRTimeStampCounter(t *testing.T) {
	h := GetAPI()
	if runtime.GOARCH == "amd64" && cpuid.HasFeature(cpuid.TSC) {
		timestamp1, err := h.ReadMSR(TimeStampCounter)
		if err != nil {
			t.Errorf("ReadMSR for MSR IA32_TIMESTAMP_COUNTER failed with %v", err)
		}
		if timestamp1 == 0 {
			t.Errorf("ReadMSR got unexpected value for MSR IA32_TIMESTAMP_COUNTER: %v", timestamp1)
		}
		time.Sleep(time.Millisecond)
		timestamp2, err := h.ReadMSR(TimeStampCounter)
		if err != nil {
			t.Errorf("ReadMSR for MSR IA32_TIMESTAMP_COUNTER failed with %v", err)
		}
		if timestamp2 == 0 {
			t.Errorf("ReadMSR got unexpected value for MSR IA32_TIMESTAMP_COUNTER: %v", timestamp2)
		}
		if timestamp1 == timestamp2 {
			t.Errorf("Timestamps are equal, but shouldn't be")
		}
	}
}

func TestReadMSREFER(t *testing.T) {
	h := GetAPI()
	if runtime.GOARCH == "amd64" {
		ia32efer, err := h.ReadMSRAllCores(Ia32Efer)
		if err != nil {
			t.Errorf("ReadMSRAllCores for MSR IA32_EFER failed with %v", err)
		} else {
			if ia32efer == 0 {
				t.Errorf("ReadMSRAllCores got unexpected value for MSR IA32_EFER: %v", ia32efer)
			}
		}
	}
}

func TestSMRR(t *testing.T) {

	h := GetAPI()

	has, err := h.HasSMRR()
	if err != nil {
		t.Errorf("HasSMRR() failed: %v", err)
	} else if has {
		t.Log("System has SMRR")

		got, err := h.GetSMRRInfo()

		if err != nil {
			t.Errorf("GetSMRRInfo() failed: %v", err)
		}
		active := true
		if got.PhysBase == 0 || got.PhysBase == 0xfffffff || got.PhysBase == 0xffffffff {
			active = false
		}
		if got.PhysMask == 0 || got.PhysMask == 0xf || got.PhysMask == 0xffffffff {
			active = false
		}
		if got.Active != active {
			t.Error("Invalid SMRR config.")
		}

		if got.Active {
			t.Logf("SMRR is active. PHYS_BASE: %x, PHYS_MASK: %x", got.PhysBase, got.PhysMask)
		} else {
			t.Log("SMRR is not active")
		}
	} else {
		t.Skip("Hardware has no SMRR support")
	}
}
