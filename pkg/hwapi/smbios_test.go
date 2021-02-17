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
		if s.Vendor != "SeaBIOS" {
			t.Errorf("Got unexpected vendor %s", s.Vendor)
		}
		if s.SystemBiosMajor != 0 {
			t.Errorf("Got unexpected bios version %d", s.SystemBiosMajor)
		}
		if s.BIOSStartingAddress != 0xe8000 {
			t.Errorf("Got unexpected starting address %x", s.BIOSStartingAddress)
		}
		if s.BIOSSize != 0x10000 {
			t.Errorf("Got unexpected bios size %x", s.BIOSSize)
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
