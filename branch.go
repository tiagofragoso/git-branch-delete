package main

type Branch struct {
	Name    string
	Merged  bool
	Current bool
}

func newBranch(name string, merged bool, current bool) *Branch {
	return &Branch{
		Name:    name,
		Merged:  merged,
		Current: current,
	}
}

func (b *Branch) print() string {
	suffix := ""
	if b.Current {
		suffix = " (current)"
	} else if b.Merged {
		suffix = " (merged)"
	}
	return b.Name + suffix
}
