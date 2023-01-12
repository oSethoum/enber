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

func (e *extension) debug(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		if e.DebugInfo.SchemaJson {
			b, _ := json.Marshal(e.jgraphy(g))
			writeFile(file{
				Path:   path.Join(e.Config.App.RootPath, "graph.json"),
				Buffer: string(b),
			})
		}
		return next.Generate(g)
	})
}

func (e *extension) generate(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		e.TemplateData = &templateData{
			Nodes:        g.Nodes,
			Config:       e.Config,
			TypesImports: []string{},
			InputImports: []string{},
		}

		e.TemplateData.InputNodes = e.parseInputNode(g.Nodes)

		files := []file{
			{
				Path:   "ent/enber_input.go",
				Buffer: parseTemplate("enber/enber_input.go.tmpl", e.TemplateData),
			},
			// {
			// 	Path:   "ent/enber_types.go",
			// 	Buffer: parseTemplate("enber/enber_types.go.tmpl", e.TemplateData),
			// },
			// {
			// 	Path:   "ent/enber_query.go",
			// 	Buffer: parseTemplate("enber/enber_query.go.tmpl", e.TemplateData),
			// },
		}

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
	if ex.DebugInfo != nil && ex.DebugInfo.DebugOnly {
		ex.hooks = append(ex.hooks, ex.debug)
	} else {
		ex.hooks = append(ex.hooks, ex.debug, ex.generate)
	}
	return ex
}
