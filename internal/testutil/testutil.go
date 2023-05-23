package testutil

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/require"
)

func MustCreateTag(t *testing.T, repo *git.Repository, name string, target plumbing.Hash) *plumbing.Reference {
	t.Helper()
	ref, err := repo.CreateTag(name, target, nil)
	require.NoError(t, err)
	return ref
}

func MustCreateAnnotatedTag(t *testing.T, repo *git.Repository, name, message string, target plumbing.Hash) *plumbing.Reference {
	t.Helper()
	ref, err := repo.CreateTag(name, target, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "Foo Bar",
			Email: "foo@bar.com",
		},
		Message: message,
	})
	require.NoError(t, err)
	return ref
}

func CreateEmptyCommit(t *testing.T, repo *git.Repository, msg string, parents []plumbing.Hash) plumbing.Hash {
	t.Helper()
	worktree, err := repo.Worktree()
	require.NoError(t, err)
	author := object.Signature{Name: "Test", Email: "test@test.com"}
	c, err := worktree.Commit(msg, &git.CommitOptions{
		Author:            &author,
		AllowEmptyCommits: true,
		Parents:           parents,
	})
	require.NoError(t, err)
	return c
}
