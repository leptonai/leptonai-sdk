package util

import (
	"crypto/md5"
	"encoding/base32"
	"fmt"
	"regexp"
	"strings"

	leptonaiv1alpha1 "github.com/leptonai/lepton/lepton-deployment-operator/api/v1alpha1"
)

const (
	NameInvalidMessage = "Name must consist of lower case alphanumeric characters or '-', and must start with an alphabetical character and end with an alphanumeric character"
)

// HexHash returns the hex encoded md5 hash of the given text.
func HexHash(text []byte) string {
	hash := md5.Sum(text)
	return strings.ToLower(base32.HexEncoding.EncodeToString(hash[:])[0:8])
}

var (
	nameRegex = regexp.MustCompile("^[a-z]([-a-z0-9]*[a-z0-9])?$")
)

// ValidateName returns true if the given name is valid.
func ValidateName(name string) bool {
	return nameRegex.MatchString(name) && len(name) <= 32
}

func DomainName(ld *leptonaiv1alpha1.LeptonDeployment, rootDomain string) string {
	return fmt.Sprintf("%s.%s", ld.GetName(), rootDomain)
}
