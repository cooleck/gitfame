package parse

import (
	"gitlab.com/slon/shad-go/gitfame/internal/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Statistics struct {
	Lines   int
	Commits int
	Files   int
}

type Author struct {
	Statistics Statistics
	Name       string
}

type AuthorSlice struct {
	Slice   []Author
	orderBy string
}

func (a AuthorSlice) Len() int {
	return len(a.Slice)
}

func (a AuthorSlice) Swap(i, j int) {
	a.Slice[i], a.Slice[j] = a.Slice[j], a.Slice[i]
}

func (a AuthorSlice) Less(i, j int) bool {
	var key1, key2 []int

	switch a.orderBy {
	case "lines":
		key1 = []int{a.Slice[i].Statistics.Lines, a.Slice[i].Statistics.Commits, a.Slice[i].Statistics.Files}
		key2 = []int{a.Slice[j].Statistics.Lines, a.Slice[j].Statistics.Commits, a.Slice[j].Statistics.Files}
	case "commits":
		key1 = []int{a.Slice[i].Statistics.Commits, a.Slice[i].Statistics.Lines, a.Slice[i].Statistics.Files}
		key2 = []int{a.Slice[j].Statistics.Commits, a.Slice[j].Statistics.Lines, a.Slice[j].Statistics.Files}
	case "files":
		key1 = []int{a.Slice[i].Statistics.Files, a.Slice[i].Statistics.Lines, a.Slice[i].Statistics.Commits}
		key2 = []int{a.Slice[j].Statistics.Files, a.Slice[j].Statistics.Lines, a.Slice[j].Statistics.Commits}
	default:
		panic("order-by flag value is incorrect")
	}

	if key1[0] == key2[0] {
		if key1[1] == key2[1] {
			if key1[2] == key2[2] {
				return strings.ToLower(a.Slice[i].Name) < strings.ToLower(a.Slice[j].Name)
			}
			return key1[2] > key2[2]
		}
		return key1[1] > key2[1]
	}
	return key1[0] > key2[0]
}

func RestrictTo(files, restrictTo []string) []string {
	ans := make([]string, 0)
	for _, i := range files {
		check := false
		for _, j := range restrictTo {
			if b, _ := filepath.Match(j, i); b {
				check = true
				break
			}
		}

		if check {
			ans = append(ans, i)
		}
	}

	return ans
}

func Exclude(files, exclude []string) []string {
	ans := make([]string, 0)
	for _, i := range files {
		check := true
		for _, j := range exclude {
			if b, _ := filepath.Match(j, i); b {
				check = false
				break
			}
		}

		if check {
			ans = append(ans, i)
		}
	}

	return ans
}

func Blame(out []string, useCommitter bool, rep, commit, fileName string) (map[string][]string, map[string]int) {
	authors := make(map[string][]string)
	commits := make(map[string]int)

	if len(out) == 0 {
		hash, a := exec.GitLog(rep, commit, fileName)
		commits[hash] = 0
		authors[a] = append(authors[a], hash)
	}

	isNextHash := true
	itr := 0
	var isWaitForAuthor bool
	var lastHash string
	for _, i := range out {
		if isNextHash {
			isNextHash = false
			if itr == 0 {
				s := strings.Split(i, " ")
				itr, _ = strconv.Atoi(s[len(s)-1])
				commits[s[0]] += itr
				isWaitForAuthor = true
				lastHash = s[0]
			}
			itr--
		} else if i[0] != '\t' && isWaitForAuthor {
			s := strings.Split(i, " ")
			if !useCommitter && s[0] == "author" {
				name := i[len("author "):]
				authors[name] = append(authors[name], lastHash)
				isWaitForAuthor = false
			} else if useCommitter && s[0] == "committer" {
				name := i[len("committer "):]
				authors[name] = append(authors[name], lastHash)
				isWaitForAuthor = false
			}
		} else if i[0] == '\t' {
			isNextHash = true
		}
	}

	return authors, commits
}

func AuthorData(authors map[string]*Statistics, a map[string]map[string]struct{}, c, files map[string]int) {
	for i, j := range a {
		authors[i] = &Statistics{
			Lines:   0,
			Commits: len(j),
			Files:   files[i],
		}
	}
	for i, j := range a {
		lines := 0
		for k := range j {
			lines += c[k]
		}
		authors[i].Lines += lines
	}
}

func Sort(authors map[string]*Statistics, orderBy string) AuthorSlice {
	var authorSlice AuthorSlice
	authorSlice.orderBy = orderBy
	for i, j := range authors {
		authorSlice.Slice = append(authorSlice.Slice, Author{
			Statistics: *j,
			Name:       i,
		})
	}

	sort.Sort(authorSlice)
	return authorSlice
}
