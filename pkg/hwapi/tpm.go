package hwapi

import (
	"fmt"
	"strings"

	tpm1 "github.com/google/go-tpm/tpm"
	tpm2 "github.com/google/go-tpm/legacy/tpm2"
)

const (
	tpm2LockedResult = "error code 0x22"
)

// NewTPM Looks for a TPM device, returns it if one is found
func (h HwAPI) NewTPM() (*TPM, error) {
	tpm, err := NewTPM()
	if err != nil {
		return nil, err
	}
	return tpm, nil
}

// NVLocked returns true if the NV RAM is locked, otherwise false
func (h HwAPI) NVLocked(tpmCon *TPM) (bool, error) {
	var res bool
	var err error
	var flags tpm1.PermanentFlags
	switch tpmCon.Version {
	case TPMVersion12:
		flags, err = tpm1.GetPermanentFlags(tpmCon.RWC)
		if err != nil {
			return res, err
		}
		res = flags.NVLocked
		return res, nil
	case TPMVersion20:
		err = tpm2.HierarchyChangeAuth(tpmCon.RWC, tpm2.HandlePlatform, tpm2.AuthCommand{Session: tpm2.HandlePasswordSession, Attributes: tpm2.AttrContinueSession}, string(tpm2.EmptyAuth))
		if err == nil {
			return false, err
		}
		res = strings.Contains(err.Error(), tpm2LockedResult)
		if !res {
			return res, err
		}
		return res, nil
	}
	return false, fmt.Errorf("unknown TPM version: %v ", tpmCon.Version)
}

// ReadNVPublic reads public data about an NV index
func (h HwAPI) ReadNVPublic(tpmCon *TPM, index uint32) ([]byte, error) {
	return tpmCon.ReadNVPublic(index)
}

// NVReadValue reads a given NV index
func (h HwAPI) NVReadValue(tpmCon *TPM, index uint32, password string, size, offhandle uint32) ([]byte, error) {
	return tpmCon.NVReadValue(index, password, size, offhandle)
}

// ReadPCR read fom a given tpm connection a given pc register
func (h HwAPI) ReadPCR(tpmCon *TPM, pcr uint32) ([]byte, error) {
	return tpmCon.ReadPCR(pcr)
}
