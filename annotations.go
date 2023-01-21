package enber

import (
	"encoding/json"

	"entgo.io/ent/schema"
)

type enberSkip uint

const (
	enberFieldSkip            = "enber.field.skip"
	SkipCreateInput enberSkip = iota
	SkipUpdateInput
	SkipWhereInput
)

type annotation struct {
	schema.Annotation
	name string
	Data any
}

func (a *annotation) Name() string {
	return a.name
}

func Skip(skips ...enberSkip) *annotation {
	return &annotation{
		name: enberFieldSkip,
		Data: true,
	}
}

func decode(v, out any) error {
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
