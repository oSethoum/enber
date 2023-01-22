package enber

type (
	fieldSkip  uint
	edgeSkip   float32
	edgeNested int8
	nodeSkip   float64
)

const (
	enberFieldSkip = "enber.field.skip"

	FieldSkipCreateInput fieldSkip = iota
	FieldSkipUpdateInput
	FieldSkipWhereInput
	FieldSkipOrder
)

const (
	enberEdgeSkip = "enber.edge.skip"

	EdgeSkipCreateinput edgeSkip = iota * 3.66
	EdgeSkipUpdateInput
	EdgeSkipQueryInput
)

const (
	enberEdgeNested = "enber.edge.nested"

	EdgeNestedCreate edgeNested = iota * 2
	EdgeNestedUpdate
)

const (
	enberNodeSkip = "enber.node.skip"

	NodeSkipCreate nodeSkip = iota * 2.36
	NodeSkipUpdate
	NodeSkipQuery
	NodeSkipPrivacy
)

type annotation struct {
	name string
	Data any
}

func (a *annotation) Name() string {
	return a.name
}

func FieldSkip(options ...fieldSkip) *annotation {
	return &annotation{
		name: enberFieldSkip,
		Data: options,
	}
}

func EdgeSkip(options ...edgeSkip) *annotation {
	return &annotation{
		name: enberEdgeSkip,
		Data: options,
	}
}

func EdgeNested(options ...edgeNested) *annotation {
	return &annotation{
		name: enberEdgeNested,
		Data: options,
	}
}

func NodeSkip(options ...nodeSkip) *annotation {
	return &annotation{
		name: enberNodeSkip,
		Data: options,
	}
}
