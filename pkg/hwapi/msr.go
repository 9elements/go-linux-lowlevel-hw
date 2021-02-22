package hwapi

import (
	"fmt"
	"runtime"

	"github.com/fearful-symmetry/gomsr"
)

// ReadMSR returns the MSR on core #0
func (h HwAPI) ReadMSR(msr int64) (uint64, error) {
	msrCtx, err := gomsr.MSR(0)
	if err != nil {
		return 0, err
	}
	msrData, err := msrCtx.Read(msr)
	if err != nil {
		return 0, err
	}

	return msrData, nil
}

// ReadMSRAllCores returns the MSR and errors if it's not the same on all cores
func (h HwAPI) ReadMSRAllCores(msr int64) (uint64, error) {
	var data uint64
	for i := 0; i < runtime.NumCPU(); i++ {
		msrCtx, err := gomsr.MSR(i)
		if err != nil {
			return 0, err
		}
		msrData, err := msrCtx.Read(msr)
		if err != nil {
			return 0, err
		}
		if i != 0 {
			if data != msrData {
				return 0, fmt.Errorf("MSR: cores of MSR 0x%x non equal", msr)
			}
		}
		data = msrData
	}
	return data, nil
}
