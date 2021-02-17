package hwapi

import (
	"os"
	"testing"
)

func TestPCIQemu(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}

	reg16, err := h.PCIReadConfig16(0, 0, 0, 0)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if reg16 != 0x8086 {
		t.Errorf("Unexpected value: %v", reg16)
	}
	reg16, err = h.PCIReadConfig16(0, 0, 0, 2)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if reg16 != 0x29c0 {
		t.Errorf("Unexpected value: %v", reg16)
	}

	reg16, err = h.PCIReadVendorID(0, 0x1f, 0)
	if err != nil {
		t.Errorf("PCIReadVendorID failed with error %v", err)
	}
	if reg16 != 0x8086 {
		t.Errorf("Unexpected value: %v", reg16)
	}
	reg16, err = h.PCIReadDeviceID(0, 0x1f, 0)
	if err != nil {
		t.Errorf("PCIReadDeviceID failed with error %v", err)
	}
	if reg16 != 0x2918 {
		t.Errorf("Unexpected value: %v", reg16)
	}

	var class uint16

	reg8, err := h.PCIReadConfig8(0, 0, 0, 0xc)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if reg8 != 0 {
		t.Errorf("Unexpected value: %v", reg8)
	}
	class |= uint16(reg8) << 8

	reg8, err = h.PCIReadConfig8(0, 0, 0, 0xb)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if reg8 != 6 {
		t.Errorf("Unexpected value: %v", reg8)
	}
	class |= uint16(reg8)

	reg16, err = h.PCIReadConfig16(0, 0, 0, 0xb)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if reg16 != class {
		t.Errorf("Unexpected value: %v", reg16)
	}

	reg32, err := h.PCIReadConfig32(0, 1, 0, 0x10)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if reg32 != 0xfd000008 {
		t.Errorf("Unexpected value: %v", reg32)
	}

	reg32, err = h.PCIReadConfig32(0, 1, 0, 0x18)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if reg32 != 0xfebd4000 {
		t.Errorf("Unexpected value: %v", reg32)
	}

	backup, err := h.PCIReadConfig32(0, 1, 0, 0x10)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if reg32 != 0xfebd4000 {
		t.Errorf("Unexpected value: %x", reg32)
	}

	reg32 = 0xffffffff
	err = h.PCIWriteConfig32(0, 1, 0, 0x10, reg32)
	if err != nil {
		t.Errorf("PCIWriteConfig32 failed with error %v", err)
	}

	// check if bits are moving
	reg32, err = h.PCIReadConfig32(0, 1, 0, 0x10)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if reg32 == 0xfd000008 {
		t.Errorf("Unexpected value: %x", reg32)
	}

	err = h.PCIWriteConfig32(0, 1, 0, 0x10, backup)
	if err != nil {
		t.Errorf("PCIWriteConfig32 failed with error %v", err)
	}

}
