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
	snake    = gen.Funcs["snake"].(func(string) string)
	singular = gen.Funcs["singular"].(func(string) string)
	camel    = func(s string) string { return gen.Funcs["camel"].(func(string) string)(snake(s)) }
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

func (e *extension) parseInputNode(g *gen.Graph) {
	inputNodes := []*inputNode{}
	for _, n := range g.Nodes {
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
			fieldName := f.StructField()
			fieldType := f.Type.String()
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
				Name:     fieldName,
				Type:     fieldType,
				Set:      "Set" + fieldName,
				SetParam: "i." + fieldName,
			}
			updateField := &inputField{
				Name:     fieldName,
				Type:     "*" + fieldType,
				Set:      "Set" + fieldName,
				SetParam: "*i." + fieldName,
			}

			if isSlice {
				updateField.Check = true
				updateField.SetParam = "i." + fieldName
				updateField.Type = fieldType
				if f.Optional {
					clearField = &inputField{
						Name:  "Clear" + fieldName,
						Set:   "Clear" + fieldName,
						Type:  "*bool",
						Clear: true,
						Check: true,
					}
				}
			} else {
				if f.Optional {
					createField.Type = "*" + fieldType
					createField.Check = true
					createField.SetParam = "*i." + fieldName

					updateField.Set = "Set" + fieldName
					updateField.Type = "*" + fieldType
					updateField.Check = true

					clearField = &inputField{
						Name:  "Clear" + fieldName,
						Set:   "Clear" + fieldName,
						Type:  "*bool",
						Check: true,
						Clear: true,
					}
				}

				if f.Default {
					createField.SetParam = "*i." + fieldName
					createField.Type = "*" + fieldType
					createField.Check = true

					updateField.Set = "Set" + fieldName
					updateField.Check = true
				}

				if f.Type.ConstName() != "TypeJSON" {
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
			edgeName := e.StructField()
			if e.Unique {
				createEdge := &inputEdge{
					Name:     edgeName + "ID",
					Type:     e.Type.IDType.String(),
					JTag:     camel(edgeName) + "ID" + ",omitempty",
					Set:      "Set" + edgeName + "ID",
					SetParam: "i." + edgeName + "ID",
				}
				updateEdge := &inputEdge{
					Name:     edgeName + "ID",
					Type:     "*" + e.Type.IDType.String(),
					Set:      "Set" + edgeName + "ID",
					JTag:     camel(edgeName) + "ID" + ",omitempty",
					SetParam: "*i." + edgeName + "ID",
					Check:    true,
				}
				if e.Optional {
					createEdge.Set = "Set" + edgeName + "ID"
					updateEdge.Set = "Set" + edgeName + "ID"
					createEdge.Type = "*" + createEdge.Type
					createEdge.SetParam = "*i." + edgeName + "ID"
					createEdge.Check = true
				}
				in.CreateEdges = append(in.CreateEdges, createEdge)
				in.UpdateEdges = append(in.UpdateEdges, updateEdge)
			} else {
				ea := &inputEdge{
					Name:  "Add" + singular(edgeName) + "IDs",
					Set:   "Add" + singular(edgeName) + "IDs",
					Type:  "[]" + e.Type.IDType.String(),
					JTag:  camel("Add"+singular(edgeName)) + "IDs" + ",omitempty",
					Slice: true,
				}
				ea.SetParam = "i." + ea.Name + "..."
				er := &inputEdge{
					Name:  "Remove" + singular(edgeName) + "IDs",
					Set:   "Remove" + singular(edgeName) + "IDs",
					JTag:  camel("Remove"+singular(edgeName)) + "IDs" + ",omitempty",
					Type:  "[]" + e.Type.IDType.String(),
					Slice: true,
				}
				er.SetParam = "i." + er.Name + "..."
				ec := &inputEdge{
					Name:  "Clear" + edgeName,
					Set:   "Clear" + edgeName,
					JTag:  camel("Clear"+edgeName) + ",omitempty",
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
	e.TemplateData.InputNodes = inputNodes
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

func (e *extension) parseQuery(g *gen.Graph) {
	qns := []*queryNode{}

	for _, n := range g.Nodes {
		qn := &queryNode{
			Name: n.Name,
		}

		if !vin("errors", e.TemplateData.QueryImports) {
			e.TemplateData.QueryImports = append(e.TemplateData.QueryImports, "errors")
		}

		e.TemplateData.QueryImports = append(e.TemplateData.QueryImports,
			path.Join(e.Config.App.Pkg, "ent", strings.ToLower(n.Name)),
		)

		for _, f := range n.Fields {
			qf := &queryField{
				Name:            f.StructField(),
				Enum:            f.Type.ConstName() == "TypeEnum",
				Boolean:         f.Type.ConstName() == "TypeBool",
				Optional:        f.Optional,
				TypeString:      f.Type.String(),
				Comparable:      f.Type.Comparable() && f.Type.ConstName() != "TypeBool",
				String:          f.Type.ConstName() == "TypeString",
				EdgeFieldOrEnum: f.IsEdgeField() || f.Type.ConstName() == "TypeEnum",
			}
			qf.WithComment = qf.Boolean || qf.Comparable || qf.Optional
			qn.Fields = append(qn.Fields, qf)
			if f.Type.ConstName() == "TypeTime" && !vin("time", e.TemplateData.QueryImports) {
				e.TemplateData.QueryImports = append(e.TemplateData.QueryImports, "time")
			}
		}

		for _, e := range n.Edges {
			qe := &queryEdge{
				Name: e.Name,
				Node: e.Type.Name,
			}
			qn.Edges = append(qn.Edges, qe)
		}

		qns = append(qns, qn)
	}
	e.TemplateData.QueryNodes = qns
}
