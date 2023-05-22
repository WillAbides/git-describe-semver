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

func TestSemVerEqual(t *testing.T) {
	test := func(a, b string, expected bool) {
		t.Helper()
		pa := *SemVerParse(a)
		pb := *SemVerParse(b)
		actual := pa.Equal(pb)
		require.Equal(t, expected, actual)
	}

	test("0.0.0", "0.0.0", true)
	test("1.0.0", "2.0.0", false)
	test("0.1.0", "0.2.0", false)
	test("0.0.1", "0.0.2", false)
	test("0.0.0-foo", "0.0.0-foo", true)
	test("0.0.0-foo", "0.0.0-bar", false)
	test("0.0.0-foo", "0.0.0", false)
	test("0.0.0", "0.0.0-bar", false)
	test("0.0.0+foo", "0.0.0+foo", true)
	test("0.0.0+foo", "0.0.0+bar", false)
	test("0.0.0+foo", "0.0.0", false)
	test("0.0.0", "0.0.0+bar", false)
}

func mustMMNewVersion(t *testing.T, s string) *semver.Version {
	t.Helper()
	v, err := semver.StrictNewVersion(s)
	require.NoError(t, err)
	return v
}
