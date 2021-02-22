package hwapi

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/digitalocean/go-smbios/smbios"
)

type smbiosHeader struct {
	Type   uint8
	Length uint8
	Handle uint16
}

//SMBIOSType0 as defined in SMBIOS 2.0
type smbiosType0Raw20 struct {
	Vendor              uint8
	BIOSVersion         uint8
	BIOSStartingAddress uint16
	BIOSReleaseDate     uint8
	BIOSSize            uint8
	BiosCharacteristics uint32
}

//SMBIOSType0 as defined in SMBIOS 2.4
type smbiosType0Raw24 struct {
	SystemBiosMajor         uint8
	SystemBiosMinor         uint8
	EmbeddedControllerMajor uint8
	EmbeddedControllerMinor uint8
}

//SMBIOSType0 represents a decoded Type0 as defined in SMBIOS 2.4
type SMBIOSType0 struct {
	smbiosHeader
	Vendor                       string
	BIOSVersion                  string
	BIOSStartingAddress          int
	BIOSReleaseDate              string
	BIOSSize                     int
	BiosCharacteristics          uint32
	BiosCharacteristicsExtension []uint8
	SystemBiosMajor              uint8
	SystemBiosMinor              uint8
	EmbeddedControllerMajor      uint8
	EmbeddedControllerMinor      uint8
}

// IterateOverSMBIOSTables calls the callback for every SMBIOS table of specified type
func (h HwAPI) IterateOverSMBIOSTables(n uint8, callback func(s *smbios.Structure) bool) (ret bool, err error) {
	// Find SMBIOS data in operating system-specific location.
	var rc io.ReadCloser
	rc, _, err = smbios.Stream()
	if err != nil {
		return
	}
	// Be sure to close the stream!
	defer rc.Close()

	// Decode SMBIOS structures from the stream.
	d := smbios.NewDecoder(rc)
	ss, err := d.Decode()

	if err != nil {
		return
	}

	for _, s := range ss {
		// Only look at memory devices.
		if s.Header.Type != n {
			continue
		}
		ret = callback(s)
		if ret {
			return
		}
	}

	return
}

// IterateOverSMBIOSTablesType0 returns all SMBIOS tables of Type0 decoded
func (h HwAPI) IterateOverSMBIOSTablesType0(callback func(t0 *SMBIOSType0) bool) (ret bool, err error) {
	var err2 error
	ret, err = h.IterateOverSMBIOSTables(uint8(0), func(s *smbios.Structure) bool {
		var decoded SMBIOSType0

		buf := bytes.NewReader(s.Formatted)

		var raw smbiosType0Raw20
		var raw24 smbiosType0Raw24
		var extra []byte
		err2 = binary.Read(buf, binary.LittleEndian, &raw)
		if err2 != nil {
			return true
		}

		decoded.Type = s.Header.Type
		decoded.Length = s.Header.Length
		decoded.Handle = s.Header.Handle
		if int(raw.Vendor-1) < len(s.Strings) {
			decoded.Vendor = s.Strings[raw.Vendor-1]
		}
		if int(raw.BIOSVersion-1) < len(s.Strings) {
			decoded.BIOSVersion = s.Strings[raw.BIOSVersion-1]
		}
		decoded.BIOSStartingAddress = int(raw.BIOSStartingAddress) * 16
		if int(raw.BIOSReleaseDate-1) < len(s.Strings) {
			decoded.BIOSReleaseDate = s.Strings[raw.BIOSReleaseDate-1]
		}
		decoded.BIOSSize = (int(raw.BIOSSize) + 1) * 0x10000
		decoded.BiosCharacteristics = raw.BiosCharacteristics

		if int(s.Header.Length) >= binary.Size(raw)+binary.Size(raw24)+binary.Size(s.Header) {
			extrabytes := int(s.Header.Length) - (binary.Size(raw) + binary.Size(raw24) + binary.Size(s.Header))
			extra = make([]byte, extrabytes)

			err2 = binary.Read(buf, binary.LittleEndian, &extra)
			if err2 != nil {
				return true
			}
			err2 = binary.Read(buf, binary.LittleEndian, &raw24)
			if err2 != nil {
				return true
			}

			decoded.BiosCharacteristicsExtension = extra
			decoded.SystemBiosMajor = raw24.SystemBiosMajor
			decoded.SystemBiosMinor = raw24.SystemBiosMinor
			decoded.EmbeddedControllerMajor = raw24.EmbeddedControllerMajor
			decoded.EmbeddedControllerMinor = raw24.EmbeddedControllerMinor
		}
		return callback(&decoded)
	})
	if err == nil && err2 != nil {
		err = err2
	}

	return
}
