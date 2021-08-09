package main

type GcovFile struct {
	Name  string      `json:"file"`
	Funcs []*GcovFunc `json:"functions"`
	Lines []*GcovLine `json:"lines"`

	Path string
}

type Coverage struct {
	Exec       uint
	Total      uint
	Percentage float64
}

func (f *GcovFile) TotalCodeLines() uint {
	var totalCodeLines uint

	for _, line := range f.Lines {
		if !line.NoCode {
			totalCodeLines++
		}
	}

	return totalCodeLines
}

func (f *GcovFile) LineCoverage() Coverage {
	var execLines uint
	totalLines := f.TotalCodeLines()

	for _, line := range f.Lines {
		if !line.UnexecedBlock && !line.NoCode {
			execLines++
		}
	}

	p := float64(execLines) / float64(totalLines)
	return Coverage{
		Exec:       execLines,
		Total:      totalLines,
		Percentage: p * 100,
	}
}

func (f *GcovFile) SymbolicCoverage() Coverage {
	var taintLines uint
	totalLines := f.TotalCodeLines()

	for _, line := range f.Lines {
		if line.Tainted && !line.NoCode {
			taintLines++
		}
	}

	symLines := totalLines - taintLines
	p := float64(symLines) / float64(totalLines)

	return Coverage{
		Exec:       symLines,
		Total:      totalLines,
		Percentage: p * 100,
	}
}
