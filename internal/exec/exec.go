package exec

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

type L struct {
	Name       string
	Type       string
	Extensions []string
}

func GitLog(rep, commit, fileName string) (string, string) {
	cmd := exec.Command("git", "log", "--pretty=format:%H %an", commit, "--", fileName)
	cmd.Dir = rep
	b, _ := cmd.Output()
	var s strings.Builder
	s.Write(b)
	split := strings.Split(s.String(), "\n")
	s1 := strings.Split(split[0], " ")
	return s1[0], strings.Join(s1[1:], " ")
}

func GitBlame(rep, commit, name string) []string {
	cmd := exec.Command("git", "blame", "--porcelain", commit, name)
	cmd.Dir = rep
	b, _ := cmd.Output()

	var s strings.Builder
	s.Write(b)
	return strings.FieldsFunc(s.String(), func(r rune) bool {
		return r == '\n'
	})
}

func GitLsTree(rep, commit string, extensions, languages []string) []string {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", commit)
	cmd.Dir = rep
	b, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var s strings.Builder
	s.Write(b)
	out := strings.FieldsFunc(s.String(), func(r rune) bool {
		return r == '\n'
	})

	ans := out

	if len(languages) != 0 {
		getLanguagesExtensions(languages, &extensions)
	}

	if len(extensions) != 0 {
		ans = filterExtension(ans, extensions)
	}

	return ans
}

func getLanguagesExtensions(languages []string, extensions *[]string) {
	b, _ := ioutil.ReadFile("../../configs/language_extensions.json")

	var l []L
	err := json.Unmarshal(b, &l)
	if err != nil {
		panic(err)
	}

	set := make(map[string]int, len(l))

	for i, j := range l {
		set[strings.ToLower(j.Name)] = i
	}

	for _, i := range languages {
		if itr, ok := set[strings.ToLower(i)]; ok {
			*extensions = append(*extensions, l[itr].Extensions...)
		}
	}
}

func filterExtension(out, extensions []string) []string {
	e := make(map[string]struct{}, len(extensions))
	for _, i := range extensions {
		e[i] = struct{}{}
	}

	ans := make([]string, 0)
	for _, i := range out {
		if _, ok := e[filepath.Ext(i)]; ok {
			ans = append(ans, i)
		}
	}

	return ans
}
