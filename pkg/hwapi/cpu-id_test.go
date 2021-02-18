package hwapi

import (
	"os"
	"testing"
)

func TestVersionStringQemu(t *testing.T) {

	h := GetAPI()
	if os.Getenv("RUN_IN_QEMU") != "TRUE" {
		t.Skip("Not running on QEMU")
	}

	got := h.VersionString()
	if got != "GenuineIntel" && got != "AuthenticAMD" {
		t.Errorf("Got unexpected version string %s", got)
	}
	got = h.ProcessorBrandName()
	if got != "QEMU Virtual CPU version 2.5+" {
		t.Errorf("Got unexpected brand name string %s", got)
	}
}

func TestVersionString(t *testing.T) {

	h := GetAPI()

	got := h.VersionString()

	if got == "" {
		t.Error("VersionString() returned the empty string.")
	}
}
