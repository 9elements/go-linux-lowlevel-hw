package hwapi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	acpiSysfsPath       = "/sys/firmware/acpi/tables"
	acpiSysfsSystabPath = "/sys/firmware/efi/systab"

	biosRomBase = 0xe0000
	biosRomSize = 0x20000
	ebdaTop     = 0xa0000
)

// ACPIRsdpRev1 as defined in ACPI Spec 1
type ACPIRsdpRev1 struct {
	Signature [8]uint8
	Checksum  uint8
	OEMID     [6]uint8
	Revision  uint8
	RSDTPtr   uint32
}

// ACPIRsdp as defined in ACPI Spec 6.2 "5.2.5.3 Root System Description Pointer (RSDP) Structure"
type ACPIRsdp struct {
	ACPIRsdpRev1

	RSDPLen          uint32
	XSDTPtr          uint64
	ExtendedChecksum uint8
	Reserved         [3]uint8
}

type acpiHeader struct {
	Signature       [4]uint8
	Length          uint32
	Revision        uint8
	Checksum        uint8
	OEMID           [6]uint8
	OEMTableID      [8]uint8
	OEMRevision     uint32
	CreatorID       uint32
	CreatorRevision uint32
}

// ACPIRsdt as defined in ACPI Spec 6.2 "5.2.7 Root System Description Table (RSDT)"
type acpiRsdt struct {
	acpiHeader
	// Entry           []uint32 count depend on Length field
}

// ACPIXsdt as defined in ACPI Spec 6.2 "5.2.8 Extended System Description Table (XSDT)"
type acpiXsdt struct {
	acpiHeader
	// Entry           []uint64 count depend on Length field
}

func GetACPITableSysFS(h LowLevelHardwareInterfaces, n string) ([]byte, error) {
	buf, err := os.ReadFile(fmt.Sprintf("%s/%s", acpiSysfsPath, n))
	if err != nil {
		return nil, fmt.Errorf("cannot access sysfs path %s: %s", acpiSysfsPath, err)
	}
	return buf, nil
}

var (
	backupRSDT     []byte
	backupRSDTList []uint32
)

func getACPITableDevMemRSDT(address uint32, l LowLevelHardwareInterfaces) ([]uint32, []byte, error) {
	var rsdt acpiRsdt
	var hdrs []uint32

	if len(backupRSDT) > 0 {
		return backupRSDTList, backupRSDT, nil
	}

	buf := make([]byte, binary.Size(rsdt))
	err := l.ReadPhysBuf(int64(address), buf)
	if err != nil {
		return nil, nil, err
	}
	err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &rsdt)
	if err != nil {
		return nil, nil, err
	}

	if string(rsdt.Signature[:]) != "RSDT" {
		return nil, nil, fmt.Errorf("RSDT has invalid signature")
	}
	if rsdt.Length == 0 || rsdt.Length == 0xffffffff ||
		(rsdt.Length-uint32(binary.Size(acpiHeader{})))%4 > 0 {
		return nil, nil, fmt.Errorf("RSDT has invalid length")
	}
	buf = make([]byte, (rsdt.Length - uint32(binary.Size(acpiHeader{}))))
	err = l.ReadPhysBuf(int64(address)+int64(binary.Size(acpiHeader{})), buf)
	if err != nil {
		return nil, nil, err
	}

	hdrs = make([]uint32, len(buf)/4)
	err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &hdrs)
	if err != nil {
		return nil, nil, err
	}

	buf = make([]byte, rsdt.Length)
	err = l.ReadPhysBuf(int64(address), buf)
	if err != nil {
		return nil, nil, err
	}

	backupRSDTList = hdrs
	backupRSDT = buf

	return hdrs, buf, nil
}

var (
	backupXSDT     []byte
	backupXSDTList []uint64
)

func getACPITableDevMemXSDT(address uint64, l LowLevelHardwareInterfaces) ([]uint64, []byte, error) {
	var xsdt acpiXsdt
	var hdrs []uint64

	if len(backupXSDT) > 0 {
		return backupXSDTList, backupXSDT, nil
	}

	buf := make([]byte, binary.Size(xsdt))
	err := l.ReadPhysBuf(int64(address), buf)
	if err != nil {
		return nil, nil, err
	}
	err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &xsdt)
	if err != nil {
		return nil, nil, err
	}

	if string(xsdt.Signature[:]) != "XSDT" {
		return nil, nil, fmt.Errorf("XSDT has invalid signature")
	}
	if xsdt.Length == 0 || xsdt.Length == 0xffffffff ||
		(xsdt.Length-uint32(binary.Size(acpiHeader{})))%8 > 0 {
		return nil, nil, fmt.Errorf("XSDT has invalid length")
	}
	buf = make([]byte, (xsdt.Length - uint32(binary.Size(acpiHeader{}))))
	err = l.ReadPhysBuf(int64(address)+int64(binary.Size(acpiHeader{})), buf)
	if err != nil {
		return nil, nil, err
	}

	hdrs = make([]uint64, len(buf)/8)
	err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &hdrs)
	if err != nil {
		return nil, nil, err
	}

	buf = make([]byte, xsdt.Length)
	err = l.ReadPhysBuf(int64(address), buf)
	if err != nil {
		return nil, nil, err
	}

	backupXSDTList = hdrs
	backupXSDT = buf

	return hdrs, buf, nil
}

