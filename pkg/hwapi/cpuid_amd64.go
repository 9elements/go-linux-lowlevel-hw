//go:build amd64
// +build amd64

// Package hwapi provides access to low level hardware
package hwapi

import "github.com/intel-go/cpuid"

func cpuidLow(arg1, arg2 uint32) (eax, ebx, ecx, edx uint32) // implemented in cpuidlow_amd64.s

// VersionString returns the vendor ID
func (h HwAPI) VersionString() string {
	return cpuid.VendorIdentificatorString
}

// HasSMX returns true if SMX is supported
func (h HwAPI) HasSMX() bool {
	return cpuid.HasFeature(cpuid.SMX)
}

// HasVMX returns true if VMX is supported
func (h HwAPI) HasVMX() bool {
	return cpuid.HasFeature(cpuid.VMX)
}

// HasMTRR returns true if MTRR are supported
func (h HwAPI) HasMTRR() bool {
	return cpuid.HasFeature(cpuid.MTRR) || cpuid.HasExtraFeature(cpuid.MTRR_2)
}

// ProcessorBrandName returns the CPU brand name
func (h HwAPI) ProcessorBrandName() string {
	return cpuid.ProcessorBrandString
}

// CPUSignature returns CPUID=1 eax
func (h HwAPI) CPUSignature() uint32 {
	eax, _, _, _ := h.CPUSignatureFull()
	return eax
}

// CPUSignatureFull returns CPUID=1 eax, ebx, ecx, edx
func (h HwAPI) CPUSignatureFull() (uint32, uint32, uint32, uint32) {
	return cpuidLow(1, 0)
}

// CPULogCount returns number of logical CPU cores
func (h HwAPI) CPULogCount() uint32 {
	return 0
	// return uint32(cpuid.MaxLogicalCPUId)
}

// CPUID executes the CPUID instruction with the given leaf (eax) and subleaf (ecx) values
// Returns the resulting eax, ebx, ecx, and edx register values
func (h HwAPI) CPUID(leaf, subleaf uint32) (eax, ebx, ecx, edx uint32) {
	return cpuidLow(leaf, subleaf)
}
