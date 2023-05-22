package internal

import (
	"github.com/willabides/git-describe-semver/internal/testutil"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"
)

func TestGitTagMap(t *testing.T) {
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	require.NoError(t, err)

	tags, err := GitTagMap(*repo)
	require.NoError(t, err)
	require.Equal(t, map[string]string{}, *tags)

	commit1 := testutil.CreateEmptyCommit(t, repo, "first", nil)
	tag1 := testutil.MustCreateTag(t, repo, "v1.0.0", commit1)
	tags, err = GitTagMap(*repo)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		tag1.Hash().String(): "v1.0.0",
	}, *tags)

	commit2 := testutil.CreateEmptyCommit(t, repo, "second", nil)
	tag2 := testutil.MustCreateAnnotatedTag(t, repo, "v2.0.0", "Version 2.0.0", commit2)
	require.NotEqual(t, commit2.String(), tag2.Hash().String())
	tags, err = GitTagMap(*repo)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		commit1.String(): "v1.0.0",
		commit2.String(): "v2.0.0",
	}, *tags)

	commit3 := testutil.CreateEmptyCommit(t, repo, "third", nil)
	tag3 := testutil.MustCreateAnnotatedTag(t, repo, "fum", "Not a semver version tag", commit3)
	require.NotEqual(t, commit3.String(), tag3.Hash().String())
	tags, err = GitTagMap(*repo)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		commit1.String(): "v1.0.0",
		commit2.String(): "v2.0.0",
	}, *tags)
}

func TestGitDescribe(t *testing.T) {
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	require.NoError(t, err)
	_, _, _, err = GitDescribe(*repo)
	require.Error(t, err)
	test := func(expectedTagName string, expectedCounter int, expectedHeadHash string) {
		t.Helper()
		actualTagName, actualCounter, actualHeadHash, e := GitDescribe(*repo)
		require.NoError(t, e)
		require.Equal(t, expectedTagName, *actualTagName)
		require.Equal(t, expectedCounter, *actualCounter)
		require.Equal(t, expectedHeadHash, *actualHeadHash)
	}

	commit1 := testutil.CreateEmptyCommit(t, repo, "first", nil)
	test("", 1, commit1.String())

	testutil.MustCreateTag(t, repo, "v1.0.0", commit1)
	test("v1.0.0", 0, commit1.String())

	commit2 := testutil.CreateEmptyCommit(t, repo, "second", nil)
	test("v1.0.0", 1, commit2.String())

	commit3 := testutil.CreateEmptyCommit(t, repo, "third", nil)
	test("v1.0.0", 2, commit3.String())

	testutil.MustCreateTag(t, repo, "v2.0.0", commit3)
	test("v2.0.0", 0, commit3.String())
}

func TestGitDescribeWithBranch(t *testing.T) {
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	require.NoError(t, err)
	_, _, _, err = GitDescribe(*repo)
	require.Error(t, err)
	test := func(expectedTagName string, expectedCounter int, expectedHeadHash string) {
		actualTagName, actualCounter, actualHeadHash, e := GitDescribe(*repo)
		require.NoError(t, e)
		require.Equal(t, expectedTagName, *actualTagName)
		require.Equal(t, expectedCounter, *actualCounter)
		require.Equal(t, expectedHeadHash, *actualHeadHash)
	}

	commit1 := testutil.CreateEmptyCommit(t, repo, "first", nil)
	test("", 1, commit1.String())

	testutil.MustCreateTag(t, repo, "v1.0.0", commit1)
	test("v1.0.0", 0, commit1.String())

	commit2 := testutil.CreateEmptyCommit(t, repo, "second", nil)
	test("v1.0.0", 1, commit2.String())

	commit3 := testutil.CreateEmptyCommit(t, repo, "third", []plumbing.Hash{commit1})
	test("v1.0.0", 1, commit3.String())

	commit4 := testutil.CreateEmptyCommit(t, repo, "forth", []plumbing.Hash{commit2, commit3})
	test("v1.0.0", 2, commit4.String())

	testutil.MustCreateTag(t, repo, "v2.0.0", commit3)
	test("v2.0.0", 1, commit4.String())
}
