package hwapi

import (
	"os"
	"testing"
)

func TestE820ReservedCheck(t *testing.T) {
	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}

	ranges := []struct {
		start uint64
		end   uint64
	}{
		{0xf0000, 0xfffff},
		{0xfffc0000, 0xffffffff},
	}

	for _, s := range ranges {
		reserved, err := h.IsReservedInE820(s.start, s.end)
		if err != nil {
			t.Errorf("Checking range %x-%x failed: %s", s.start, s.end, err)
		}

		t.Logf("range %x-%x: %t", s.start, s.end, reserved)
	}
}
