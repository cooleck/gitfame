package flagset

import (
	flag "github.com/spf13/pflag"
	"os"
)

func NewFlagSet() *flag.FlagSet {
	set := flag.NewFlagSet("flag_set", flag.ExitOnError)
	set.String("repository", ".", "path to Git repo")
	set.String("revision", "HEAD", "pointer on commit")
	set.String("order-by", "lines", "key for result sorting")
	set.Bool("use-committer", false, "flag to switch between author and committer")
	set.String("format", "tabular", "format of the output")
	set.StringSlice("extensions", []string{}, "list of extensions")
	set.StringSlice("languages", []string{}, "list of permitted languages")
	set.StringSlice("exclude", []string{}, "set of glob patterns ")
	set.StringSlice("restrict-to", []string{}, "set of glob patterns for restrict")

	err := set.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	return set
}
