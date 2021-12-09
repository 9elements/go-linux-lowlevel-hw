package platformsecurity

import (
	"fmt"

	"github.com/klauspost/cpuid/v2"
)

// ID is an unique ID of a platform security configuration.
// Platform may combine for example Intel CPU + OpenTitan as RTM.
//
// It differs from github.com/klauspost/cpuid.Vendor, because cpuid.Vendor
// defines CPU vendor, while a platform may combine for example
// Intel and Lattice (and this combination will have an additional ID).
type ID int

const (
	// IDUndefined is a ID reserved for the zero-value only.
	IDUndefined = ID(iota)

	// IDIntelTXT is an ID corresponds to "Intel" TXT (pre-CBnT) platforms.
	IDIntelTXT

	// IDIntelCBnT is an ID corresponds to "Intel" CBnT platforms.
	IDIntelCBnT

	// IDAMDMilan is an ID corresponds to "AMD Milan".
	IDAMDMilan

	// EndOfID is a limiter for loops to iterate over ID-s.
	EndOfID
)

// String implements fmt.Stringer.
func (id ID) String() string {
	switch id {
	case IDUndefined:
		return "<undefined>"
	case IDIntelTXT:
		return "Intel TXT"
	case IDIntelCBnT:
		return "Intel CBnT"
	case IDAMDMilan:
		return "AMD Milan"
	}
	return fmt.Sprintf("unknown_ID_%d", id)
}

// CPUVendorID return the vendor ID of the CPU used on the platform.
func (id ID) CPUVendorID() cpuid.Vendor {
	switch id {
	case IDIntelTXT, IDIntelCBnT:
		return cpuid.Intel
	case IDAMDMilan:
		return cpuid.AMD
	}
	return cpuid.VendorUnknown
}
