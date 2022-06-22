package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const mergedOption = "--merged"
const noMergedOption = "--no-merged"
const remote = "origin"
const defaultBranch = "master"

type Branch struct {
	Name    string
	Merged  bool
	Current bool
}

func newBranch(name string, merged bool) *Branch {
	return &Branch{
		Name:    strings.ReplaceAll(name, "*", ""),
		Merged:  merged,
		Current: isCurrentBranch(name),
	}
}

func (b *Branch) print() string {
	return fmt.Sprintf("%s (merged: %t)", b.Name, b.Merged)
}

func gitBranch(merged bool) (branches []*Branch) {
	var option string

	if merged {
		option = mergedOption
	} else {
		option = noMergedOption
	}

	commit := fmt.Sprintf("%s/%s", remote, defaultBranch)

	cmd := exec.Command("git", "branch", option, commit)

	output, err := cmd.Output()

	if err != nil {
		log.Fatalf("Error running git branch %s: %v", option, err)
	}

	trimmedString := strings.Trim(string(output), "\n")
	trimmedString = strings.ReplaceAll(trimmedString, " ", "")
	branchNames := strings.Split(trimmedString, "\n")

	for _, name := range branchNames {
		branches = append(branches, newBranch(name, merged))
	}

	return
}

func isCurrentBranch(name string) bool {
	re := regexp.MustCompile(`^\*\w+`)
	return re.Match([]byte(name))
}

func gitBranchDelete(branches []string) {
	args := []string{"branch", "-D"}
	args = append(args, branches...)
	cmd := exec.Command("git", args...)

	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error while deleting branches: %v", err)
	}
	fmt.Printf("%s", string(output))
}

func gitFetch() {
	cmd := exec.Command("git", "fetch", remote)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error while fetching %s: %v", remote, err)
	}
}

func main() {
	gitFetch()
	var branches []*Branch
	mergedBranches := gitBranch(true)
	unmergedBranches := gitBranch(false)
	branches = append(branches, mergedBranches...)
	branches = append(branches, unmergedBranches...)
	sort.Slice(branches, func(i, j int) bool {
		b1, b2 := branches[i], branches[j]
		if b1.Current {
			return true
		}
		return b1.Name < b2.Name
	})

	var options []string
	for _, b := range branches {
		options = append(options, b.print())
	}

	m := initialModel(options)
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	var selected []string

	for i := range m.selected {
		selected = append(selected, branches[i].Name)
	}

	if len(selected) > 0 {
		gitBranchDelete(selected)
		fmt.Printf("Deleted %d branches", len(selected))
	} else {
		fmt.Println("No branches selected")
	}
}
