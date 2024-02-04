package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jason-meredith/warships/base26"
)

const evenBack = lipgloss.Color("#2f2f2f")
const oddBack = lipgloss.Color("#000000")
const headerColor = lipgloss.Color("#222222")
const selectedCell = lipgloss.Color("#af00af")
const selectedRC = lipgloss.Color("#3f003f")

const selectMargin = 3

var headerStyle = lipgloss.NewStyle().Background(headerColor).Width(10).Align(lipgloss.Right).Inline(true).MaxWidth(10)
var evenStyle = lipgloss.NewStyle().Background(evenBack).Width(10).Align(lipgloss.Right).Inline(true).MaxWidth(10)
var oddStyle = lipgloss.NewStyle().Background(oddBack).Width(10).Align(lipgloss.Right).Inline(true).MaxWidth(10)
var selectedRCStyle = lipgloss.NewStyle().Background(selectedRC).Width(10).Align(lipgloss.Right).Inline(true).MaxWidth(10)
var selectedCellStyle = lipgloss.NewStyle().Background(selectedCell).Width(10).Align(lipgloss.Right).Inline(true).MaxWidth(10)

type Coords struct {
	x int
	y int
}

type Viewport struct {
	xmin int
	ymin int
	xmax int
	ymax int
}

type model struct {
	cells    *Grid
	selected Coords
	table    table.Table
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("h"))):
			m.selected.x--
			if m.selected.x < m.cells.viewport.xmin+selectMargin {
				m.cells.viewport.xmin--
				m.cells.viewport.xmax--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("j"))):
			m.selected.y++
			if m.selected.y >= m.cells.viewport.ymax-selectMargin {
				m.cells.viewport.ymin++
				m.cells.viewport.ymax++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("k"))):
			m.selected.y--
			if m.selected.y < m.cells.viewport.ymin+selectMargin {
				m.cells.viewport.ymin--
				m.cells.viewport.ymax--
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("l"))):
			m.selected.x++
			if m.selected.x >= m.cells.viewport.xmax-1-selectMargin {
				m.cells.viewport.xmin++
				m.cells.viewport.xmax++
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("q"))):
			return m, tea.Quit
		}
		// log.Printf("new selected %#v", m.selected)
		m.table.StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row+m.cells.viewport.ymin == m.selected.y+1 && col+m.cells.viewport.xmin == m.selected.x+1:
				return selectedCellStyle
			case row+m.cells.viewport.ymin == m.selected.y+1,
				col+m.cells.viewport.xmin == m.selected.x+1:
				return selectedRCStyle
			case row == 0:
				return headerStyle
			case row%2 == 0:
				return evenStyle
			default:
				return oddStyle
			}
		})
		return m, nil
	case tea.WindowSizeMsg:
		m.cells.viewport.xmax, m.cells.viewport.ymax = msg.Width/10, msg.Height-4
		return m, nil
	case tea.MouseMsg:
		m.selected.x, m.selected.y = msg.X, msg.Y
		return m, nil
	default:
		return m, nil
	}
}

func absi(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func (m model) View() string {
	m.table.Data(m.cells)
	headers := []string{""}
	for x := m.cells.viewport.xmin; x < m.cells.viewport.xmax; x++ {
		sig := ""
		if x < 0 {
			sig = "-"
		}
		headers = append(headers, fmt.Sprintf("%s%s", sig, base26.ConvertToBase26(absi(x))))
	}
	m.table.Headers(headers...)

	return m.table.Render()
}

var width = flag.Int("w", 100, "the number of columns")
var height = flag.Int("h", 1000, "the number of rows")

func main() {
	broker := NewBroker()
	grid := newGrid()
	broker.grid = grid
	m := newModel(grid)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func newModel(grid *Grid) model {
	m := model{
		cells:    grid,
		selected: Coords{x: 0, y: 0},
	}
	m.table = *table.New().
		Border(lipgloss.Border{
			Top:          "",
			Bottom:       "",
			Left:         "",
			Right:        "",
			TopLeft:      "",
			TopRight:     "",
			BottomLeft:   "",
			BottomRight:  "",
			MiddleLeft:   "",
			MiddleRight:  "",
			Middle:       "",
			MiddleTop:    "",
			MiddleBottom: "",
		}).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row+m.cells.viewport.ymin == m.selected.y+1 && col+m.cells.viewport.xmin == m.selected.x+1:
				return selectedCellStyle
			case row+m.cells.viewport.ymin == m.selected.y+1,
				col+m.cells.viewport.xmin == m.selected.x+1:
				return selectedRCStyle
			case row == 0:
				return headerStyle
			case row%2 == 0:
				return evenStyle
			default:
				return oddStyle
			}
		})
	m.table.Data(m.cells)
	return m
}
