package hwapi

import (
	"fmt"
	"runtime"

	"github.com/fearful-symmetry/gomsr"
)

func (h HwApi) ReadMSR(msr int64) (uint64, error) {
	var data uint64
	for i := 0; i < runtime.NumCPU(); i++ {
		msrCtx, err := gomsr.MSR(i)
		if err != nil {
			return 0, fmt.Errorf("MSR: Selected core %d doesn't exist", i)
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
