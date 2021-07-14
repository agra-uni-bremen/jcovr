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
	dest = flag.String("d", "./www", "output directory for HTML files")
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
	tmpl, err := tmpl.ParseFS(templates, "tmpl/source.tmpl")
	if err != nil {
		return nil, err
	}

	var accCov Gcov
	for _, file := range files {
		c, err := OpenGcov(file)
		if err != nil {
			return nil, err
		}

		accCov.CWD = c.CWD
		accCov.Files = append(accCov.Files, c.Files...)

		for _, f := range c.Files {
			fn := filepath.Base(f.Name) + ".html"
			fp := filepath.Join(*dest, fn)
			file, err := os.Create(fp)
			if err != nil {
				return nil, err
			}

			err = tmpl.Execute(file, &CoveragePair{c, f})
			if err != nil {
				return nil, err
			}
		}
	}

	return &accCov, nil
}

func buildIndex(accCov *Gcov) error {
	tmpl := template.New("index.tmpl")
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
