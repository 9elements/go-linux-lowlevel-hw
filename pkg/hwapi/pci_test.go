package hwapi

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"
)

func TestPCIQemu(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}

	d0f0 := PCIDevice{
		Bus:      0,
		Device:   0,
		Function: 0,
	}

	reg16, err := h.PCIReadConfigSpace(d0f0, 0, 2)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}

	if binary.LittleEndian.Uint16(reg16) != 0x8086 {
		t.Errorf("Unexpected value: %v", reg16)
	}
	reg16, err = h.PCIReadConfigSpace(d0f0, 2, 2)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}

	if binary.LittleEndian.Uint16(reg16) != 0x29c0 {
		t.Errorf("Unexpected value: %v", reg16)
	}

	d1ff0 := PCIDevice{
		Bus:      0,
		Device:   0x1f,
		Function: 0,
	}

	reg16, err = h.PCIReadConfigSpace(d1ff0, 0, 2)
	if err != nil {
		t.Errorf("PCIReadVendorID failed with error %v", err)
	}
	if binary.LittleEndian.Uint16(reg16) != 0x8086 {
		t.Errorf("Unexpected value: %v", reg16)
	}
	reg16, err = h.PCIReadConfigSpace(d1ff0, 2, 2)
	if err != nil {
		t.Errorf("PCIReadDeviceID failed with error %v", err)
	}
	if binary.LittleEndian.Uint16(reg16) != 0x2918 {
		t.Errorf("Unexpected value: %v", reg16)
	}

	var class uint16

	reg8, err := h.PCIReadConfigSpace(d0f0, 0xc, 1)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if binary.LittleEndian.Uint16(reg8) != 0 {
		t.Errorf("Unexpected value: %v", reg8)
	}
	class |= binary.LittleEndian.Uint16(reg8) << 8

	reg8, err = h.PCIReadConfigSpace(d0f0, 0xb, 1)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if byte(reg8[0]) != 6 {
		t.Errorf("Unexpected value: %v", reg8)
	}
	class |= binary.LittleEndian.Uint16(reg8)

	reg16, err = h.PCIReadConfigSpace(d0f0, 0xb, 16)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if binary.LittleEndian.Uint16(reg16) != class {
		t.Errorf("Unexpected value: %v", reg16)
	}

	d1f0 := PCIDevice{
		Bus:      0,
		Device:   0x1,
		Function: 0,
	}

	reg32, err := h.PCIReadConfigSpace(d1f0, 0x10, 4)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if binary.LittleEndian.Uint32(reg32) == 0 {
		t.Errorf("Unexpected value: %x", reg32)
	}

	reg32, err = h.PCIReadConfigSpace(d1f0, 0x18, 4)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if binary.LittleEndian.Uint32(reg32) == 0 {
		t.Errorf("Unexpected value: %x", reg32)
	}

	backup, err := h.PCIReadConfigSpace(d1f0, 0x10, 4)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if binary.LittleEndian.Uint32(reg32) == 0 || binary.LittleEndian.Uint32(reg32) == 0xffffffff {
		t.Errorf("Unexpected value: %x", reg32)
	}
	reg32 = []byte{0xff, 0xff, 0xff, 0xff}
	if err := h.PCIWriteConfigSpace(d1f0, 0x10, reg32); err != nil {
		t.Errorf("PCIWriteConfig32 failed with error %v", err)
	}

	// check if bits are moving
	reg32, err = h.PCIReadConfigSpace(d1f0, 0x10, 4)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if binary.LittleEndian.Uint32(reg32) == 0xfd000008 {
		t.Errorf("Unexpected value: %x", reg32)
	}

	if err := h.PCIWriteConfigSpace(d1f0, 0x10, backup); err != nil {
		t.Errorf("PCIWriteConfig32 failed with error %v", err)
	}
}

