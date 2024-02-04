package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jason-meredith/warships/base26"
)

type Grid struct {
	cells    map[int]map[int]*Cell
	viewport Viewport
}

func (g *Grid) getcell(c Coords) *Cell {
	if g.cells[c.x] == nil {
		return nil
	}
	return g.cells[c.x][c.y]
}

func (g *Grid) cellToStr(c Coords) string {
	if g.cells[c.x] == nil {
		minus := ""
		if c.x < 0 {
			minus = "-"
		}
		return fmt.Sprintf("%s%s%d", minus, base26.ConvertToBase26(absi(c.x)), c.y)
	}
	return g.cells[c.x][c.y].displayString
}

func toCoords(s string) (Coords, error) {
	fst := strings.IndexAny(s[1:], "-0123456789") + 1
	var x int
	if s[0] == '-' {
		x = -base26.ConvertToDecimal(s[1:fst])
	} else {
		x = base26.ConvertToDecimal(s[:fst])
	}
	y, err := strconv.Atoi(s[fst:])
	if err != nil {
		return Coords{0, 0}, err
	}
	return Coords{x, y}, nil
}

func (g *Grid) getcellstr(s string) (*Cell, error) {
	coords, err := toCoords(s)
	if err != nil {
		return nil, err
	}
	return g.getcell(coords), nil
}

func (g *Grid) setcell(x int, y int, cell *Cell) {

	if g.cells[x] == nil {
		g.cells[x] = make(map[int]*Cell)
	}
	g.cells[x][y] = cell
}

func newGrid() *Grid {
	return &Grid{
		cells:    make(map[int]map[int]*Cell),
		viewport: Viewport{xmin: 0, xmax: 3, ymin: 0, ymax: 10},
	}
}
func (g *Grid) At(row, cell int) string {
	if cell == 0 {
		return fmt.Sprint(row + g.viewport.ymin)
	}
	return g.cellToStr(Coords{x: cell + g.viewport.xmin - 1, y: row + g.viewport.ymin})
}
func (g *Grid) Rows() int {
	return g.viewport.ymax - g.viewport.ymin
}
func (g *Grid) Columns() int {
	return g.viewport.xmax - g.viewport.xmin
}
