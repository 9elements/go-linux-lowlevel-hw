package hwapi

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// PCIDevice represents a PCI device
type PCIDevice struct {
	Bus      int
	Device   int
	Function int
	// True if device is hidden
	Hidden bool
	// BARs currently decoded by the device
	BAR map[int]uint64
	// ROM BAR currently decoded by the device
	ROM uint64
}

//pciReadConfigSpace reads from PCI config space into buf
func (h HwApi) pciReadConfigSpace(d PCIDevice, off int, buf interface{}) (err error) {
	var path string
	var f *os.File
	path = fmt.Sprintf("/sys/bus/pci/devices/0000:%02x:%02x.%1x/config", d.Bus, d.Device, d.Function)

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
func (h HwApi) PCIReadConfig8(d PCIDevice, off int) (reg8 uint8, err error) {

	err = h.pciReadConfigSpace(d, off, &reg8)

	return
}

//PCIReadConfig16 reads 16bits from PCI config space
func (h HwApi) PCIReadConfig16(d PCIDevice, off int) (reg16 uint16, err error) {

	err = h.pciReadConfigSpace(d, off, &reg16)

	return
}

//PCIReadConfig32 reads 32bits from PCI config space
func (h HwApi) PCIReadConfig32(d PCIDevice, off int) (reg32 uint32, err error) {

	err = h.pciReadConfigSpace(d, off, &reg32)

	return
}

//PCIReadVendorID reads the device vendor ID from PCI config space
func (h HwApi) PCIReadVendorID(d PCIDevice) (id uint16, err error) {
	id, err = h.PCIReadConfig16(d, 0)

	return
}

//PCIReadDeviceID reads the device ID from PCI config space
func (h HwApi) PCIReadDeviceID(d PCIDevice) (id uint16, err error) {
	id, err = h.PCIReadConfig16(d, 2)

	return
}

//pciWriteConfigSpace writes to PCI config space from buf
func (h HwApi) pciWriteConfigSpace(d PCIDevice, off int, buf interface{}) (err error) {
	var path string
	var f *os.File
	path = fmt.Sprintf("/sys/bus/pci/devices/0000:%02x:%02x.%1x/config", d.Bus, d.Device, d.Function)

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
func (h HwApi) PCIWriteConfig8(d PCIDevice, off int, val uint8) (err error) {

	err = h.pciWriteConfigSpace(d, off, val)

	return
}

//PCIWriteConfig16 writes 16bits to PCI config space
func (h HwApi) PCIWriteConfig16(d PCIDevice, off int, val uint16) (err error) {

	err = h.pciWriteConfigSpace(d, off, val)

	return
}

//PCIWriteConfig32 writes 32bits to PCI config space
func (h HwApi) PCIWriteConfig32(d PCIDevice, off int, val uint32) (err error) {

	err = h.pciWriteConfigSpace(d, off, val)

	return
}
