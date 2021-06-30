package main

// byLine sorts RepoFiles by their object type (directories first).
type byLine []*GcovLine

func (t byLine) Len() int {
	return len(t)
}

func (t byLine) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t byLine) Less(i, j int) bool {
	return t[i].LineNumber < t[j].LineNumber
}
