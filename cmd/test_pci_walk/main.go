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
			venid, err := h.PCIReadConfigSpace(d, 0, 2)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
			}
			devid, err := h.PCIReadConfigSpace(d, 2, 2)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
			}

			fmt.Printf("Found device: %02x:%02x.%x [%04x:%04x]\n",
				d.Bus, d.Device, d.Function, venid, devid)
			return false
		}); err != nil {
		fmt.Fprintf(os.Stderr, "PCIEnumerateVisibleDevices failed wiht: %v", err)
	}

}
