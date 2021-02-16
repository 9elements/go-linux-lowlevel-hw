package hwapi

import (
	"os"
	"testing"
)

func TestSMBIOSQemu(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}
	count := 0
	_, err := h.IterateOverSMBIOSTablesType0(func(s *SMBIOSType0) bool {
		count++
		if s.Type != 0 {
			t.Errorf("Got unexpected type %d", s.Type)
		}
		if s.Vendor != "Emulation" {
			t.Errorf("Got unexpected vendor %s", s.Vendor)
			t.Errorf("%+v", s)

		}
		return false
	})
	if err != nil {
		t.Errorf("IterateOverSMBIOSTables failed with %v", err)
	}
	if count == 0 {
		t.Errorf("Type 0 not found")
	}
}
