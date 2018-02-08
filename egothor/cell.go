package egothor

// Cell in a Row of a Trie
type Cell struct {
	ref  int32
	cmd  int32
	cnt  int32
	skip int32
}

//NewCell returns new cell
func NewCell(ref, cmd, cnt, skip int32) *Cell {
	return &Cell{ref: ref, cmd: cmd, cnt: cnt, skip: skip}
}

// Cmd returns cmd from cell
func (c Cell) Cmd() int32 {
	return c.cmd
}

// Ref returns ref from cell
func (c Cell) Ref() int32 {
	return c.ref
}

// Skip returns skip from cell
func (c Cell) Skip() int32 {
	return c.skip
}

// Cnt returns cnt from cell
func (c Cell) Cnt() int32 {
	return c.cnt
}
