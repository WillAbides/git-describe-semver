package internal

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
)

func TestSemVerString(t *testing.T) {
	for _, s := range []string{
		"0.0.0",
		"v0.0.0",
		"1.2.3",
		"0.0.0-rc.1",
		"0.0.0-alpha-version.1",
		"0.0.0+foo.bar",
		"v1.2.3-rc.1+foo.bar",
	} {
		t.Run(s, func(t *testing.T) {
			require.Equal(t, s, SemVerParse(s).String())
		})
	}
}

func TestSemVerParse(t *testing.T) {
	test := func(input string, expected *SemVer) {
		t.Helper()
		actual := SemVerParse(input)
		require.Equal(t, expected, actual)
	}

	test("0.0.0", &SemVer{Version: *mustMMNewVersion(t, "0.0.0")})
	test("v0.0.0", &SemVer{Prefix: "v", Version: *mustMMNewVersion(t, "0.0.0")})
	test("1.2.3", &SemVer{Version: *mustMMNewVersion(t, "1.2.3")})
	test("0.0.0-rc.1", &SemVer{Version: *mustMMNewVersion(t, "0.0.0-rc.1")})
	test("0.0.0-alpha-version.1", &SemVer{Version: *mustMMNewVersion(t, "0.0.0-alpha-version.1")})
	test("0.0.0+foo.bar", &SemVer{Version: *mustMMNewVersion(t, "0.0.0+foo.bar")})
	test("v1.2.3-rc.1+foo.bar", &SemVer{Prefix: "v", Version: *mustMMNewVersion(t, "1.2.3-rc.1+foo.bar")})
	test("invalid", nil)
}

func mustMMNewVersion(t *testing.T, s string) *semver.Version {
	t.Helper()
	v, err := semver.StrictNewVersion(s)
	require.NoError(t, err)
	return v
}
