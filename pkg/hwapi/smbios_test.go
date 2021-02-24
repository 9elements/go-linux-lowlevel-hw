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

func TestSMBIOSType17(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}
	count := 0
	_, err := h.IterateOverSMBIOSTablesType17(func(s *SMBIOSType17) bool {
		count++
		if s.Type != 17 {
			t.Errorf("Got unexpected type %d", s.Type)
		}
		if s.TotalWidth < s.DataWidth {
			t.Errorf("TotalWidth (%x) < DataWidth (%x\n", s.TotalWidth, s.DataWidth)
		}
		if s.DataWidth == 0 {
			t.Errorf("DataWidth is zero\n")
		}
		if s.Size != 0 && s.Size != 0xFFFF && s.Size != 0x7FFF {
			if s.Size&0xFFFF > 0 {
				t.Errorf("Size is not multiple of 1 KiB\n")
			}
		}
		t.Logf("TotalWidth 0x%x\n", s.TotalWidth)
		t.Logf("DataWidth 0x%x\n", s.DataWidth)
		t.Logf("DeviceLocator %s\n", s.DeviceLocator)
		t.Logf("BankLocator %s\n", s.BankLocator)
		t.Logf("Size 0x%x\n", s.Size)

		return false
	})
	if err != nil {
		t.Errorf("IterateOverSMBIOSTables failed with %v", err)
	}
	if count == 0 {
		t.Errorf("Type 17 not found")
	}
}
