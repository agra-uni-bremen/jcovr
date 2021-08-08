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

func (f *GcovFile) TotalCodeLines() []*GcovLine {
	var totalCodeLines uint

	codeLines := make([]*GcovLine, len(f.Lines))
	for _, line := range f.Lines {
		if line.NoCode {
			continue
		}

		codeLines[totalCodeLines] = line
		totalCodeLines++
	}

	// Shrink to appropriate size
	codeLines = codeLines[0:totalCodeLines]

	return codeLines
}

func (f *GcovFile) LineCoverage() Coverage {
	codeLines := f.TotalCodeLines()
	totalLines := uint(len(codeLines))

	var execLines uint
	for _, line := range codeLines {
		if !line.UnexecedBlock {
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
	codeLines := f.TotalCodeLines()
	totalLines := uint(len(codeLines))

	var taintLines uint
	for _, line := range codeLines {
		if line.Tainted {
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
