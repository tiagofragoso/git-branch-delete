package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

const mergedOption = "--merged"
const noMergedOption = "--no-merged"
const remote = "origin"
const branch = "master"

func gitBranch(merged bool) (branches []*Branch) {
	option := noMergedOption

	if merged {
		option = mergedOption
	}

	commit := fmt.Sprintf("%s/%s", remote, branch)

	cmd := exec.Command("git", "branch", option, commit)
	output, err := cmd.Output()

	if err != nil {
		log.Fatalf("Error running git branch %s: %v", option, err)
	}

	trimmedString := strings.ReplaceAll(string(output), " ", "")
	branchNames := strings.Split(trimmedString, "\n")

	for _, name := range branchNames {
		if name == "" {
			continue
		}

		current := isCurrentBranch(name)
		strippedName := strings.ReplaceAll(name, "*", "")
		branches = append(branches, newBranch(strippedName, merged, current))
	}

	return
}

func gitBranchDelete(branches []string) string {
	args := []string{"branch", "-D"}
	args = append(args, branches...)
	cmd := exec.Command("git", args...)

	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error while deleting branches: %v", err)
	}

	return string(output)
}

func gitFetch() {
	cmd := exec.Command("git", "fetch", remote)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error while fetching %s: %v", remote, err)
	}
}

func isCurrentBranch(name string) bool {
	re := regexp.MustCompile(`^\*\w+`)
	return re.Match([]byte(name))
}