func TestPCIDeviceVendorIDQemu(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}

	d0f0 := PCIDevice{
		Bus:      0,
		Device:   0,
		Function: 0,
	}

	reg16, err := h.PCIReadConfigSpace(d0f0, 0, 2)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if binary.LittleEndian.Uint16(reg16) != 0x8086 {
		t.Errorf("Unexpected value: %v", reg16)
	}
	reg16, err = h.PCIReadConfigSpace(d0f0, 2, 2)
	if err != nil {
		t.Errorf("PCIReadConfig failed with error %v", err)
	}
	if binary.LittleEndian.Uint16(reg16) != 0x29c0 {
		t.Errorf("Unexpected value: %v", reg16)
	}

	d1ff0 := PCIDevice{
		Bus:      0,
		Device:   0x1f,
		Function: 0,
	}
	reg16, err = h.PCIReadConfigSpace(d1ff0, 0, 2)
	if err != nil {
		t.Errorf("PCIReadVendorID failed with error %v", err)
	}
	if binary.LittleEndian.Uint16(reg16) != 0x8086 {
		t.Errorf("Unexpected value: %v", reg16)
	}
	reg16, err = h.PCIReadConfigSpace(d1ff0, 2, 2)
	if err != nil {
		t.Errorf("PCIReadDeviceID failed with error %v", err)
	}
	if binary.LittleEndian.Uint16(reg16) != 0x2918 {
		t.Errorf("Unexpected value: %v", reg16)
	}
}

func TestPCIBusMasterQEMU(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}

	d1f0 := PCIDevice{
		Bus:      0,
		Device:   0x1,
		Function: 0,
	}

	backup, err := h.PCIReadConfigSpace(d1f0, 4, 1)
	if err != nil {
		t.Errorf("PCIReadConfig8 failed with error %v", err)
	}

	reg8 := []byte{byte(backup[0]) ^ 4}
	if err = h.PCIWriteConfigSpace(d1f0, 4, &reg8); err != nil {
		if err := h.PCIWriteConfigSpace(d1f0, 4, backup); err != nil {
			t.Error(err)
		}
		t.Errorf("PCIWriteConfig8 failed with error %v", err)
	}

	reg8, err = h.PCIReadConfigSpace(d1f0, 4, 1)
	if err != nil {
		if err := h.PCIWriteConfigSpace(d1f0, 4, backup); err != nil {
			t.Error(err)
		}
		t.Errorf("PCIReadConfig8 failed with error %v", err)
	}
	if byte(reg8[0]) == byte(backup[0]) {
		if err := h.PCIWriteConfigSpace(d1f0, 4, backup); err != nil {
			t.Error(err)
		}
		t.Errorf("PCIWriteConfig8 failed. Register content is unchanged.")
	}
}

func TestPCIBusMaster2QEMU(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}

	d1f0 := PCIDevice{
		Bus:      0,
		Device:   0x1,
		Function: 0,
	}

	backup, err := h.PCIReadConfigSpace(d1f0, 4, 2)
	if err != nil {
		t.Errorf("PCIReadConfig16 failed with error %v", err)
	}

	err = h.PCIWriteConfigSpace(d1f0, 5, 0xff)
	if err != nil {
		t.Errorf("PCIWriteConfig8 failed with error %v", err)
	}

	reg16, err := h.PCIReadConfigSpace(d1f0, 4, 2)
	if err != nil {
		t.Errorf("PCIReadConfig16 failed with error %v", err)
	}
	if !bytes.Equal(reg16, backup) {
		t.Errorf("PCIWriteConfig8 failed. Register content is unchanged.")
	}
	// restore register
	err = h.PCIWriteConfigSpace(d1f0, 4, backup)
	if err != nil {
		t.Errorf("PCIWriteConfig16 failed with error %v", err)
	}
	reg16, err = h.PCIReadConfigSpace(d1f0, 4, 16)
	if err != nil {
		t.Errorf("PCIReadConfig16 failed with error %v", err)
	}
	if !bytes.Equal(reg16, backup) {
		t.Errorf("PCIWriteConfig16 failed. Register content is unchanged.")
	}
}

func TestPCIEnum2QEMU(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}

	l := []PCIDevice{}

	err := h.PCIEnumerateVisibleDevices(func(d PCIDevice) (abort bool) {
		l = append(l, d)
		return false
	})
	if err != nil {
		t.Errorf("PCIEnumerateVisibleDevices failed with error %v", err)
	}
	reference := []PCIDevice{
		{Bus: 0, Device: 0, Function: 0},
		{Bus: 0, Device: 1, Function: 0},
		{Bus: 0, Device: 2, Function: 0},
		{Bus: 0, Device: 0x1f, Function: 0},
		{Bus: 0, Device: 0x1f, Function: 2},
		{Bus: 0, Device: 0x1f, Function: 3},
	}

	for _, r := range reference {
		found := false
		for i := range l {
			if l[i].Device == r.Device &&
				l[i].Function == r.Function &&
				l[i].Bus == r.Bus {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("PCI device %v is missing", r)
		}
	}
}
