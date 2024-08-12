package hwapi

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// iterateOverE820Ranges iterates over all e820 entries and invokes the callback for every matching type
func (h HwAPI) IterateOverE820Ranges(target string, callback func(start uint64, end uint64) bool) (bool, error) {
	dir, err := os.Open("/sys/firmware/memmap")
	if err != nil {
		return false, fmt.Errorf("cannot access e820 table: %s", err)
	}

	subdirs, err := dir.Readdir(0)
	if err != nil {
		return false, fmt.Errorf("cannot access e820 table: %s", err)
	}

	for _, subdir := range subdirs {
		if subdir.IsDir() {

			path := fmt.Sprintf("/sys/firmware/memmap/%s/type", subdir.Name())
			buf, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			if strings.Contains(strings.ToLower(string(buf)), strings.ToLower(target)) {
				path := fmt.Sprintf("/sys/firmware/memmap/%s/start", subdir.Name())
				thisStart, err := readHexInteger(path)
				if err != nil {
					continue
				}

				path = fmt.Sprintf("/sys/firmware/memmap/%s/end", subdir.Name())
				thisEnd, err := readHexInteger(path)
				if err != nil {
					continue
				}

				if callback(thisStart, thisEnd) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// UsableMemoryAbove4G returns the usable memory above 4GiB
func UsableMemoryAbove4G(l LowLevelHardwareInterfaces) (size uint64, err error) {
	_, err = l.IterateOverE820Ranges("system ram", func(rstart uint64, rend uint64) bool {
		if rstart > 0xffffffff {
			size += (rend - rstart)
		}
		return false
	})

	return
}

// UsableMemoryBelow4G returns the usable memory below 4GiB
func UsableMemoryBelow4G(l LowLevelHardwareInterfaces) (size uint64, err error) {
	_, err = l.IterateOverE820Ranges("system ram", func(rstart uint64, rend uint64) bool {
		if rstart <= 0xffffffff {
			size += (rend - rstart)
		}
		return false
	})

	return
}

// IsReservedInE820 reads the e820 table exported via /sys/firmware/memmap and checks whether
// the range [start; end] is marked as reserved. Returns true if it is reserved,
// false if not.
func IsReservedInE820(l LowLevelHardwareInterfaces, start uint64, end uint64) (bool, error) {
	if start > end {
		return false, fmt.Errorf("invalid range")
	}

	contains, err := l.IterateOverE820Ranges("reserved", func(rstart uint64, rend uint64) bool {
		if rstart <= start && rend >= end {
			return true
		}
		return false
	})
	if err != nil {
		return false, err
	}
	return contains, nil
}

func readHexInteger(path string) (uint64, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	ret, err := strconv.ParseUint(string(buf[:len(buf)-1]), 0, 64)
	if err != nil {
		return 0, err
	}

	return ret, nil
}
