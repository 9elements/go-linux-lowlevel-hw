package hwapi

import (
	"fmt"
)

//Model specific registers
const (
	msrSMBase             int64 = 0x9e
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
func (h HwAPI) HasSMRR() (bool, error) {
	mtrrcap, err := h.ReadMSR(msrMTRRCap)
	if err != nil {
		return false, fmt.Errorf("Cannot access MSR IA32_MTRRCAP: %s", err)
	}

	return (mtrrcap>>11)&1 != 0, nil
}

// SMRR for the SMM code.
type SMRR struct {
	Active   bool
	PhysBase uint64
	PhysMask uint64
}

// GetSMRRInfo returns SMRR config of the platform
func (h HwAPI) GetSMRRInfo() (SMRR, error) {
	var ret SMRR

	smrrPhysbase, err := h.ReadMSR(msrSMRRPhysBase)
	if err != nil {
		return ret, fmt.Errorf("Cannot access MSR IA32_SMRR_PHYSBASE: %s", err)
	}

	smrrPhysmask, err := h.ReadMSR(msrSMRRPhysMask)
	if err != nil {
		return ret, fmt.Errorf("Cannot access MSR IA32_SMRR_PHYSMASK: %s", err)
	}

	ret.Active = (smrrPhysmask>>11)&1 != 0
	ret.PhysBase = (smrrPhysbase >> 12) & 0xfffff
	ret.PhysMask = (smrrPhysmask >> 12) & 0xfffff

	return ret, nil
}

//IA32FeatureControlIsLocked returns true if the IA32_FEATURE_CONTROL msr is locked
func (h HwAPI) IA32FeatureControlIsLocked() (bool, error) {
	featCtrl, err := h.ReadMSR(msrFeatureControl)
	if err != nil {
		return false, fmt.Errorf("Cannot access MSR IA32_FEATURE_CONTROL: %s", err)
	}

	return featCtrl&1 != 0, nil
}

//IA32PlatformID returns the IA32_PLATFORM_ID msr
func (h HwAPI) IA32PlatformID() (uint64, error) {
	pltID, err := h.ReadMSR(msrPlatformID)
	if err != nil {
		return 0, fmt.Errorf("Cannot access MSR IA32_PLATFORM_ID: %s", err)
	}

	return pltID, nil
}

//AllowsVMXInSMX returns true if VMX is allowed in SMX
func (h HwAPI) AllowsVMXInSMX() (bool, error) {
	featCtrl, err := h.ReadMSR(msrFeatureControl)
	if err != nil {
		return false, fmt.Errorf("Cannot access MSR IA32_FEATURE_CONTROL: %s", err)
	}

	var mask uint64 = (1 << 1) & (1 << 5) & (1 << 6)
	return (mask & featCtrl) == mask, nil
}

//TXTLeavesAreEnabled returns true if all TXT leaves are enabled
func (h HwAPI) TXTLeavesAreEnabled() (bool, error) {
	featCtrl, err := h.ReadMSR(msrFeatureControl)
	if err != nil {
		return false, fmt.Errorf("Cannot access MSR IA32_FEATURE_CONTROL: %s", err)
	}

	txtBits := (featCtrl >> 8) & 0x1ff
	return (txtBits&0xff == 0xff) || (txtBits&0x100 == 0x100), nil
}

//IA32DebugInterfaceEnabledOrLocked returns the enabled, locked and pchStrap state of IA32_DEBUG_INTERFACE msr
func (h HwAPI) IA32DebugInterfaceEnabledOrLocked() (*IA32Debug, error) {
	var debugMSR IA32Debug
	debugInterfaceCtrl, err := h.ReadMSR(msrIA32DebugInterface)
	if err != nil {
		return nil, fmt.Errorf("Cannot access MSR IA32_DEBUG_INTERFACE: %s", err)
	}

	debugMSR.Enabled = (debugInterfaceCtrl>>0)&1 != 0
	debugMSR.Locked = (debugInterfaceCtrl>>30)&1 != 0
	debugMSR.PCHStrap = (debugInterfaceCtrl>>31)&1 != 0
	return &debugMSR, nil
}
