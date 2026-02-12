package composite

type ContextCell struct {
	Type    CellIdentity
	Marks   bool
	Visited bool
}

func NewContextCell(identity CellIdentity, marked bool, visited bool) ContextCell {
	return ContextCell{
		Type:    identity,
		Marks:   marked,
		Visited: visited,
	}
}

func (ctx *ContextCell) ToCell() Cell {
	return NewCell(ctx.Type, ctx.Marks, ctx.Visited)
}
