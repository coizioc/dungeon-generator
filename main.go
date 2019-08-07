package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	d := GenerateDungeon(80, 20)
	d.Print()
}

type Dungeon struct {
	Tiles  [][]byte
	Height int
	Width  int
}

type Point struct {
	x int
	y int
}

// GenerateDungeon generates the dungeon.
func GenerateDungeon(w, h int) *Dungeon {
	d := &Dungeon{make([][]byte, h), h, w}
	sizeCutoff := 3 * w * h / 5
	for i := 0; i < h; i++ {
		d.Tiles[i] = make([]byte, w)
	}

	for {
		d.InitializeDungeon()

		// Create caves
		for n := 0; n < 2; n++ {
			d.SimulateStep(3, 5)
		}

		// Create corridors
		for n := 0; n < 3; n++ {
			d.SimulateStep(6, 2)
		}

		d.FillBoarder()
		size := d.FloodFill()
		if size == -1 {
			continue
		} else if size > sizeCutoff {
			break
		}
	}

	return d
}

// InitializeDungeon randomly fills Tiles with walls.
func (d *Dungeon) InitializeDungeon() {
	aliveProb := 0.45
	for y := 0; y < d.Height; y++ {
		for x := 0; x < d.Width; x++ {
			if rand.Float64() <= aliveProb {
				d.Tiles[y][x] = '#'
			} else {
				d.Tiles[y][x] = ' '
			}
		}
	}
}

// SimulateStep simulates a step in the cellular automata process.
func (d *Dungeon) SimulateStep(birthLim, deathLim int) {
	// Create new tilemap.
	newTiles := make([][]byte, d.Height)
	for i := 0; i < d.Height; i++ {
		newTiles[i] = make([]byte, d.Width)
	}

	// For each tile, count its neighbors.
	for y := 0; y < d.Height; y++ {
		for x := 0; x < d.Width; x++ {
			nbs := d.CountAliveNeighbors(x, y)
			// If wall and number of neighbors is less than deathLim, turn to floor.
			if d.Tiles[y][x] == '#' {
				if nbs < deathLim {
					newTiles[y][x] = ' '
				} else {
					newTiles[y][x] = '#'
				}
				// If floor and number of neighbors is greater than birthLim, turn to wall.
			} else {
				if nbs > birthLim {
					newTiles[y][x] = '#'
				} else {
					newTiles[y][x] = ' '
				}
			}
		}
	}
	// Deep copy newTiles to d.Tiles.
	for y := 0; y < d.Height; y++ {
		for x := 0; x < d.Width; x++ {
			d.Tiles[y][x] = newTiles[y][x]
		}
	}
}

// CountAliveNeighbors counts the number of neighbors around point (x, y)
func (d *Dungeon) CountAliveNeighbors(x, y int) int {
	count := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			nx := x + i
			ny := y + j
			if j == 0 && i == 0 {
				continue
				// Check if ny and nx are outside the dimensions of d.Tiles.
			} else if ny < 0 || nx < 0 || ny >= d.Height || nx >= d.Width {
				count += 1
			} else if d.Tiles[ny][nx] == '#' {
				count += 1
			}
		}
	}
	return count
}

// FillBoarder turns all the floors on the border tiles into walls.
func (d *Dungeon) FillBoarder() {
	for y := 0; y < d.Height; y++ {
		for x := 0; x < d.Width; x++ {
			if y == 0 || y == d.Height-1 {
				if d.Tiles[y][x] == ' ' {
					d.Tiles[y][x] = '#'
				}
			} else {
				if x == 0 || x == d.Width-1 {
					if d.Tiles[y][x] == ' ' {
						d.Tiles[y][x] = '#'
					}
				}
			}
		}
	}
}

// FloodFill finds the number of tiles in the center chamber, and removes all floor tiles not in the center chamber.
func (d *Dungeon) FloodFill() int {
	// Create new set of tiles.
	newTiles := make([][]byte, d.Height)
	for i := 0; i < d.Height; i++ {
		newTiles[i] = make([]byte, d.Width)
		for j := 0; j < d.Width; j++ {
			newTiles[i][j] = '#'
		}
	}

	visited := make(map[Point]bool)
	q := []Point{}

	// Find a floor tile in the middle of the tiles.
	for y := d.Height/2 - 5; y < d.Height/2+5; y++ {
		for x := d.Width/2 - 5; x < d.Width/2+5; x++ {
			if d.Tiles[y][x] == ' ' {
				q = append(q, Point{x, y})
				break
			}
		}
	}
	// If not found, return -1; no chamber in the center.
	if len(q) == 0 {
		return -1
	}

	// Flood Fill algorithm.
	size := 0
	for len(q) > 0 {
		// Remove a point from the queue and if it has been visited, skip it.
		n := q[0]
		q = q[1:]
		if visited[n] {
			continue
		}
		visited[n] = true
		newTiles[n.y][n.x] = ' '
		// If a neighbor to the point is a floor, add it to the queue.
		if n.x > 0 && d.Tiles[n.y][n.x-1] == ' ' {
			q = append(q, Point{n.x - 1, n.y})
			size += 1
		}
		if n.x < d.Width-1 && d.Tiles[n.y][n.x+1] == ' ' {
			q = append(q, Point{n.x + 1, n.y})
			size += 1
		}
		if n.y > 0 && d.Tiles[n.y-1][n.x] == ' ' {
			q = append(q, Point{n.x, n.y - 1})
			size += 1
		}
		if n.y < d.Height-1 && d.Tiles[n.y+1][n.x] == ' ' {
			q = append(q, Point{n.x, n.y + 1})
			size += 1
		}
	}

	// Deep copy nezTiles to the d.Tiles:
	for y := 0; y < d.Height; y++ {
		for x := 0; x < d.Width; x++ {
			d.Tiles[y][x] = newTiles[y][x]
		}
	}

	return size
}

// Print prints the dungeon.
func (d *Dungeon) Print() {
	for y := 0; y < d.Height; y++ {
		for x := 0; x < d.Width; x++ {
			fmt.Printf("%c", d.Tiles[y][x])
		}
		fmt.Printf("\n")
	}
}
