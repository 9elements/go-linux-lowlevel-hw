package hwapi

import "github.com/digitalocean/go-smbios/smbios"

// LowLevelHardwareInterfaces provides methods to access hardware found on modern x86_64 platforms
type LowLevelHardwareInterfaces interface {
	// cpuid.go
	VersionString() string
	HasSMX() bool
	HasVMX() bool
	HasMTRR() bool
	ProcessorBrandName() string
	CPUSignature() uint32
	CPUSignatureFull() (uint32, uint32, uint32, uint32)
	CPULogCount() uint32
	CPUID(uint32, uint32) (uint32, uint32, uint32, uint32)

	// e820.go
	IterateOverE820Ranges(target string, callback func(start uint64, end uint64) bool) (bool, error)

	// iommu.go
	LookupIOAddress(addr uint64, regs VTdRegisters) ([]uint64, error)

	// msr.go
	ReadMSR(msr int64) []uint64

	// pci.go
	PCIEnumerateVisibleDevices(cb func(d PCIDevice) (abort bool)) (err error)
	PCIReadConfigSpace(d PCIDevice, off int, len int) ([]byte, error)
	PCIWriteConfigSpace(d PCIDevice, off int, val interface{}) error

	// phys.go
	ReadPhys(addr int64, data UintN) error
	ReadPhysBuf(addr int64, buf []byte) error
	WritePhys(addr int64, data UintN) error

	// tpm.go
	NewTPM() (*TPM, error)
	NVLocked(tpmCon *TPM) (bool, error)
	ReadNVPublic(tpmCon *TPM, index uint32) ([]byte, error)
	NVReadValue(tpmCon *TPM, index uint32, password string, size, offhandle uint32) ([]byte, error)
	ReadPCR(tpmCon *TPM, pcr uint32) ([]byte, error)

	// acpi.go
	GetACPITable(n string) ([]byte, error)

	// smbios.go
	IterateOverSMBIOSTables(n uint8, callback func(s *smbios.Structure) bool) (ret bool, err error)
}

// HwAPI The context object for low level hardware api
type HwAPI struct{}

// GetAPI Returns an initialized TxtApi object
func GetAPI() LowLevelHardwareInterfaces {
	return HwAPI{}
}
