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
	CPULogCount() uint32

	// e820.go
	IsReservedInE820(start uint64, end uint64) (bool, error)

	// iommu.go
	LookupIOAddress(addr uint64, regs VTdRegisters) ([]uint64, error)
	AddressRangesIsDMAProtected(first, end uint64) (bool, error)

	// msr.go
	ReadMSR(msr int64) (uint64, error)

	// msr_intel.go
	HasSMRR() (bool, error)
	GetSMRRInfo() (SMRR, error)
	IA32FeatureControlIsLocked() (bool, error)
	IA32PlatformID() (uint64, error)
	AllowsVMXInSMX() (bool, error)
	TXTLeavesAreEnabled() (bool, error)
	IA32DebugInterfaceEnabledOrLocked() (*IA32Debug, error)

	// pci.go
	PCIReadConfigSpace(bus int, device int, devFn int, off int, buf interface{}) error
	PCIWriteConfigSpace(bus int, device int, devFn int, off int, buf interface{}) error
	PCIReadConfig8(bus int, device int, devFn int, off int) (uint8, error)
	PCIReadConfig16(bus int, device int, devFn int, off int) (uint16, error)
	PCIReadConfig32(bus int, device int, devFn int, off int) (uint32, error)
	PCIWriteConfig8(bus int, device int, devFn int, off int, val uint8) error
	PCIWriteConfig16(bus int, device int, devFn int, off int, val uint16) error
	PCIWriteConfig32(bus int, device int, devFn int, off int, val uint32) error
	PCIReadVendorID(bus int, device int, devFn int) (uint16, error)
	PCIReadDeviceID(bus int, device int, devFn int) (uint16, error)

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
}

//HwApi The context object for low level hardware api
type HwApi struct{}

//GetAPI Returns an initialized TxtApi object
func GetAPI() LowLevelHardwareInterfaces {
	return HwApi{}
}
