package main

import (
	"github.com/nicolas-martin/scoreUpdate/Interfaces"
	"github.com/nicolas-martin/scoreUpdate/Sports"
)

func main() {
	//TODO: Make crawler to get all the game ID, teams and datetime
	nhl := Sports.Nhl{Game: Interfaces.Game{URL: "https://statsapi.web.nhl.com/api/v1/game/2015021078/feed/live"}}
	nhl.Loop()

}
