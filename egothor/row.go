package egothor

// Row struct
type Row struct {
	cells map[rune]*Cell
}

// NewRow returns row
func NewRow() *Row {
	return &Row{cells: make(map[rune]*Cell)}
}

// AddCell appends cell
func (r *Row) AddCell(char rune, cell *Cell) {
	r.cells[char] = cell
}

// Cmd returns cmd or -1
func (r Row) Cmd(c rune) int32 {
	return r.getCellValue(c, "cmd")
}

// Ref returns ref or -1
func (r Row) Ref(c rune) int32 {
	return r.getCellValue(c, "ref")
}

// getCellValue returns value from cell
// or -1 if there is no cell
func (r Row) getCellValue(ch rune, field string) int32 {
	c, ok := r.cells[ch]
	if !ok {
		return -1
	}
	switch field {
	case "cmd":
		return c.Cmd()
	case "cnt":
		return c.Cnt()
	case "ref":
		return c.Ref()
	case "skip":
		return c.Skip()
	default:
		return -1
	}
}
