package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

type Gcov struct {
	CWD      string      `json:"current_working_directory"`
	DataFile string      `json:"data_file"`
	Version  string      `json:"format_version"`
	Files    []*GcovFile `json:"files"`
}

type GcovFunc struct {
	Blocks        uint   `json:"blocks"`
	ExecedBlocks  uint   `json:"blocks_executed"`
	DemangledName string `json:"demangled_name"`
	EndCol        uint   `json:"end_column"`
	EndLine       uint   `json:"end_line"`
	ExecCount     uint   `json:"execution_count"`
	Name          string `json:"name"`
	StartCol      uint   `json:"start_column"`
	StartLine     uint   `json:"start_line"`
}

type GcovLine struct {
	Branches      []*GcovBranch `json:"branches"`
	Count         uint          `json:"count"`
	LineNumber    uint          `json:"line_number"`
	UnexecedBlock bool          `json:"unexecuted_block"`
	FuncName      string        `json:"function_name"`

	// Extensions used only by symex-vp
	Tainted bool `json:"symex/tainted_instr"`
	Symbolic bool `json:"symex/symbolic_once"`

	// Not available in JSON, added separatly.
	SourceCode string
	NoCode     bool
}

type GcovBranch struct {
	Count       uint `json:"count"`
	Fallthrough bool `json:"fallthrough"`
	Throw       bool `json:"throw"`
}

func OpenGcov(fp string) (*Gcov, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var cov Gcov
	dec := json.NewDecoder(reader)
	err = dec.Decode(&cov)
	if err != nil {
		return nil, err
	}

	err = cov.addLines()
	if err != nil {
		return nil, err
	}

	return &cov, nil
}

func (g *Gcov) addLines() error {
	for _, file := range g.Files {
		cfp := filepath.Join(g.CWD, file.Name)
		f, err := os.Open(cfp)
		if err != nil {
			return err
		}

		lnum := uint(1)
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			code := scanner.Text()
			for _, line := range file.Lines {
				if line.LineNumber == lnum {
					line.SourceCode = code
					goto nextLine
				}
			}

			file.Lines = append(file.Lines, &GcovLine{
				Branches:      []*GcovBranch{},
				Count:         0,
				LineNumber:    lnum,
				UnexecedBlock: false,
				FuncName:      "", // XXX
				SourceCode:    code,
				NoCode:        true,
			})

		nextLine:
			lnum++
		}

		// Make sure lines are sorted by line number
		sort.Sort(byLine(file.Lines))

		f.Close()
		err = scanner.Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Gcov) LineCoverage() Coverage {
	var execLines uint
	var totalLines uint

	for _, file := range g.Files {
		totalLines += uint(len(file.TotalCodeLines()))
		for _, line := range file.Lines {
			if !line.UnexecedBlock && !line.NoCode {
				execLines++
			}
		}
	}

	p := float64(execLines) / float64(totalLines)
	return Coverage{
		Exec:       execLines,
		Total:      totalLines,
		Percentage: p * 100,
	}
}

func (g *Gcov) SymbolicCoverage() Coverage {
	var taintLines uint
	var totalLines uint

	for _, file := range g.Files {
		totalLines += uint(len(file.TotalCodeLines()))
		for _, line := range file.Lines {
			if line.Tainted && !line.NoCode {
				taintLines++
			}
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
