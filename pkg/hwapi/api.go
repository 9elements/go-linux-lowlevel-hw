package hwapi

import "github.com/digitalocean/go-smbios/smbios"

//LowLevelHardwareInterfaces provides methods to access hardware found on modern x86_64 platforms
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

	// e820.go
	IsReservedInE820(start uint64, end uint64) (bool, error)
	UsableMemoryAbove4G() (size uint64, err error)
	UsableMemoryBelow4G() (size uint64, err error)

	// iommu.go
	LookupIOAddress(addr uint64, regs VTdRegisters) ([]uint64, error)
	AddressRangesIsDMAProtected(first, end uint64) (bool, error)

	// msr.go
	ReadMSR(msr int64, core int) (uint64, error)

	// msr_intel.go
	HasSMRR() (bool, error)
	GetSMRRInfo() (SMRR, error)
	IA32FeatureControlIsLocked() (bool, error)
	IA32PlatformID() (uint64, error)
	AllowsVMXInSMX() (bool, error)
	TXTLeavesAreEnabled() (bool, error)
	IA32DebugInterfaceEnabledOrLocked() (*IA32Debug, error)

	// pci.go
	PCIEnumerateVisibleDevices(cb func(d PCIDevice) (abort bool)) (err error)
	PCIReadConfigSpace(d PCIDevice, off int, out interface{}) error
	PCIWriteConfigSpace(d PCIDevice, off int, val interface{}) error
	PCIReadVendorID(d PCIDevice) (uint16, error)
	PCIReadDeviceID(d PCIDevice) (uint16, error)

	// hostbridge.go
	ReadHostBridgeTseg() (uint32, uint32, error)
	ReadHostBridgeDPR() (DMAProtectedRange, error)

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
	GetACPITableDevMem(n string) ([]byte, error)
	GetACPITableSysFS(n string) ([]byte, error)

	// smbios.go
	IterateOverSMBIOSTables(n uint8, callback func(s *smbios.Structure) bool) (ret bool, err error)
	IterateOverSMBIOSTablesType0(callback func(t0 *SMBIOSType0) bool) (ret bool, err error)
	IterateOverSMBIOSTablesType17(callback func(t17 *SMBIOSType17) bool) (ret bool, err error)
}

//HwAPI The context object for low level hardware api
type HwAPI struct{}

//GetAPI Returns an initialized TxtApi object
func GetAPI() LowLevelHardwareInterfaces {
	return HwAPI{}
}
