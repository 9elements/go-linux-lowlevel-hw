package hwapi

import (
	"github.com/micgor32/go-msr"
)

// ReadMSR returns the MSR on core #0
func (h HwAPI) ReadMSR(msrAddr int64) uint64 {
	ret, err := msr.ReadMSR(0, msrAddr)
	if err != nil {
		return 0xff
	}

	return ret
}
