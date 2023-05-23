package internal

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type SemVer struct {
	Prefix string
	semver.Version
}

func (v *SemVer) String() string {
	return fmt.Sprintf("%s%s", v.Prefix, v.Version.String())
}

func (v *SemVer) SetPrerelease(prerelease string) error {
	ver, err := v.Version.SetPrerelease(prerelease)
	if err != nil {
		return err
	}
	v.Version = ver
	return nil
}

// NextPatch is like IncPatch but preserves the prefix and metadata
func (v *SemVer) NextPatch() {
	metadata := v.Metadata()
	prerelease := v.Prerelease()
	v.Version = v.Version.IncPatch()
	vv, err := v.Version.SetMetadata(metadata)
	if err != nil {
		panic(err)
	}
	v.Version = vv
	err = v.SetPrerelease(prerelease)
	if err != nil {
		panic(err)
	}
}

func SemVerParse(str string) *SemVer {
	// prefix cannot contain a digit so the end of the prefix is the first digit
	prefixEnd := strings.IndexFunc(str, func(r rune) bool {
		return r >= '0' && r <= '9'
	})
	if prefixEnd == -1 {
		return nil
	}
	sv, err := semver.StrictNewVersion(str[prefixEnd:])
	if err != nil {
		return nil
	}

	return &SemVer{
		Prefix:  str[:prefixEnd],
		Version: *sv,
	}
}
