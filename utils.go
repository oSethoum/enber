package enber

import (
	"bytes"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"entgo.io/ent/entc/gen"
)

var (
	_camel   = gen.Funcs["camel"].(func(string) string)
	snake    = gen.Funcs["snake"].(func(string) string)
	singular = gen.Funcs["singular"].(func(string) string)
	camel    = func(s string) string { return _camel(snake(s)) }
)

func vin[T comparable](v T, a []T) bool {
	for _, va := range a {
		if v == va {
			return true
		}
	}
	return false
}

func catch(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// func printStruct(v interface{}) {
// 	b, _ := json.Marshal(v)
// 	println(string(b))
// }

func parseTemplate(filename string, data *templateData) string {
	buffer, err := assets.ReadFile("templates/" + filename)
	catch(err)
	t, err := template.New(filename).Funcs(gen.Funcs).Parse(string(buffer))
	catch(err)
	out := bytes.Buffer{}
	err = t.Execute(&out, data)
	catch(err)
	return out.String()
}

func (e *extension) parseInputNode(nodes []*gen.Type) []*inputNode {
	inputNodes := []*inputNode{}
	for _, n := range nodes {
		if len(n.Edges)+len(n.Fields) == 0 {
			continue
		}
		in := &inputNode{
			Name:        n.Name,
			ShouldType:  true,
			ShouldInput: true,
		}
		in.CreateFields = []*inputField{}
		in.UpdateFields = []*inputField{}
		in.CreateEdges = []*inputEdge{}
		in.UpdateEdges = []*inputEdge{}

		for _, f := range n.Fields {
			if f.IsEdgeField() {
				continue
			}
			var clearField *inputField
			isSlice := strings.HasPrefix(f.Type.String(), "[]")
			if f.Type.PkgPath != "" && !vin(f.Type.PkgPath, e.TemplateData.TypesImports) {
				e.TemplateData.TypesImports = append(e.TemplateData.TypesImports, f.Type.PkgPath)
			}
			if f.Enums != nil {
				pkgName := path.Join(e.Config.App.Pkg, "ent", strings.ToLower(n.Name))
				if !vin(pkgName, e.TemplateData.TypesImports) {
					e.TemplateData.TypesImports = append(e.TemplateData.TypesImports, pkgName)
				}
			}
			createField := &inputField{
				Name:     f.StructField(),
				Type:     f.Type.String(),
				Set:      "Set" + f.StructField(),
				SetParam: "i." + f.StructField(),
			}
			updateField := &inputField{
				Name:     f.StructField(),
				Type:     "*" + f.Type.String(),
				Set:      "Set" + f.StructField(),
				SetParam: "*i." + f.StructField(),
			}

			if isSlice {
				updateField.Check = true
				updateField.SetParam = "i." + f.StructField()
				updateField.Type = f.Type.String()
				if f.Optional {
					clearField = &inputField{
						Name:  "Clear" + f.StructField(),
						Set:   "Clear" + f.StructField(),
						Type:  "*bool",
						Clear: true,
						Check: true,
					}
				}
			} else {
				if f.Optional {
					createField.Type = "*" + f.Type.String()
					createField.Check = true
					createField.SetParam = "*i." + f.StructField()

					updateField.Set = "Set" + f.StructField()
					updateField.Type = "*" + f.Type.String()
					updateField.Check = true

					clearField = &inputField{
						Name:  "Clear" + f.StructField(),
						Set:   "Clear" + f.StructField(),
						Type:  "*bool",
						Check: true,
						Clear: true,
					}
				}

				if f.Default {
					createField.SetParam = "*i." + f.StructField()
					createField.Type = "*" + f.Type.String()
					createField.Check = true

					updateField.Set = "Set" + f.StructField()
					updateField.Check = true
				}
			}

			in.CreateFields = append(in.CreateFields, createField)
			in.UpdateFields = append(in.UpdateFields, updateField)
			if clearField != nil {
				in.UpdateFields = append(in.UpdateFields, clearField)
			}
		}

		for _, e := range n.EdgesWithID() {
			if e.Unique {
				createEdge := &inputEdge{
					Name:     e.StructField() + "ID",
					Type:     e.Type.IDType.String(),
					JTag:     camel(e.StructField()) + "ID" + ",omitempty",
					Set:      "Set" + e.StructField() + "ID",
					SetParam: "i." + e.StructField() + "ID",
				}
				updateEdge := &inputEdge{
					Name:     e.StructField() + "ID",
					Type:     "*" + e.Type.IDType.String(),
					Set:      "Set" + e.StructField() + "ID",
					JTag:     camel(e.StructField()) + "ID" + ",omitempty",
					SetParam: "*i." + e.StructField() + "ID",
					Check:    true,
				}
				if e.Optional {
					// ex: SetNillable
					createEdge.Set = "Set" + e.StructField() + "ID"
					updateEdge.Set = "Set" + e.StructField() + "ID"
					createEdge.Type = "*" + createEdge.Type
					// ex: delete two lines
					createEdge.SetParam = "*i." + e.StructField() + "ID"
					createEdge.Check = true
				}
				in.CreateEdges = append(in.CreateEdges, createEdge)
				in.UpdateEdges = append(in.UpdateEdges, updateEdge)
			} else {
				ea := &inputEdge{
					Name:  "Add" + singular(e.StructField()) + "IDs",
					Set:   "Add" + singular(e.StructField()) + "IDs",
					Type:  "[]" + e.Type.IDType.String(),
					JTag:  camel("Add"+singular(e.StructField())) + "IDs" + ",omitempty",
					Slice: true,
				}
				ea.SetParam = "i." + ea.Name + "..."
				er := &inputEdge{
					Name:  "Remove" + singular(e.StructField()) + "IDs",
					Set:   "Remove" + singular(e.StructField()) + "IDs",
					JTag:  camel("Remove"+singular(e.StructField())) + "IDs" + ",omitempty",
					Type:  "[]" + e.Type.IDType.String(),
					Slice: true,
				}
				er.SetParam = "i." + er.Name + "..."
				ec := &inputEdge{
					Name:  "Clear" + e.StructField(),
					Set:   "Clear" + e.StructField(),
					JTag:  camel("Clear"+e.StructField()) + ",omitempty",
					Type:  "*bool",
					Check: true,
					Clear: true,
				}
				in.CreateEdges = append(in.CreateEdges, ea)
				in.UpdateEdges = append(in.UpdateEdges, ea, er, ec)
			}
		}
		inputNodes = append(inputNodes, in)
	}
	return inputNodes
}

func writeFile(f file) {
	err := os.MkdirAll(path.Dir(f.Path), 0666)
	catch(err)
	os.WriteFile(f.Path, []byte(f.Buffer), 0666)
}

func (e *extension) jgraphy(g *gen.Graph) *jgraph {
	jg := &jgraph{}
	jns := []*jnode{}
	for _, n := range g.Nodes {
		jn := &jnode{
			Name: n.Name,
		}
		jfs := []*jfield{}
		for _, f := range n.Fields {
			jf := &jfield{
				Name:          f.Name,
				StructField:   f.StructField(),
				Optional:      f.Optional,
				Nillable:      f.Nillable,
				Unique:        f.Unique,
				Default:       f.Default,
				Enums:         f.Enums,
				TypePkgName:   f.Type.PkgName,
				TypePkgPath:   f.Type.PkgPath,
				Type:          f.Type.String(),
				TypeConstName: f.Type.ConstName(),
				TypeIdent:     f.Type.Ident,
			}
			jfs = append(jfs, jf)
		}

		jes := []*jedge{}
		for _, e := range n.Edges {
			je := &jedge{
				Name:     e.Name,
				Optional: e.Optional,
				Inverse:  e.Inverse,
				Unique:   e.Unique,
			}

			jes = append(jes, je)
		}
		jn.Edges = jes
		jn.Fields = jfs
		jns = append(jns, jn)
	}
	jg.Nodes = jns
	return jg
}
