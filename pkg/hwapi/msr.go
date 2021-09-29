package hwapi

import (
	"github.com/fearful-symmetry/gomsr"
)

// ReadMSR returns the MSR on core #0
func (h HwAPI) ReadMSR(msr int64, core int) (uint64, error) {
	msrCtx, err := gomsr.MSR(core)
	if err != nil {
		return 0, err
	}
	msrData, err := msrCtx.Read(msr)
	if err != nil {
		return 0, err
	}

	return msrData, nil
}
