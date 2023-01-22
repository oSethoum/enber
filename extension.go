package enber

import (
	"embed"
	"encoding/json"
	"path"
	"strings"

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
		}

		if e.Config.tsConfig != nil {
			if !strings.HasSuffix(e.Config.tsConfig.Path, ".ts") {
				e.Config.tsConfig.Path += ".ts"
			}
			files = append(files,
				file{
					Path:   e.Config.tsConfig.Path,
					Buffer: parseTemplate("enber/enber_typescript.go.tmpl", e.TemplateData),
				},
			)
		}
		e.writeFiles(files)
		return next.Generate(g)
	})
}

func (e *extension) debug() {
	b, _ := json.Marshal(e.TemplateData)
	writeFile(file{
		Path:   path.Join(e.Config.App.RootPath, "debug.json"),
		Buffer: fixString(string(b)),
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
	ex.debug()
	gen.Funcs["camel"] = camel
	ex.hooks = append(ex.hooks, ex.generate)
	return ex
}
