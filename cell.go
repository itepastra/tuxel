package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/PaesslerAG/gval"
)

var language = gval.Full()

var cellmatch = regexp.MustCompile(`[A-Z]{1-3}[\-0-9]{1-6}`)

type Dependencies map[*Cell]string

type Cell struct {
	ftext         string
	formula       gval.Evaluable
	deps          Dependencies
	value         interface{}
	displayString string
	coords        Coords
	broker        *Broker
}

func newCell(x int, y int, broker *Broker) Cell {
	c := Cell{
		formula:       nil,
		deps:          make(Dependencies, 0),
		value:         nil,
		displayString: "",
		coords:        Coords{x: x, y: y},
		broker:        broker,
	}
	return c
}

func (c *Cell) changeFormula(new string) error {
	for o := range c.deps {
		c.broker.unsubscribe <- CellSub{
			subcell:  o,
			thiscell: c,
		}
		delete(c.deps, o)
	}
	if new == "" {
		c = nil
	}
	f, err := c.parseFomula()
	if err != nil {
		return err
	}
	c.formula = f
	for o := range c.deps {
		c.broker.subscribe <- CellSub{
			subcell:  o,
			thiscell: c,
		}
	}
	return nil
}

func (c *Cell) Exec() {
	depvals := make(map[string]interface{}, len(c.deps))
	for dep, name := range c.deps {
		depvals[name] = dep.value
	}
	res, err := c.formula(context.Background(), depvals)
	if err != nil {
		c.value = fmt.Sprint(err)
	} else if res == c.value {
		return
	} else {
		c.value = res
	}
	c.displayString = fmt.Sprint(c.value)
	c.broker.runevent <- CellMsg{from: c}
}

func (c *Cell) parseFomula() (gval.Evaluable, error) {
	args := removeDuplicate(cellmatch.FindAllString(c.ftext, -1))
	deps := make(map[*Cell]string, len(args))
	grid := c.broker.grid
	for _, arg := range args {
		dep, err := grid.getcellstr(arg)
		if err != nil {
			return nil, err
		}
		deps[dep] = arg
	}

	return language.NewEvaluable(c.ftext)
}
