package enber

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"entgo.io/ent/entc/gen"
)

var (
	plural   = gen.Funcs["plural"].(func(string) string)
	snake    = gen.Funcs["snake"].(func(string) string)
	singular = gen.Funcs["singular"].(func(string) string)
	oldCamel = gen.Funcs["camel"].(func(string) string)
	camel    = func(s string) string { return oldCamel(snake(s)) }
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
			Name: n.Name,
		}

		in.CreateFields = []*inputField{}
		in.UpdateFields = []*inputField{}
		in.CreateEdges = []*inputEdge{}
		in.UpdateEdges = []*inputEdge{}

		fields := n.Fields
		in.IDType = n.IDType.String()

		if n.ID.UserDefined {
			fields = append([]*gen.Field{n.ID}, fields...)
		}

		for _, f := range fields {
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
				TsType:   gots[fieldType],
				Set:      "Set" + fieldName,
				SetParam: "i." + fieldName,
			}
			updateField := &inputField{
				Name:     fieldName,
				Type:     "*" + fieldType,
				TsType:   gots[fieldType],
				Set:      "Set" + fieldName,
				SetParam: "*i." + fieldName,
			}

			if isSlice {
				createField.Type = fieldType
				createField.TsType = gots[fieldType]

				updateField.Check = true
				updateField.TsCheck = true
				updateField.SetParam = "i." + fieldName
				updateField.Type = fieldType
				updateField.TsType = gots[fieldType]
				if f.Optional {
					clearField = &inputField{
						Name:    "Clear" + fieldName,
						Set:     "Clear" + fieldName,
						Type:    "*bool",
						TsType:  gots[fieldType],
						Clear:   true,
						Check:   true,
						TsCheck: true,
					}
				}
			} else {
				if f.Optional {
					createField.Type = "*" + fieldType
					createField.Check = true
					createField.TsCheck = true
					createField.SetParam = "*i." + fieldName
					createField.TsType = gots[fieldType]

					updateField.Set = "Set" + fieldName
					updateField.Type = "*" + fieldType
					updateField.Check = true
					updateField.TsCheck = true
					updateField.TsType = gots[fieldType]

					clearField = &inputField{
						Name:    "Clear" + fieldName,
						Set:     "Clear" + fieldName,
						Type:    "*bool",
						TsType:  gots["bool"],
						TsCheck: true,
						Check:   true,
						Clear:   true,
					}
				}

				if f.Default {
					createField.SetParam = "*i." + fieldName
					createField.Type = "*" + fieldType
					createField.Check = true
					createField.TsCheck = true

					updateField.Set = "Set" + fieldName
					updateField.Check = true
					updateField.TsCheck = true
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
					TsType:   gots[e.Type.IDType.String()],
					SetParam: "i." + edgeName + "ID",
				}
				updateEdge := &inputEdge{
					Name:     edgeName + "ID",
					Type:     "*" + e.Type.IDType.String(),
					Set:      "Set" + edgeName + "ID",
					JTag:     camel(edgeName) + "ID" + ",omitempty",
					SetParam: "*i." + edgeName + "ID",
					TsType:   gots[e.Type.IDType.String()],
					Check:    true,
					TsCheck:  true,
				}
				if e.Optional {
					createEdge.Set = "Set" + edgeName + "ID"
					updateEdge.Set = "Set" + edgeName + "ID"
					createEdge.Type = "*" + createEdge.Type
					createEdge.SetParam = "*i." + edgeName + "ID"
					createEdge.Check = true
					createEdge.TsCheck = true
				}

				in.CreateEdges = append(in.CreateEdges, createEdge)
				in.UpdateEdges = append(in.UpdateEdges, updateEdge)
			} else {
				ea := &inputEdge{
					Name:    "Add" + singular(edgeName) + "IDs",
					Set:     "Add" + singular(edgeName) + "IDs",
					Type:    "[]" + e.Type.IDType.String(),
					JTag:    camel("Add"+singular(edgeName)) + "IDs" + ",omitempty",
					TsType:  gots[e.Type.IDType.String()] + "[]",
					TsCheck: true,
					Slice:   true,
				}
				ea.SetParam = "i." + ea.Name + "..."
				er := &inputEdge{
					Name:    "Remove" + singular(edgeName) + "IDs",
					Set:     "Remove" + singular(edgeName) + "IDs",
					JTag:    camel("Remove"+singular(edgeName)) + "IDs" + ",omitempty",
					Type:    "[]" + e.Type.IDType.String(),
					TsType:  gots[e.Type.IDType.String()] + "[]",
					TsCheck: true,
					Slice:   true,
				}
				er.SetParam = "i." + er.Name + "..."
				ec := &inputEdge{
					Name:    "Clear" + edgeName,
					Set:     "Clear" + edgeName,
					JTag:    camel("Clear"+edgeName) + ",omitempty",
					Type:    "*bool",
					TsType:  gots[e.Type.IDType.String()],
					Check:   true,
					TsCheck: true,
					Clear:   true,
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

func (e *extension) parseQuery(g *gen.Graph) {
	qns := []*queryNode{}

	for _, n := range g.Nodes {
		qn := &queryNode{
			Name: n.Name,
		}
		qn.IDType = n.IDType.String()
		fields := append([]*gen.Field{n.ID}, n.Fields...)

		if !vin("errors", e.TemplateData.QueryImports) {
			e.TemplateData.QueryImports = append(e.TemplateData.QueryImports, "errors")
		}

		e.TemplateData.QueryImports = append(e.TemplateData.QueryImports,
			path.Join(e.Config.App.Pkg, "ent", strings.ToLower(n.Name)),
		)

		for _, f := range fields {
			if f.Sensitive() {
				continue
			}
			qf := &queryField{
				Name:            f.StructField(),
				Enum:            f.Type.ConstName() == "TypeEnum",
				Boolean:         f.Type.ConstName() == "TypeBool",
				Optional:        f.Optional,
				TsType:          gots[f.Type.String()],
				TypeString:      f.Type.String(),
				Comparable:      f.Type.Comparable() && f.Type.ConstName() != "TypeBool",
				String:          f.Type.ConstName() == "TypeString",
				EdgeFieldOrEnum: f.IsEdgeField() || f.Type.ConstName() == "TypeEnum",
			}
			qf.WithComment = qf.Boolean || qf.Comparable || qf.Optional
			if f.Type.ConstName() == "TypeTime" && !vin("time", e.TemplateData.QueryImports) {
				e.TemplateData.QueryImports = append(e.TemplateData.QueryImports, "time")
			}

			skips := []fieldSkip{}
			if err := decodeAnnotation(f.Annotations[enberFieldSkip], &skips); err != nil {
				log.Fatalln(err.Error())
			}
			qf.Order = f.Type.Comparable() && !vin(FieldSkipOrder, skips)
			qn.Fields = append(qn.Fields, qf)
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

func fixString(s string) string {
	s = strings.ReplaceAll(s, "\\u0026", "?")
	return s
}

func enberorder(n queryNode) string {
	fields := []string{}
	for _, f := range n.Fields {
		if f.Order {
			fields = append(fields, "\""+snake(f.Name)+"\"")
		}
	}
	return strings.Join(fields, " | ")
}
func enberselect(n queryNode) string {
	fields := []string{}
	for _, f := range n.Fields {
		fields = append(fields, "\""+snake(f.Name)+"\"")
	}
	return strings.Join(fields, " | ")
}

func decodeAnnotation(v, out any) error {
	a := &annotation{}
	buffer, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buffer, a)
	if err != nil {
		return err
	}
	buffer, err = json.Marshal(a.Data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buffer, out)
	return err
}
