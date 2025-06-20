package hwapi

import (
	"fmt"
	"os"

	"github.com/fearful-symmetry/gomsr"
)

// ReadMSR returns the MSR on core #0
func (h HwAPI) ReadMSR(msr int64) []uint64 {
	count := 0
	var ret []uint64
	for {
		fmt.Printf("ReadMSR - gomsr.MSR count: %d\n", count)
		msrCtx, err := gomsr.MSR(count)
		if err != nil {
			fmt.Fprintf(os.Stdout, "ReadMSR - gomsr.MSR context aborted with: %v\n", err)
			break
		}
		msrData, err := msrCtx.Read(msr)
		if err != nil {
			fmt.Fprintf(os.Stdout, "ReadMSR - msrCtx.Read aborted with: %v\n", err)
			break
		}
		ret = append(ret, msrData)
		count++
	}

	fmt.Printf("ReadMSR - gomsr.MSR count: %d, ret len: %d\n", count, len(ret))
	return ret
}
