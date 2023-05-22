package internal

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"strings"
)

// SemVer ...
type SemVer struct {
	Prefix string
	semver.Version
}

// Equal ...
func (v *SemVer) Equal(v2 SemVer) bool {
	return v.Version.Equal(&v2.Version) &&
		v.Prefix == v2.Prefix &&
		v.Metadata() == v2.Metadata()
}

// String ...
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

func (v *SemVer) SetMetadata(metadata string) error {
	ver, err := v.Version.SetMetadata(metadata)
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
	err := v.SetMetadata(metadata)
	if err != nil {
		panic(err)
	}
	err = v.SetPrerelease(prerelease)
	if err != nil {
		panic(err)
	}
}

// SemVerParse ...
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
