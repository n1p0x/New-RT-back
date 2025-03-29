package main

import "time"

var (
	RoundDuration = 3 * time.Minute
	MinPlayers    = 2
)

func main() {
	runGame()
}
