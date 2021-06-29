package main

import (
	"embed"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

var (
	dest = "./www"
)

//go:embed tmpl
var templates embed.FS

type SourcePair struct {
	Coverage   *Gcov
	SourceCode string
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

	funcMap := template.FuncMap{
		"increment": increment,
		"getLines":  getLines,
	}
	tmpl = tmpl.Funcs(funcMap)

	tmpl, err = tmpl.ParseFS(templates, "tmpl/*.tmpl")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func main() {
	files := os.Args[1:]

	err := os.MkdirAll(dest, 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = createCSS(filepath.Join(dest, "style.css"))
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := buildHTML()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		c, err := OpenGcov(file)
		if err != nil {
			log.Fatal(err)
		}

		if len(c.Files) != 1 {
			panic("support for mutiple files not implemented yet")
		}
		f := c.Files[0]

		fn := filepath.Base(f.Name) + ".html"
		fp := filepath.Join(dest, fn)
		file, err := os.Create(fp)
		if err != nil {
			log.Fatal(err)
		}

		cfp := filepath.Join(c.CWD, f.Name)
		data, err := os.ReadFile(cfp)
		if err != nil {
			log.Fatal(err)
		}

		p := &SourcePair{Coverage: c, SourceCode: string(data)}
		err = tmpl.Execute(file, p)
		if err != nil {
			log.Fatal(err)
		}
	}
}
