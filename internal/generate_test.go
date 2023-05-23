package internal

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateVersion(t *testing.T) {
	now, err := time.Parse(time.RFC822Z, "01 Jan 20 02:03 -0000")
	require.NoError(t, err)
	for _, td := range []struct {
		inputTagName string
		inputCounter int
		inputOpts    GenerateVersionOptions
		expected     string
		err          bool
	}{
		{inputTagName: "0.0.0", inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev"}, expected: "0.0.0"},
		{inputTagName: "0.0.0", inputCounter: 1, inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev"}, expected: "0.0.1-dev.1.gabc1234"},
		{inputTagName: "0.0.0-rc1", inputCounter: 1, inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev"}, expected: "0.0.0-rc1.dev.1.gabc1234"},
		{inputTagName: "0.0.0-rc.1", inputCounter: 1, inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev"}, expected: "0.0.0-rc.1.dev.1.gabc1234"},
		{inputTagName: "0.0.0-rc.1+foobar", inputCounter: 1, inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev"}, expected: "0.0.0-rc.1.dev.1.gabc1234+foobar"},
		{inputTagName: "v0.0.0-rc.1+foobar", inputCounter: 1, inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev"}, expected: "v0.0.0-rc.1.dev.1.gabc1234+foobar"},
		{inputTagName: "", inputCounter: 1, inputOpts: GenerateVersionOptions{FallbackTagName: "0.0.0", PrereleasePrefix: "dev"}, expected: "0.0.0-dev.1.gabc1234"},
		{inputTagName: "", inputCounter: 1, inputOpts: GenerateVersionOptions{FallbackTagName: "v0.0.0", PrereleasePrefix: "dev"}, expected: "v0.0.0-dev.1.gabc1234"},
		{inputTagName: "v0.0.0", inputCounter: 0, inputOpts: GenerateVersionOptions{PrereleaseSuffix: "SNAPSHOT", PrereleasePrefix: "dev"}, expected: "v0.0.0"},
		{inputTagName: "v0.0.0", inputCounter: 1, inputOpts: GenerateVersionOptions{PrereleaseSuffix: "SNAPSHOT", PrereleasePrefix: "dev"}, expected: "v0.0.1-dev.1.gabc1234-SNAPSHOT"},
		{inputTagName: "v0.0.0", inputCounter: 0, inputOpts: GenerateVersionOptions{DropTagNamePrefix: true, PrereleasePrefix: "dev"}, expected: "0.0.0"},
		{inputTagName: "v0.0.0-rc.1", inputCounter: 1, inputOpts: GenerateVersionOptions{DropTagNamePrefix: true, PrereleasePrefix: "dev"}, expected: "0.0.0-rc.1.dev.1.gabc1234"},
		{inputTagName: "v0.0.0-rc.1+foobar", inputCounter: 1, inputOpts: GenerateVersionOptions{DropTagNamePrefix: true, PrereleasePrefix: "dev"}, expected: "0.0.0-rc.1.dev.1.gabc1234+foobar"},
		{inputTagName: "", inputCounter: 1, inputOpts: GenerateVersionOptions{FallbackTagName: "v0.0.0", DropTagNamePrefix: true, PrereleasePrefix: "dev"}, expected: "0.0.0-dev.1.gabc1234"},
		{inputTagName: "0.0.0", inputCounter: 0, inputOpts: GenerateVersionOptions{PrereleasePrefix: "custom"}, expected: "0.0.0"},
		{inputTagName: "0.0.0", inputCounter: 1, inputOpts: GenerateVersionOptions{PrereleasePrefix: "custom"}, expected: "0.0.1-custom.1.gabc1234"},
		{inputTagName: "0.0.0", inputCounter: 0, inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev", PrereleaseTimestamped: true}, expected: "0.0.0"},
		{inputTagName: "0.0.0", inputCounter: 1, inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev", PrereleaseTimestamped: false}, expected: "0.0.1-dev.1.gabc1234"},
		{inputTagName: "0.0.0", inputCounter: 1, inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev", PrereleaseTimestamped: true}, expected: "0.0.1-dev.1577844180.gabc1234"},
		{inputTagName: "", inputCounter: 1, inputOpts: GenerateVersionOptions{FallbackTagName: "0.0.0", PrereleasePrefix: "dev", PrereleaseTimestamped: true}, expected: "0.0.0-dev.1577844180.gabc1234"},
		{inputTagName: "0.0.0", inputCounter: 0, inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev", PrereleaseTimestamped: true}, expected: "0.0.0"},
		{inputTagName: "0.0.0", inputCounter: 1, inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev", PrereleaseTimestamped: true}, expected: "0.0.1-dev.1577844180.gabc1234"},
		{inputTagName: "", inputCounter: 1, inputOpts: GenerateVersionOptions{PrereleasePrefix: "dev"}, err: true},
	} {
		t.Run(fmt.Sprintf("%s-%d", td.inputTagName, td.inputCounter), func(t *testing.T) {
			actual, err := GenerateVersion(td.inputTagName, td.inputCounter, "abc1234", now, &td.inputOpts)
			if td.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, td.expected, *actual)
			}
		})
	}
}