var (
	backupRSDP    ACPIRsdp
	backupRawRSDP []byte
)

func parseSystab(l LowLevelHardwareInterfaces) ([]byte, ACPIRsdp, error) {
	var rsdp ACPIRsdp

	buf := make([]byte, binary.Size(rsdp))

	// Try to get the RSDP address from systab
	file, err := os.Open(acpiSysfsSystabPath)
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if !strings.Contains(scanner.Text(), "=") {
				continue
			}

			kv := strings.SplitAfter(scanner.Text(), "=")

			if kv[0] == "ACPI20=" || kv[0] == "ACPI=" {
				address, err := strconv.ParseUint(kv[1], 0, 64)
				if err != nil || address == 0 {
					continue
				}
				log.Printf("address =%x", address)

				err = l.ReadPhysBuf(int64(address), buf)
				if err != nil {
					log.Printf("%v", err)
					continue
				}
				err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &rsdp)
				if err != nil {
					return nil, rsdp, err
				}
				log.Printf("signature %s", string(rsdp.Signature[:]))
				if string(rsdp.Signature[:]) == "RSD PTR " {
					return buf, rsdp, nil
				}
			}
		}
	}
	return nil, rsdp, fmt.Errorf("RSDP not found in systab")
}

func scanLowMem(l LowLevelHardwareInterfaces) ([]byte, ACPIRsdp, error) {
	var rsdp ACPIRsdp

	buf := make([]byte, binary.Size(rsdp))

	for i := int64(biosRomBase); i < biosRomBase+biosRomSize-int64(binary.Size(rsdp)); i += 16 {
		err := l.ReadPhysBuf(i, buf)
		if err != nil {
			return nil, rsdp, fmt.Errorf("failed to read physical memory at %x: %v", i, err)
		}
		err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &rsdp)
		if err != nil {
			return nil, rsdp, err
		}
		if i == 0xF05B0 {
			log.Printf("string(rsdp.Signature[:]) = %s", string(rsdp.Signature[:]))
		}
		if string(rsdp.Signature[:]) == "RSD PTR " {
			return buf, rsdp, nil
		}
	}
	return nil, rsdp, fmt.Errorf("RSDP not found in low memory")
}

func scanEBDA(l LowLevelHardwareInterfaces) ([]byte, ACPIRsdp, error) {
	var rsdp ACPIRsdp

	buf := make([]byte, binary.Size(rsdp))

	for i := int64(ebdaTop - biosRomSize); i < ebdaTop-int64(binary.Size(rsdp)); i += 16 {
		err := l.ReadPhysBuf(i, buf)
		if err != nil {
			log.Printf("%v", err)

			return nil, rsdp, fmt.Errorf("failed to read physical memory at %x: %v", i, err)
		}
		err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &rsdp)
		if err != nil {
			return nil, rsdp, err
		}

		if string(rsdp.Signature[:]) == "RSD PTR " {
			break
		}
	}

	return nil, rsdp, fmt.Errorf("RSDP not found in low memory")
}

func scanReservedMem(l LowLevelHardwareInterfaces) ([]byte, ACPIRsdp, error) {
	var rsdp ACPIRsdp

	buf := make([]byte, binary.Size(rsdp))

	_, err := l.IterateOverE820Ranges("ACPI Tables", func(start uint64, end uint64) bool {
		for i := int64(start); i < int64(end)-int64(binary.Size(rsdp)); i += 16 {
			err := l.ReadPhysBuf(i, buf)
			if err != nil {
				log.Printf("%v", err)

				return false
			}
			err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &rsdp)
			if err != nil {
				return false
			}

			if string(rsdp.Signature[:]) == "RSD PTR " {
				return true
			}
		}
		return false
	})
	if err != nil {
		return nil, rsdp, err
	}
	return nil, rsdp, fmt.Errorf("RSDP not found in low memory")
}

func verifyRSDP(buf []byte, rsdp ACPIRsdp) error {
	var old ACPIRsdpRev1

	if rsdp.Revision == 0 {
		if rsdp.RSDPLen != uint32(binary.Size(old)) {
			return fmt.Errorf("ACPI RSDP has unexpected length %d", rsdp.RSDPLen)
		}
	}
	/* Validate first checksum */
	chksum := byte(0)
	for i, b := range buf {
		if i >= binary.Size(old) {
			break
		}
		chksum = chksum + b
	}
	if chksum > 0 {
		return fmt.Errorf("ACPI RSDP has invalid checksum")
	}

	if rsdp.Revision == 0 {
		return nil
	}
	if rsdp.Revision == 2 {
		if rsdp.RSDPLen != uint32(binary.Size(rsdp)) {
			return fmt.Errorf("ACPI RSDP has unexpected length %d", rsdp.RSDPLen)
		}
	}
	/* Validate second checksum */
	chksum = byte(0)
	for _, i := range buf {
		chksum = chksum + i
	}

	if chksum > 0 {
		return fmt.Errorf("ACPI RSDP has invalid checksum")
	}
	return nil
}

