package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type GenerateVersionOptions struct {
	FallbackTagName       string
	DropTagNamePrefix     bool
	PrereleaseSuffix      string
	PrereleasePrefix      string
	PrereleaseTimestamped bool
	Format                string
}

func GenerateVersion(tagName string, counter int, headHash string, timestamp time.Time, opts *GenerateVersionOptions) (*string, error) {
	devPrerelease := strings.Join([]string{opts.PrereleasePrefix, strconv.Itoa(counter), "g" + (headHash)[0:7]}, ".")
	if opts.PrereleaseTimestamped {
		timestampUTC := timestamp.UTC()
		timestampSegments := []string{
			strconv.FormatInt(timestampUTC.UnixMilli()/1000, 10),
		}
		devPrerelease = strings.Join([]string{opts.PrereleasePrefix, strings.Join(timestampSegments, ""), "g" + (headHash)[0:7]}, ".")
	}
	if opts.PrereleaseSuffix != "" {
		devPrerelease += "-" + opts.PrereleaseSuffix
	}
	version := &SemVer{}
	if tagName == "" {
		version = SemVerParse(opts.FallbackTagName)
		if version == nil {
			return nil, fmt.Errorf("unable to parse fallback tag")
		}
		err := version.SetPrerelease(devPrerelease)
		if err != nil {
			return nil, err
		}
	} else {
		version = SemVerParse(tagName)
		if version == nil {
			return nil, fmt.Errorf("unable to parse tag")
		}
		if counter > 0 {
			if version.Prerelease() != "" {
				devPrerelease = version.Prerelease() + "." + devPrerelease
			} else {
				version.NextPatch()
			}
			err := version.SetPrerelease(devPrerelease)
			if err != nil {
				return nil, err
			}
		}
	}
	if opts.DropTagNamePrefix {
		version.Prefix = ""
	}
	result := version.String()
	if opts.Format != "" {
		result = strings.ReplaceAll(opts.Format, "<version>", result)
	}
	return &result, nil
}
