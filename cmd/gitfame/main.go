// +build !solution

package main

import (
	"gitlab.com/slon/shad-go/gitfame/internal/exec"
	"gitlab.com/slon/shad-go/gitfame/internal/flagset"
	"gitlab.com/slon/shad-go/gitfame/internal/format"
	"gitlab.com/slon/shad-go/gitfame/internal/parse"
)

func main() {
	set := flagset.NewFlagSet()

	rep, _ := set.GetString("repository")
	commit, _ := set.GetString("revision")
	extensions, _ := set.GetStringSlice("extensions")
	languages, _ := set.GetStringSlice("languages")
	files := exec.GitLsTree(rep, commit, extensions, languages)

	exclude, _ := set.GetStringSlice("exclude")
	if len(exclude) != 0 {
		files = parse.Exclude(files, exclude)
	}

	restrictTo, _ := set.GetStringSlice("restrict-to")
	if len(restrictTo) != 0 {
		files = parse.RestrictTo(files, restrictTo)
	}

	useCommitter, _ := set.GetBool("use-committer")

	authors := make(map[string]*parse.Statistics)
	a1 := make(map[string]map[string]struct{}, 100)
	c1 := make(map[string]int)
	f1 := make(map[string]int)
	for _, i := range files {
		s := exec.GitBlame(rep, commit, i)
		a, c := parse.Blame(s, useCommitter, rep, commit, i)
		UpdateData(a1, c1, f1, a, c)
	}

	parse.AuthorData(authors, a1, c1, f1)

	orderBy, _ := set.GetString("order-by")
	sortSLice := parse.Sort(authors, orderBy)

	formatString, _ := set.GetString("format")
	format.Format(sortSLice, formatString)
}

func UpdateData(a1 map[string]map[string]struct{}, c1, f1 map[string]int, a map[string][]string, c map[string]int) {
	for i, j := range c {
		c1[i] += j
	}

	for i, j := range a {
		f1[i]++
		for _, k := range j {
			if _, ok := a1[i]; !ok {
				a1[i] = make(map[string]struct{})
			}
			a1[i][k] = struct{}{}
		}
	}
}
