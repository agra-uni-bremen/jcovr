package main

import (
	"embed"
	"flag"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

var (
	dest  = flag.String("d", "./www", "output directory for HTML files")
	relfp = flag.String("r", "", "treat all paths relative to this directory")
)

//go:embed tmpl
var templates embed.FS

type CoveragePair struct {
	Coverage *Gcov
	File     *GcovFile
}

func createCSS(path string) error {
	const name = "base.tmpl"
	stylesheet := template.New(name)

	t, err := stylesheet.ParseFS(templates, "tmpl/css/*.tmpl")
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = t.Execute(file, nil)
	if err != nil {
		return err
	}

	return nil
}

func buildFiles(files []string) (*Gcov, error) {
	tmpl := template.New("source.tmpl")
	funcMap := template.FuncMap{
		"relIndex": relIndex,
	}
	tmpl = tmpl.Funcs(funcMap)

	tmpl, err := tmpl.ParseFS(templates, "tmpl/source.tmpl")
	if err != nil {
		return nil, err
	}

	var accCov Gcov
	accCov.CWD = *relfp
	if accCov.CWD == "" {
		accCov.CWD = filepath.Dir(files[0])
	}

	for _, file := range files {
		c, err := OpenGcov(file)
		if err != nil {
			return nil, err
		}

		for _, f := range c.Files {
			fp, err := filepath.Rel(accCov.CWD, filepath.Join(c.CWD, f.Name+".html"))
			if err != nil {
				return nil, err
			}
			f.Path = fp
			destFp := filepath.Join(*dest, fp)

			err = os.MkdirAll(filepath.Dir(destFp), 0755)
			if err != nil {
				log.Fatal(err)
			}
			file, err := os.Create(destFp)
			if err != nil {
				return nil, err
			}

			err = tmpl.Execute(file, &CoveragePair{c, f})
			if err != nil {
				return nil, err
			}
		}

		accCov.Files = append(accCov.Files, c.Files...)
	}

	return &accCov, nil
}

func buildIndex(accCov *Gcov) error {
	tmpl := template.New("index.tmpl")
	funcMap := template.FuncMap{
		"toSlash":      filepath.ToSlash,
	}
	tmpl = tmpl.Funcs(funcMap)

	tmpl, err := tmpl.ParseFS(templates, "tmpl/index.tmpl")
	if err != nil {
		return err
	}

	fp := filepath.Join(*dest, "index.html")
	indexFile, err := os.Create(fp)
	if err != nil {
		return err
	}

	err = tmpl.Execute(indexFile, accCov)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)

	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("Missing GCOV JSON file argument")
	}
	files := flag.Args()

	err := os.MkdirAll(*dest, 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = createCSS(filepath.Join(*dest, "style.css"))
	if err != nil {
		log.Fatal(err)
	}

	accCov, err := buildFiles(files)
	if err != nil {
		log.Fatal(err)
	}
	err = buildIndex(accCov)
	if err != nil {
		log.Fatal(err)
	}
}
