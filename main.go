package main

import (
	"fmt"
	"log"
	"sort"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

func getBranches() (branches []*Branch) {
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, merged := range [2]bool{true, false} {
		wg.Add(1)
		go func(merged bool) {
			defer wg.Done()
			result := gitBranch(merged)
			mutex.Lock()
			branches = append(branches, result...)
			mutex.Unlock()
		}(merged)
	}
	wg.Wait()
	return
}

func sortBranches(branches []*Branch) {
	sort.Slice(branches, func(i, j int) bool {
		b1, b2 := branches[i], branches[j]
		if b1.Current != b2.Current {
			return b1.Current // current branch first
		} else if b1.Merged != b2.Merged {
			return b1.Merged // merged branches second
		}
		return b1.Name < b2.Name // sort alphabetically otherwise
	})
}

func launchInterface(options []*Branch) *model {
	m := initialModel(options)
	p := tea.NewProgram(&m)
	if err := p.Start(); err != nil {
		log.Fatalf("Error: %v", err)
	}
	return &m
}

func getSelectedBranchNames(model *model) (selected []string) {
	for i := range model.selected {
		selected = append(selected, model.options[i].Name)
	}
	return
}

func deleteBranches(selected []string) {
	if len(selected) > 0 {
		output := gitBranchDelete(selected)
		fmt.Printf("%s\nDeleted %d branches", output, len(selected))
	} else {
		fmt.Println("No branches selected")
	}
}

func main() {
	gitFetch()

	branches := getBranches()
	sortBranches(branches)

	model := launchInterface(branches)
	selected := getSelectedBranchNames(model)

	deleteBranches(selected)
}
