package hwapi

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

// PCIEnumerateVisibleDevices enumerates all visible PCI devices
func (h HwAPI) PCIEnumerateVisibleDevices(cb func(d PCIDevice) (abort bool)) (err error) {
	dir := "/sys/bus/pci/devices/"
	err = filepath.Walk(dir, func(path string, info os.FileInfo, interr error) error {
		if interr != nil || path == dir {
			return nil
		}
		if strings.HasPrefix(path, dir) {
			path = strings.Replace(path, dir, "", 1)
		}
		if strings.HasPrefix(info.Mode().String(), "L") {
			domain := 0
			bus := 0
			device := 0
			function := 0
			_, err = fmt.Sscanf(path, "%4x:%2x:%2x.%1x", &domain, &bus, &device, &function)
			if err != nil {
				return err
			}
			d := PCIDevice{Bus: bus,
				Device:   device,
				Function: function}

			if cb(d) {
				return filepath.SkipDir
			}
		}
		return nil
	})
	return
}

//pciReadConfigSpace reads from PCI config space into out
func (h HwAPI) PCIReadConfigSpace(d PCIDevice, off int, out interface{}) (err error) {
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
	err = binary.Read(f, binary.LittleEndian, out)

	return
}

//PCIReadVendorID reads the device vendor ID from PCI config space
func (h HwAPI) PCIReadVendorID(d PCIDevice) (id uint16, err error) {
	if err = h.PCIReadConfigSpace(d, 0, &id); err != nil {
		return 0, err
	}

	return
}

//PCIReadDeviceID reads the device ID from PCI config space
func (h HwAPI) PCIReadDeviceID(d PCIDevice) (id uint16, err error) {
	if err = h.PCIReadConfigSpace(d, 2, &id); err != nil {
		return 0, err
	}

	return
}

//pciWriteConfigSpace writes to PCI config space from in
func (h HwAPI) PCIWriteConfigSpace(d PCIDevice, off int, in interface{}) (err error) {
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
	err = binary.Write(f, binary.LittleEndian, in)

	return
}
