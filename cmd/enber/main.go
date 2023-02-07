package main

import (
	"bytes"
	"embed"
	"os"
	"path"
	"text/template"
)

type file struct {
	Path   string
	Buffer string
}

//go:embed templates
var assets embed.FS

func main() {

	writeFile(file{
		Path:   "ent/generate/entc.go",
		Buffer: parseTemplate("entc.go.tmpl", nil),
	})
	writeFile(file{
		Path:   "ent/generate/generate.go",
		Buffer: parseTemplate("generate.go.tmpl", nil),
	})
	os.Mkdir("ent/schema", 0666)
}

func parseTemplate(filename string, data interface{}) string {
	buffer, err := assets.ReadFile("templates/" + filename)
	catch(err)
	t, err := template.New(filename).Parse(string(buffer))
	catch(err)
	out := bytes.Buffer{}
	err = t.Execute(&out, data)
	catch(err)
	return out.String()
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func writeFile(f file) {
	err := os.MkdirAll(path.Dir(f.Path), 0666)
	catch(err)
	os.WriteFile(f.Path, []byte(f.Buffer), 0666)
}
