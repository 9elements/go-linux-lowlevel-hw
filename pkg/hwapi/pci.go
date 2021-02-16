package hwapi

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

//PCIReadConfigSpace reads from PCI config space into buf
func (h HwApi) PCIReadConfigSpace(bus int, device int, devFn int, off int, buf interface{}) (err error) {
	var path string
	var f *os.File
	path = fmt.Sprintf("/sys/bus/pci/devices/0000:%02x:%02x.%1x/config", bus, device, devFn)

	f, err = os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Seek(int64(off), io.SeekStart)
	if err != nil {
		return
	}
	err = binary.Read(f, binary.LittleEndian, buf)

	return
}

//PCIReadConfig8 reads 8bits from PCI config space
func (h HwApi) PCIReadConfig8(bus int, device int, devFn int, off int) (reg8 uint8, err error) {

	err = h.PCIReadConfigSpace(bus, device, devFn, off, &reg8)

	return
}

//PCIReadConfig16 reads 16bits from PCI config space
func (h HwApi) PCIReadConfig16(bus int, device int, devFn int, off int) (reg16 uint16, err error) {

	err = h.PCIReadConfigSpace(bus, device, devFn, off, &reg16)

	return
}

//PCIReadConfig32 reads 32bits from PCI config space
func (h HwApi) PCIReadConfig32(bus int, device int, devFn int, off int) (reg32 uint32, err error) {

	err = h.PCIReadConfigSpace(bus, device, devFn, off, &reg32)

	return
}

//PCIReadVendorID reads the device vendor ID from PCI config space
func (h HwApi) PCIReadVendorID(bus int, device int, devFn int) (id uint16, err error) {
	id, err = h.PCIReadConfig16(bus, device, devFn, 0)

	return
}

//PCIReadDeviceID reads the device ID from PCI config space
func (h HwApi) PCIReadDeviceID(bus int, device int, devFn int) (id uint16, err error) {
	id, err = h.PCIReadConfig16(bus, device, devFn, 2)

	return
}

//PCIWriteConfigSpace writes to PCI config space from buf
func (h HwApi) PCIWriteConfigSpace(bus int, device int, devFn int, off int, buf interface{}) (err error) {
	var path string
	var f *os.File
	path = fmt.Sprintf("/sys/bus/pci/devices/0000:%02x:%02x.%1x/config", bus, device, devFn)

	f, err = os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Seek(int64(off), io.SeekStart)
	if err != nil {
		return
	}
	err = binary.Write(f, binary.LittleEndian, buf)

	return
}

//PCIWriteConfig8 writes 8bits to PCI config space
func (h HwApi) PCIWriteConfig8(bus int, device int, devFn int, off int, val uint8) (err error) {

	err = h.PCIWriteConfigSpace(bus, device, devFn, off, val)

	return
}

//PCIWriteConfig16 writes 16bits to PCI config space
func (h HwApi) PCIWriteConfig16(bus int, device int, devFn int, off int, val uint16) (err error) {

	err = h.PCIWriteConfigSpace(bus, device, devFn, off, val)

	return
}

//PCIWriteConfig32 writes 32bits to PCI config space
func (h HwApi) PCIWriteConfig32(bus int, device int, devFn int, off int, val uint32) (err error) {

	err = h.PCIWriteConfigSpace(bus, device, devFn, off, val)

	return
}
