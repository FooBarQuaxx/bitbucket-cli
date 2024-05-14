package cli

import (
	"encoding/json"
	"fmt"
	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"regexp"
	"strings"
)

type RepoBranchListCmd struct {
	Base    string `arg:"--base" help:"base branch or tag to compare each branch to (for the metadata providers that uses that information)"`
	Details bool   `arg:"--details" help:"whether to retrieve plugin-provided metadata about each branch"`
	OrderBy string `arg:"--order-by" help:"ordering of refs either ALPHABETICAL (by name) or MODIFICATION (last updated)"`
	Filter  string `arg:"-f,--filter" help:"Filter to match branch names against (contains)"`
	Prefix  string `arg:"-p,--prefix" help:"Only list branches that start with this prefix"`
	Regex   string `arg:"-r,--regex" help:"Only list branches that start with this prefix"`
	Limit   int    `arg:"--limit" default:"-1"`
	Start   int    `arg:"--start" default:"-1"`
	Json    bool   `arg:"--json"`
}

func (b *BitbucketCLI) branchCmdList(cmd *RepoCmd) {
	if cmd == nil || cmd.BranchCmd == nil || cmd.BranchCmd.List == nil {
		return
	}

	filterFunction := func(branch bitbucketv1.Branch) bool { return true }

	list := cmd.BranchCmd.List
	if list.Prefix != "" {
		prevFilter := filterFunction
		filterFunction = func(branch bitbucketv1.Branch) bool {
			return prevFilter(branch) && strings.HasPrefix(branch.DisplayID, list.Prefix)
		}
	}

	if list.Regex != "" {
		regex, err := regexp.Compile(list.Regex)

		if err != nil {
			b.logger.Warnf("Regex %s is not valid, will be skipped", list.Regex)
		} else {
			prevFilter := filterFunction
			filterFunction = func(branch bitbucketv1.Branch) bool {
				return prevFilter(branch) && regex.MatchString(branch.DisplayID)
			}
		}
	}
	optionals := make(map[string]interface{})
	if list.Filter != "" {
		optionals["filterText"] = list.Filter
	}
	if list.Base != "" {
		optionals["base"] = list.Base
	}
	optionals["details"] = list.Details
	if list.OrderBy != "" {
		optionals["orderBy"] = list.OrderBy
	}
	if list.Limit != -1 {
		optionals["limit"] = list.Limit
	}
	if list.Start != -1 {
		optionals["start"] = list.Start
	}
	response, err := b.client.DefaultApi.GetBranches(cmd.ProjectKey, cmd.Slug, optionals)

	if err != nil {
		b.logger.Fatalf("Failed to fetch branches %s", err.Error())
		return
	}
	branches, err := bitbucketv1.GetBranchesResponse(response)

	if err != nil {
		b.logger.Fatalf("Failed to parse branches response %s", err.Error())
		return
	}

	filterdBranches := []bitbucketv1.Branch{}
	for _, b := range branches {
		if filterFunction(b) {
			filterdBranches = append(filterdBranches, b)
		}
	}

	if list.Json {
		jsonOutput, err := json.Marshal(filterdBranches)
		if err != nil {
			b.logger.Fatalf("unable to marshal JSON: %v", err)
		}
		fmt.Printf("%s", jsonOutput)

	} else {
		for _, branch := range filterdBranches {
			fmt.Printf("%s \n", branch.DisplayID)
		}
	}

}
