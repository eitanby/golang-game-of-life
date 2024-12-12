package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type Universe [][]bool

type Screen struct {
	lifeIcon    string
	deathIcon   string
	refreshRate time.Duration
}

func (screen Screen) Show(universe Universe) {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	text := universe.GetUniverseAsText(screen.lifeIcon, screen.deathIcon)
	fmt.Print(text)
	time.Sleep(screen.refreshRate)
}

func (universe Universe) clone() Universe {
	clone := make([][]bool, len(universe))
	for i := range universe {
		clone[i] = make([]bool, len(universe[i]))
		copy(clone[i], universe[i])
	}
	return clone
}

func NewUniverse(width, height int) Universe {
	universe := make([][]bool, height)
	for i := range universe {
		universe[i] = make([]bool, width)
	}
	return universe
}

func (universe Universe) GetUniverseAsText(lifeIcon, deathIcon string) string {

	getRowAsText := func(row []bool) string {
		text := ""
		for _, cell := range row {
			if cell {
				text += lifeIcon
			} else {
				text += deathIcon
			}
		}
		return text + "\n"
	}

	text := ""
	for _, row := range universe {
		text += getRowAsText(row)
	}
	return text
}

func (universe Universe) Seed(seed int64, lifeProbability int) Universe {
	universeSeeded := universe.clone()
	randomGen := rand.New(rand.NewSource(seed))
	for row := range universeSeeded {
		for cell := range universeSeeded[row] {
			random := randomGen.Intn(100)
			if random <= lifeProbability {
				universeSeeded[row][cell] = true
			} else {
				universeSeeded[row][cell] = false
			}
		}
	}

	return universeSeeded

}

func (universe Universe) Next() Universe {

	isCellAlive := func(universe Universe, x, y int) bool {
		adjust := func(index, limit int) int {
			if index < 0 {
				index = index + limit
			} else if index >= limit {
				index = index % limit
			}
			return index

		}

		x = adjust(x, len(universe[0]))
		y = adjust(y, len(universe))
		return universe[y][x] == true
	}

	countNeighbors := func(universe Universe, x, y int) int {
		neighbors := 0
		for xOffset := -1; xOffset <= 1; xOffset++ {
			for yOffset := -1; yOffset <= 1; yOffset++ {
				if xOffset == 0 && yOffset == 0 {
					continue
				}
				if isCellAlive(universe, x+xOffset, y+yOffset) {
					neighbors++
				}
			}
		}
		return neighbors

	}

	nextCell := func(universe Universe, x, y int) bool {
		neighbors := countNeighbors(universe, x, y)
		isAlive := isCellAlive(universe, x, y)
		if isAlive && (neighbors == 2 || neighbors == 3) {
			return true
		}
		if isAlive == false && neighbors == 3 {
			return true
		}
		return false
	}

	nextUniverse := universe.clone()
	for x, row := range universe {
		for y := range row {
			nextUniverse[x][y] = nextCell(universe, x, y)
		}
	}
	return nextUniverse
}

var (
	width, height, lifeProbability, refresh, totalSteps int
	lifeIcon, deathIcon                                 string
	seed                                                int64
	screen                                              Screen
	universes                                           []Universe
)

func main() {
	flag.IntVar(&width, "width", 10, "universe width")
	flag.IntVar(&height, "height", 10, "universe height")
	flag.IntVar(&refresh, "refresh", 1, "game refresh rate (in seconds)")
	flag.IntVar(&totalSteps, "steps", 1000, "how many steps the game run")
	flag.IntVar(&lifeProbability, "lifeProbability", 25, "initial probability for life")
	flag.StringVar(&lifeIcon, "lifeIcon", "💜", "life icon")
	flag.StringVar(&deathIcon, "deathIcon", "💀", "death icon")
	flag.Parse()
	seed = time.Now().UnixNano()
	screen = Screen{lifeIcon, deathIcon, time.Duration(refresh) * time.Second}

	universe := NewUniverse(width, height)
	universe = universe.Seed(seed, lifeProbability)

	for step := 0; step < totalSteps; step++ {
		universes = append(universes, universe)
		universe = universe.Next()
	}

	for _, universe := range universes {
		screen.Show(universe)
	}
}
