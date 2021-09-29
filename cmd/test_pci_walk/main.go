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
			var venid, devid uint16
			if err := h.PCIReadConfigSpace(d, 0, &venid); err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
			}
			if err := h.PCIReadConfigSpace(d, 2, &devid); err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
			}

			fmt.Printf("Found device: %02x:%02x.%x [%04x:%04x]\n",
				d.Bus, d.Device, d.Function, venid, devid)
			return false
		}); err != nil {
		fmt.Fprintf(os.Stderr, "PCIEnumerateVisibleDevices failed wiht: %v", err)
	}

}
