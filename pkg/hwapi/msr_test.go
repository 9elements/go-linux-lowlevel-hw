package hwapi

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/micgor32/go-msr"
	"github.com/u-root/cpuid"
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
}

func TestReadMSR(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}
	err := msr.MSR(0, func(dev msr.MSRDev) error {
		// This closure is dummy, and left here only to
		// satisfy msr.MSR function definition. We can
		// safelu return nil here, as the function won't
		// ever reach the point of returning the clousure
		// if the caller permissions are not sufficient to
		// read /dev/cpu/n/msr
		return nil
	})
	if err != nil {
		t.Skip("Not enough permissions to do test")
	}

	tests := []msrs{
		{"MSR_FSB_FREQ", MsrFsbFreq},
		{"MSR_PLATFORM_INFO", MsrPlatformInfo},
		{"IA32_PLATFORM_ID", MsrPlatformID},
	}
	for _, test := range tests {
		val := h.ReadMSR(test.msr)

		if val == 0xff {
			t.Errorf("MSR value for %s is not being read correctly.", test.name)
		}

		// We can't check the consistency since ReadMSR reads only from core 0,
		// but we can check whether SOME value other than 0xff is being read.
		// Left for previous test for reference.
		// for iterator, value := range vals {
		// 	if iterator < len(vals) && value != vals[iterator+1] {
		// 		if value != vals[iterator+1] {
		// 			t.Errorf("MSR value are not the same for all cores. Core: %d, Value: 0x%x, Previous value: 0x%x", iterator, value, vals[iterator+1])
		// 		}
		// 		if value == 0 || vals[iterator+1] == 0xffffffffffffffff {
		// 			t.Errorf("ReadMSR got unexpected value for MSR %s %v", test.name, vals)
		// 		}
		// 	}
		// }
	}

}

func TestReadMSRTimeStampCounter(t *testing.T) {
	h := GetAPI()
	err := msr.MSR(0, func(dev msr.MSRDev) error {
		return nil
	})

	if err != nil {
		t.Skip("Not enough permissions to do test")
	}
	if runtime.GOARCH == "amd64" && cpuid.HasFeature(cpuid.TSC) {
		timestamp1 := h.ReadMSR(TimeStampCounter)

		if timestamp1 == 0 {
			t.Errorf("ReadMSR got unexpected value for MSR IA32_TIMESTAMP_COUNTER: %v", timestamp1)
		}
		time.Sleep(time.Millisecond)
		timestamp2 := h.ReadMSR(TimeStampCounter)
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

	err := msr.MSR(0, func(dev msr.MSRDev) error {
		return nil
	})
	if err != nil {
		t.Skip("Not enough permissions to do test")
	}
	if runtime.GOARCH == "amd64" {
		val := h.ReadMSR(Ia32Efer)
		if val == 0xff {
			t.Error("MSR value is not being read correctly.")
		}

		// Left for previous test for reference.
		// for iterator, value := range vals {
		// 	if iterator < len(vals) && value != vals[iterator+1] {
		// 		t.Errorf("MSR value are not the same for all cores. Core: %d, Value: 0x%x, Next value: 0x%x", iterator, value, vals[iterator+1])
		// 	}
		// }
	}
}

func TestSMRR(t *testing.T) {
	h := GetAPI()

	err := msr.MSR(0, func(dev msr.MSRDev) error {
		return nil
	})
	if err != nil {
		t.Skip("Not enough permissions to do test")
	}

	has, err := HasSMRR(h)
	if err != nil {
		t.Errorf("HasSMRR() failed: %v", err)
	} else if has {
		t.Log("System has SMRR")

		got, err := GetSMRRInfo(h)

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
