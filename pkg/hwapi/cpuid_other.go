// +build !amd64

// Package hwapi provides access to low level hardware
package hwapi

//VersionString returns the vendor ID
func (h HwAPI) VersionString() string {
	return "null"
}

//HasSMX returns true if SMX is supported
func (h HwAPI) HasSMX() bool {
	return false
}

//HasVMX returns true if VMX is supported
func (h HwAPI) HasVMX() bool {
	return false
}

//HasMTRR returns true if MTRR are supported
func (h HwAPI) HasMTRR() bool {
	return false
}

//ProcessorBrandName returns the CPU brand name
func (h HwAPI) ProcessorBrandName() string {
	return "not intel"
}

//CPUSignature returns CPUID=1 eax
func (h HwAPI) CPUSignature() uint32 {
	return 0
}

//CPUSignatureFull returns CPUID=1 eax, ebx, ecx, edx
func (h HwAPI) CPUSignatureFull() (uint32, uint32, uint32, uint32) {
	return 0, 0, 0, 0
}

//CPULogCount returns number of logical CPU cores
func (h HwAPI) CPULogCount() uint32 {
	return 0
}
