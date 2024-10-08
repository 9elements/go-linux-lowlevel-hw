
![image](assets/logo.png)

go-linux-lowevel-hw provides low level access to common hardware on UNIX like platforms.

Description
-----------
This package provides low level access to certain hardware typically found
on modern x86 PCs. Some information are only available when run as
most priviledged user. Thus this library is to be used in preproduction
and testing envirnoments with relaxed kernel security.

**Be warned, you could brick your system.**

How to use this library
-----------------------

```
package main

import (

	"github.com/9elements/go-linux-lowlevel-hw/pkg/hwapi"
)

func main() {
	h := hwapi.GetAPI()

	//...
}
```

Interfaces
----------
The GetAPI call returns an interface providing the following methods:
```
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
	ReadMSRAllCores(msr int64) (uint64, error)

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
	PCIReadConfig8(d PCIDevice, off int) (uint8, error)
	PCIReadConfig16(d PCIDevice, off int) (uint16, error)
	PCIReadConfig32(d PCIDevice, off int) (uint32, error)
	PCIWriteConfig8(d PCIDevice, off int, val uint8) error
	PCIWriteConfig16(d PCIDevice, off int, val uint16) error
	PCIWriteConfig32(d PCIDevice, off int, val uint32) error
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
```
