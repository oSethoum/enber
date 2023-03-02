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

		if e.Config.DBConfig != nil {
			files = append(files, file{
				Path:   "db/db.go",
				Buffer: parseTemplate("db/db.go.tmpl", e.TemplateData),
			})
		}

		if e.Config.Server != nil {
			files = append(files, file{
				Path:   "routes/enber.go",
				Buffer: parseTemplate("fiber/routes.go.tmpl", e.TemplateData),
			}, file{
				Path:   snake(e.TemplateData.Config.Server.FileName + ".go"),
				Buffer: parseTemplate("fiber/server.go.tmpl", e.TemplateData),
			})

			for i, n := range g.Nodes {
				e.TemplateData.InputNode = e.TemplateData.InputNodes[i]
				e.TemplateData.QueryNode = e.TemplateData.QueryNodes[i]

				files = append(files, file{
					Path:   "handlers/" + plural(snake(n.Name)) + ".go",
					Buffer: parseTemplate("fiber/handlers.go.tmpl", e.TemplateData),
				})
			}
		}

		if e.Config.TsConfig != nil {
			if !strings.HasSuffix(e.Config.TsConfig.Path, ".ts") {
				e.Config.TsConfig.Path += ".ts"
			}
			files = append(files,
				file{
					Path:   e.Config.TsConfig.Path,
					Buffer: parseTemplate("enber/enber_typescript.go.tmpl", e.TemplateData),
				},
			)
		}
		if e.Config.Debug {
			e.debug()
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
	gen.Funcs["camel"] = camel
	gen.Funcs["enberorder"] = enberorder
	gen.Funcs["enberselect"] = enberselect
	ex.hooks = append(ex.hooks, ex.generate)
	return ex
}
