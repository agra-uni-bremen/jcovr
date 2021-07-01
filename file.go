package main

type GcovFile struct {
	Name  string      `json:"file"`
	Funcs []*GcovFunc `json:"functions"`
	Lines []*GcovLine `json:"lines"`
}

type LineCoverage struct {
	Exec       uint
	Total      uint
	Percentage float64
}

func (f *GcovFile) LineCoverage() LineCoverage {
	var execLines uint
	totalLines := uint(len(f.Lines))

	for _, line := range f.Lines {
		if !line.UnexecedBlock {
			execLines++
		}
	}

	p := float64(execLines) / float64(totalLines)
	return LineCoverage{
		Exec:       execLines,
		Total:      totalLines,
		Percentage: p * 100,
	}
}
