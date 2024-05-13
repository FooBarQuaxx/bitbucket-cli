package cli

import (
	"encoding/json"
	"fmt"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
)

type CommitListCmd struct {
	FollowRenames bool   `arg:"--follow-renames"`
	IgnoreMissing bool   `arg:"--ignore-missing"`
	Merges        string `arg:"--merges"`
	Path          string `arg:"--path"`
	Since         string `arg:"--since"`
	Until         string `arg:"--until"`
	WithCounts    bool   `arg:"--with-counts"`
	Start         int    `arg:"--start" default:"-1"`
	Limit         int    `arg:"--limit" default:"-1"`

	Json bool `arg:"--json"`
}

func (b *BitbucketCLI) commitListCmd(cmd *RepoCmd) {
	if cmd.CommitCmd.List == nil {
		return
	}

	lCmd := cmd.CommitCmd.List

	var err error

	options := map[string]interface{}{}

	options["followRenames"] = lCmd.FollowRenames

	options["IgnoreMissing"] = lCmd.IgnoreMissing

	options["withCounts"] = lCmd.WithCounts

	if lCmd.Path != "" {
		options["path"] = lCmd.Path
	}

	if lCmd.Since != "" {
		options["since"] = lCmd.Since
	}

	if lCmd.Until != "" {
		options["until"] = lCmd.Until
	}

	if lCmd.Start != -1 {
		options["start"] = lCmd.Start
	}

	if lCmd.Limit != -1 {
		options["limit"] = lCmd.Limit
	}

	res, err := b.client.DefaultApi.GetCommits(cmd.ProjectKey, cmd.Slug, options)
	if err != nil {
		b.logger.Fatalf("unable to list commits: %v", err)
		return
	}

	commits, err := bitbucket.GetCommitsResponse(res)
	if err != nil {
		b.logger.Fatalf("unable to parse commits list: %v", err)
		return
	}

	if lCmd.Json {
		jsonOutput, err := json.Marshal(commits)
		if err != nil {
			b.logger.Fatalf("unable to marshal JSON: %v", err)
		}
		fmt.Printf("%s", jsonOutput)
	} else {
		for _, c := range commits {
			fmt.Printf("%s\n", c.ID)
		}
	}
}
