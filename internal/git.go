package internal

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GitTagMap(repo git.Repository) (map[string]string, error) {
	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	tagMap := map[string]string{}
	err = iter.ForEach(func(r *plumbing.Reference) error {
		if SemVerParse(r.Name().Short()) == nil {
			// Filter out tags that are not semver
			return nil
		}
		var tag *object.Tag
		tag, err = repo.TagObject(r.Hash())
		switch err {
		case nil:
			var c *object.Commit
			c, err = tag.Commit()
			if err != nil {
				return err
			}
			tagMap[c.Hash.String()] = r.Name().Short()
		case plumbing.ErrObjectNotFound:
			tagMap[r.Hash().String()] = r.Name().Short()
		default:
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tagMap, nil
}

func GitDescribe(repo git.Repository) (tagPtr *string, counterPtr *int, headPtr *string, _ error) {
	type gitDescribeNode struct {
		Commit   object.Commit
		Distance int
	}

	head, err := repo.Head()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to find head: %v", err)
	}
	headHash := head.Hash().String()
	tags, err := GitTagMap(repo)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to get tags: %v", err)
	}
	commits, err := repo.Log(&git.LogOptions{
		From:  head.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to get log: %v", err)
	}
	state := map[string]gitDescribeNode{}
	counter := 0
	tagHash := ""
	err = commits.ForEach(func(c *object.Commit) error {
		node, found := state[c.Hash.String()]
		if !found {
			node = gitDescribeNode{
				Commit:   *c,
				Distance: 0,
			}
			state[c.Hash.String()] = node
		}
		err = c.Parents().ForEach(func(p *object.Commit) error {
			_, foundIt := state[p.Hash.String()]
			if !foundIt {
				state[p.Hash.String()] = gitDescribeNode{
					Commit:   *p,
					Distance: node.Distance + 1,
				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		_, foundTag := tags[c.Hash.String()]
		if tagHash == "" && foundTag {
			counter = state[c.Hash.String()].Distance
			tagHash = c.Hash.String()
		}
		return nil
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to get log: %v", err)
	}
	if tagHash == "" {
		for i := range state {
			if state[i].Distance+1 > counter {
				counter = state[i].Distance + 1
			}
		}
		tagName := ""
		return &tagName, &counter, &headHash, nil
	}
	tagName := tags[tagHash]
	return &tagName, &counter, &headHash, nil
}
