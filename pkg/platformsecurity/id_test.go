package platformsecurity

import (
	"strings"
	"testing"

	"github.com/klauspost/cpuid/v2"
)

func TestVendorIDString(t *testing.T) {
	for id := IDUndefined; id < EndOfID; id++ {
		if strings.Contains(id.String(), "unknown") {
			t.Fatalf("id %d has no String", id)
		}
	}
	if !strings.Contains(EndOfID.String(), "unknown") {
		t.Fatalf("EndOfID has String")
	}
}

func TestVendorIDCPUVendorID(t *testing.T) {
	for id := IDUndefined + 1; id < EndOfID; id++ {
		if id.CPUVendorID() == cpuid.VendorUnknown {
			t.Fatalf("id %d has no CPUVendorID", id)
		}
	}
	if EndOfID.CPUVendorID() != cpuid.VendorUnknown {
		t.Fatalf("EndOfID has CPUVendorID")
	}
}
