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

type Page struct {
	Name     string
	Coverage *Gcov
	File     *GcovFile // May be nil
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

func buildHTML() (*template.Template, error) {
	var err error

	const name = "base.tmpl"
	tmpl := template.New(name)

	tmpl, err = tmpl.ParseFS(templates, "tmpl/*.tmpl")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func buildIndex(tmpl *template.Template, path string, accCov *Gcov) error {
	indexFile, err := os.Create(path)
	if err != nil {
		return err
	}

	indexPage := &Page{
		Name: "index",
		Coverage: accCov,
		File: nil,
	}

	err = tmpl.Execute(indexFile, indexPage)
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

	tmpl, err := buildHTML()
	if err != nil {
		log.Fatal(err)
	}

	var accCov Gcov
	for _, file := range files {
		c, err := OpenGcov(file)
		if err != nil {
			log.Fatal(err)
		}
		accCov.Files = append(accCov.Files, c.Files...)

		for _, f := range c.Files {
			fn := filepath.Base(f.Name) + ".html"
			fp := filepath.Join(*dest, fn)
			file, err := os.Create(fp)
			if err != nil {
				log.Fatal(err)
			}

			p := &Page{
				Name: "source",
				Coverage: c,
				File: f,
			}

			err = tmpl.Execute(file, p)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fp := filepath.Join(*dest, "index.html")
	err = buildIndex(tmpl, fp, &accCov)
	if err != nil {
		log.Fatal(err)
	}
}
