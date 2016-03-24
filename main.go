package main

import (
	"github.com/nicolas-martin/scoreUpdate/Interfaces"
	"github.com/nicolas-martin/scoreUpdate/NHL"
)

func main() {
	nhl := NHL.Nhl{Interfaces.Game{URL: "https://statsapi.web.nhl.com/api/v1/game/2015021078/feed/live"}}
	nhl.Loop()

}
