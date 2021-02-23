package main

import (
	"fmt"
	"os"

	"github.com/9elements/go-linux-lowlevel-hw/pkg/hwapi"
)

func main() {
	h := hwapi.GetAPI()
	if err := h.PCIEnumerateVisibleDevices(
		func(d hwapi.PCIDevice) (abort bool) {
			vend, _ := h.PCIReadVendorID(d)
			devid, _ := h.PCIReadDeviceID(d)
			fmt.Printf("Found device: %02x:%02x.%x [%04x:%04x]\n",
				d.Bus, d.Device, d.Function, vend, devid)
			return false
		}); err != nil {
		fmt.Fprintf(os.Stderr, "PCIEnumerateVisibleDevices failed wiht: %v", err)
	}

}
