package enber

import (
	"embed"
	"encoding/json"
	"path"

	"entgo.io/ent/entc/gen"
)

//go:embed templates
var assets embed.FS

func (e *extension) Hooks() []gen.Hook {
	return e.hooks
}

func (e *extension) generate(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		e.TemplateData = &templateData{
			Config:       e.Config,
			TypesImports: []string{},
			InputImports: []string{},
		}
		e.parseInputNode(g)
		e.parseQuery(g)

		files := []file{
			{
				Path:   "ent/enber_input.go",
				Buffer: parseTemplate("enber/enber_input.go.tmpl", e.TemplateData),
			},
			{
				Path:   "ent/enber_query.go",
				Buffer: parseTemplate("enber/enber_query.go.tmpl", e.TemplateData),
			},
			{
				Path:   "ent/enber_typescript.ts",
				Buffer: parseTemplate("enber/enber_typescript.go.tmpl", e.TemplateData),
			},
		}

		b, _ := json.Marshal(e.TemplateData)
		writeFile(file{
			Path:   path.Join(e.Config.App.RootPath, "debug.json"),
			Buffer: fixString(string(b)),
		})

		e.writeFiles(files)
		return next.Generate(g)
	})
}

func (e *extension) writeFiles(files []file) {
	for _, f := range files {
		f.Path = path.Join(e.Config.App.RootPath, f.Path)
		writeFile(f)
	}
}

func NewExtension(options ...extensionOption) *extension {
	ex := &extension{
		Config: config{
			App: AppConfig{
				Pkg:      "app",
				RootPath: "../../",
			},
		},
	}
	for i := range options {
		options[i](ex)
	}
	gen.Funcs["snake"] = snake
	gen.Funcs["camel"] = camel
	ex.hooks = append(ex.hooks, ex.generate)
	return ex
}
