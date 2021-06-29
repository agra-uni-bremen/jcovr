package main

import (
	"compress/gzip"
	"encoding/json"
	"os"
)

type Gcov struct {
	CWD      string      `json:"current_working_directory"`
	DataFile string      `json:"data_file"`
	Version  string      `json:"format_version"`
	Files    []GcovFiles `json:"files"`
}

type GcovFiles struct {
	Name  string     `json:"file"`
	Funcs []GcovFunc `json:"functions"`
	Lines []GcovLine `json:"lines"`
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
	Branches      []GcovBranch `json:"branches"`
	Count         uint         `json:"count"`
	LineNumber    uint         `json:"line_number"`
	UnexecedBlock bool         `json:"unexecuted_block"`
	FuncName      string       `json:"function_name"`
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

	return &cov, nil
}
