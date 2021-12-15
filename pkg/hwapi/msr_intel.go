package hwapi

//Model specific registers
const (
	msrSMBase             int64 = 0x9e //nolint
	msrMTRRCap            int64 = 0xfe
	msrSMRRPhysBase       int64 = 0x1F2
	msrSMRRPhysMask       int64 = 0x1F3
	msrFeatureControl     int64 = 0x3A
	msrPlatformID         int64 = 0x17
	msrIA32DebugInterface int64 = 0xC80
)

// IA32Debug feature msr
type IA32Debug struct {
	Enabled  bool
	Locked   bool
	PCHStrap bool
}

//HasSMRR returns true if the CPU supports SMRR
func HasSMRR(h LowLevelHardwareInterfaces) (bool, error) {
	mtrrcap := h.ReadMSR(msrMTRRCap)

	return (mtrrcap[0]>>11)&1 != 0, nil
}

// SMRR for the SMM code.
type SMRR struct {
	Active   bool
	PhysBase uint64
	PhysMask uint64
}

// GetSMRRInfo returns SMRR config of the platform
func GetSMRRInfo(h LowLevelHardwareInterfaces) (SMRR, error) {
	var ret SMRR

	smrrPhysbase := h.ReadMSR(msrSMRRPhysBase)

	smrrPhysmask := h.ReadMSR(msrSMRRPhysMask)

	ret.Active = (smrrPhysmask[0]>>11)&1 != 0
	ret.PhysBase = (smrrPhysbase[0] >> 12) & 0xfffff
	ret.PhysMask = (smrrPhysmask[0] >> 12) & 0xfffff

	return ret, nil
}

//IA32FeatureControlIsLocked returns true if the IA32_FEATURE_CONTROL msr is locked
func IA32FeatureControlIsLocked(h LowLevelHardwareInterfaces) (bool, error) {
	featCtrl := h.ReadMSR(msrFeatureControl)

	return featCtrl[0]&1 != 0, nil
}

//IA32PlatformID returns the IA32_PLATFORM_ID msr
func IA32PlatformID(h LowLevelHardwareInterfaces) (uint64, error) {
	pltID := h.ReadMSR(msrPlatformID)

	return pltID[0], nil
}

//AllowsVMXInSMX returns true if VMX is allowed in SMX
func AllowsVMXInSMX(h LowLevelHardwareInterfaces) (bool, error) {
	featCtrl := h.ReadMSR(msrFeatureControl)

	var mask uint64 = (1 << 1) & (1 << 5) & (1 << 6)
	return (mask & featCtrl[0]) == mask, nil
}

//TXTLeavesAreEnabled returns true if all TXT leaves are enabled
func TXTLeavesAreEnabled(h LowLevelHardwareInterfaces) (bool, error) {
	featCtrl := h.ReadMSR(msrFeatureControl)

	txtBits := (featCtrl[0] >> 8) & 0x1ff
	return (txtBits&0xff == 0xff) || (txtBits&0x100 == 0x100), nil
}

//IA32DebugInterfaceEnabledOrLocked returns the enabled, locked and pchStrap state of IA32_DEBUG_INTERFACE msr
func IA32DebugInterfaceEnabledOrLocked(h LowLevelHardwareInterfaces) (*IA32Debug, error) {
	var debugMSR IA32Debug
	debugInterfaceCtrl := h.ReadMSR(msrIA32DebugInterface)

	debugMSR.Enabled = (debugInterfaceCtrl[0]>>0)&1 != 0
	debugMSR.Locked = (debugInterfaceCtrl[0]>>30)&1 != 0
	debugMSR.PCHStrap = (debugInterfaceCtrl[0]>>31)&1 != 0
	return &debugMSR, nil
}