func getACPITableDevMemRSDP(l LowLevelHardwareInterfaces) ([]byte, ACPIRsdp, error) {
	var rsdp ACPIRsdp

	if string(backupRSDP.Signature[:]) == "RSD PTR " {
		return backupRawRSDP, backupRSDP, nil
	}

	// Try to get the RSDP address from systab
	data, rsdp, err := parseSystab(l)
	if err == nil {
		err = verifyRSDP(data, rsdp)
		if err == nil {
			backupRSDP = rsdp
			backupRawRSDP = data
			return data, rsdp, nil
		}
	}

	data, rsdp, err = scanLowMem(l)
	if err == nil {
		err = verifyRSDP(data, rsdp)
		if err == nil {
			backupRSDP = rsdp
			backupRawRSDP = data
			return data, rsdp, nil
		}
	}

	data, rsdp, err = scanEBDA(l)
	if err == nil {
		err = verifyRSDP(data, rsdp)
		if err == nil {
			backupRSDP = rsdp
			backupRawRSDP = data
			return data, rsdp, nil
		}
	}

	data, rsdp, err = scanReservedMem(l)

	if err == nil {
		err = verifyRSDP(data, rsdp)
		if err == nil {
			backupRSDP = rsdp
			backupRawRSDP = data
			return data, rsdp, nil
		}
	}

	return nil, rsdp, fmt.Errorf("RSDP not found")
}

func GetACPITableDevMem(h LowLevelHardwareInterfaces, n string) ([]byte, error) {
	rsdpBuf, rsdp, err := getACPITableDevMemRSDP(h)
	if err != nil {
		return nil, err
	}

	if string(rsdp.Signature[:]) != "RSD PTR " {
		return nil, fmt.Errorf("RSDP not found")
	}

	if n == "RSDP" {
		return rsdpBuf, nil
	}

	rsdtHeaders, rsdtBuf, err1 := getACPITableDevMemRSDT(rsdp.RSDTPtr, h)
	if err1 == nil && n == "RSDT" {
		return rsdtBuf, nil
	}

	xsdtHeaders, xsdtBuf, err2 := getACPITableDevMemXSDT(rsdp.XSDTPtr, h)
	if err2 == nil && n == "XSDT" {
		return xsdtBuf, nil
	}

	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("RSDT and XSDT are invalid")
	}

	buf := make([]byte, binary.Size(acpiHeader{}))

	acpitables := map[uint64]string{}
	if len(rsdtHeaders) > 0 {
		for i := range rsdtHeaders {
			var header acpiHeader
			if _, ok := acpitables[uint64(rsdtHeaders[i])]; ok {
				continue
			}

			err := h.ReadPhysBuf(int64(rsdtHeaders[i]), buf)
			if err != nil {
				return nil, err
			}
			err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &header)
			if err != nil {
				return nil, err
			}
			acpitables[uint64(rsdtHeaders[i])] = string(header.Signature[:])
		}
	}

	if len(xsdtHeaders) > 1 {
		for i := range xsdtHeaders {
			var header acpiHeader
			if _, ok := acpitables[xsdtHeaders[i]]; ok {
				continue
			}

			err := h.ReadPhysBuf(int64(xsdtHeaders[i]), buf)
			if err != nil {
				return nil, err
			}
			err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &header)
			if err != nil {
				return nil, err
			}
			acpitables[xsdtHeaders[i]] = string(header.Signature[:])
		}
	}

	for k, v := range acpitables {
		if v == n { // FIXME: Handle duplicated entries like SSDT
			var header acpiHeader

			err := h.ReadPhysBuf(int64(k), buf)
			if err != nil {
				return nil, err
			}
			err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &header)
			if err != nil {
				return nil, err
			}
			buf = make([]byte, header.Length)
			err = h.ReadPhysBuf(int64(k), buf)
			if err != nil {
				return nil, err
			}

			return buf, nil
		}
	}
	return nil, fmt.Errorf("ACPI table not found")
}

// GetACPITable returns the requested ACPI table, for DSDT use argument "DSDT"
func (h HwAPI) GetACPITable(n string) ([]byte, error) {
	if n == "" || len(n) > 6 {
		return nil, fmt.Errorf("invalid ACPI name")
	}

	// Try SYSFS first, but it doesn't has RSDP
	tbl, err := GetACPITableSysFS(h, n)
	if err != nil {
		tbl, err = GetACPITableDevMem(h, n)
	}
	return tbl, err
}
