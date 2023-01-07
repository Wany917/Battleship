package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	// map dimension
	rows     = 10
	cols     = 10
	numShips = 4
	minSize  = 2
	maxSize  = 5

	// color
	RESET = "\033[0m"  /* no color for reset the color*/
	RED   = "\033[31m" /* Red color for Sunken ships*/
	GREEN = "\033[32m" /* Green color for ships*/
	CYAN  = "\033[36m" /* Cyan color for watter*/
)

type Coordinate struct {
	// for ships Coordinate
	row int
	col int
}

type Ship struct {
	start    Coordinate
	end      Coordinate
	boat     []rune
	isSunken bool
}

type Playermap_s struct {
	playermap [][]rune
	Ships     []Ship
}

func generateMap() Playermap_s {
	// Initialize the map with all ~
	var playermap Playermap_s
	playermap.playermap = make([][]rune, rows)
	for i := range playermap.playermap {
		playermap.playermap[i] = make([]rune, cols)
		for j := 0; j != cols; j++ {
			playermap.playermap[i][j] = '~'
		}
	}
	// Place the ships on the playermap.playermap
	playermap.Ships = make([]Ship, numShips)
	rand.Seed(time.Now().UnixNano())
	for size, i := 3, 0; i < numShips; i++ {
		if i == -1 {
			size--
		}
		for {
			//size := minSize + rand.Intn(maxSize-minSize+1)
			isHorizontal := rand.Intn(2) == 0
			startRow := rand.Intn(rows)
			startCol := rand.Intn(cols)
			endRow := startRow
			endCol := startCol
			if isHorizontal {
				endCol += size - 1
				if endCol >= cols {
					continue // if out of set it restat boucle
				}
			} else {
				endRow += size - 1
				if endRow >= rows {
					continue // if out of set it restat boucle
				}
			}

			// Check if the area is clear
			// make sure that they are not another boat in the place that i want plcae the ships
			isClear := true
			for row := startRow; row <= endRow; row++ {
				for col := startCol; col <= endCol; col++ {
					if playermap.playermap[row][col] != '~' {
						isClear = false
						break
					}
				}
				if !isClear {
					break
				}
			}
			if !isClear {
				continue
			}

			// create ships: start cooordinate, end coordinate and a string boat
			playermap.Ships[i] = Ship{
				start: Coordinate{startRow, startCol},
				end:   Coordinate{endRow, endCol},
			}
			playermap.Ships[i].boat = make([]rune, size, size)
			for j := 0; j != size; j++ {
				playermap.Ships[i].boat[j] = '0'
			}
			playermap.Ships[i].isSunken = false
			// Place the ship on the playermap.playermap
			for row := startRow; row <= endRow; row++ {
				for col := startCol; col <= endCol; col++ {
					playermap.playermap[row][col] = '0'
				}
			}

			size++ // for size from 2 to five;
			break
		}
	}
	// Return player map structure with all ship and a map --> get a look to Playermap_s struct
	return playermap
}

func DisplayPlayerMap(user player, me bool) {
	str1 := "  | A  B  C  D  E  F  G  H  I  J"
	str2 := "--+-----------------------------"

	// print player name
	for i := 0; i != (len(str1)/2 - len(user.playerName)/2); i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("%s\n", user.playerName)

	// display the map
	fmt.Printf("%s\n%s\n", str1, str2)
	for i := 0; i < 10; i++ {
		fmt.Printf("%d |", i)
		for j := 0; j < 10; j++ {
			if user.playermap.playermap[i][j] == '0' { // if ships print with green color
				if me == true {
					fmt.Printf("%s[%c]%s", GREEN, user.playermap.playermap[i][j], RESET)
				} else {
					fmt.Printf("%s[%c]%s", RED, '~', RESET)
				}
			} else if user.playermap.playermap[i][j] == '*' { // if hit ships print with red color
				fmt.Printf("%s[%c]%s", RED, user.playermap.playermap[i][j], RESET)
			} else if user.playermap.playermap[i][j] == '~' { // if watter print with cyan color
				fmt.Printf("%s[%c]%s", CYAN, user.playermap.playermap[i][j], RESET)
			} else { // and print missed cell with no color
				fmt.Printf("[%c]", user.playermap.playermap[i][j])
			}
		}
		fmt.Printf("\n")
	}
}

func updateMap(playermap *Playermap_s) {
	// Reset the map to all av ~
	for i := range playermap.playermap {
		for j := 0; j != cols; j++ {
			if playermap.playermap[i][j] != ' ' {
				playermap.playermap[i][j] = '~'
			}
		}
	}

	// Place the ships on the map
	for i := 0; i < len(playermap.Ships); i++ {
		ship := playermap.Ships[i]
		len := 0
		for row := ship.start.row; row <= ship.end.row; row++ {
			for col := ship.start.col; col <= ship.end.col; col++ {
				playermap.playermap[row][col] = rune(ship.boat[len])
				len++
			}
		}
	}
}

func moveShips(playermap *Playermap_s) {
	for i := 0; i < len(playermap.Ships); i++ {
		ship := &playermap.Ships[i]

		// Generate a random direction (0: up, 1: down, 2: left, 3: right)
		direction := rand.Intn(4)
		// Update the position of the ship based on the direction
		if direction == 0 {
			// Move the ship up
			if ship.start.row > 0 {
				ship.start.row--
				ship.end.row--
			}
		} else if direction == 1 {
			// Move the ship down
			if ship.end.row < rows-1 {
				ship.start.row++
				ship.end.row++
			}
		} else if direction == 2 {
			// Move the ship left
			if ship.start.col > 0 {
				ship.start.col--
				ship.end.col--
			}
		} else {
			// Move the ship right
			if ship.end.col < cols-1 {
				ship.start.col++
				ship.end.col++
			}
		}
	}
	updateMap(playermap)
}

func shikShip(playermap *Playermap_s) {
	for i := 0; i != len(playermap.Ships); i++ {
		for j := 0; j != len(playermap.Ships[i].boat); j++ {
			if playermap.Ships[i].boat[j] == '0' {
				playermap.Ships[i].isSunken = false
				break
			}
			playermap.Ships[i].isSunken = true
		}
	}
}

func Display_All_Board(playerLits []player) {
	DisplayPlayerMap(playerList[0], true)
	for i := 1; i != len(playerLits); i++ {
		DisplayPlayerMap(playerList[i], false)
	}
}
