package main

import "flag"

var remote string
var branch string

func parseFlags() {
	flag.StringVar(&remote, "remote", "origin", "Remote to fetch from")
	flag.StringVar(&branch, "branch", "master", "Branch to compare against")
	flag.Parse()
}
