package cmd

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"
	"github.com/willabides/git-describe-semver/internal"
	"github.com/willabides/git-describe-semver/internal/testutil"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()
	_, err := run(dir, internal.GenerateVersionOptions{PrereleasePrefix: "dev"})
	require.Error(t, err)

	repo, err := git.PlainInit(dir, false)
	require.NoError(t, err)
	_, err = run(dir, internal.GenerateVersionOptions{PrereleasePrefix: "dev"})
	require.Error(t, err)

	commit1 := testutil.CreateEmptyCommit(t, repo, "first", nil)
	testutil.MustCreateTag(t, repo, "invalid", commit1)
	got, err := run(dir, internal.GenerateVersionOptions{PrereleasePrefix: "dev"})
	require.Error(t, err)
	require.Nil(t, got)

	commit2 := testutil.CreateEmptyCommit(t, repo, "first", []plumbing.Hash{commit1})
	testutil.MustCreateTag(t, repo, "v1.0.0", commit2)

	commit3 := testutil.CreateEmptyCommit(t, repo, "second", []plumbing.Hash{commit2})
	result, err := run(dir, internal.GenerateVersionOptions{PrereleasePrefix: "dev"})
	require.NoError(t, err)
	require.Equal(t, "v1.0.1-dev.1.g"+commit3.String()[0:7], *result)
}
