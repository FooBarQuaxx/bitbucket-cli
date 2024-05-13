package cli

type CommitCmd struct {
	List *CommitListCmd `arg:"subcommand:list"`
}

func (b *BitbucketCLI) commitCmd(cmd *RepoCmd) {
	if cmd.CommitCmd.List != nil {
		b.commitListCmd(cmd)
		return
	}
	b.logger.Fatal(errSpecifySubcommand) // a constant.
}
