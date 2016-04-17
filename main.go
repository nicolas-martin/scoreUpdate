package main

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nicolas-martin/scoreUpdate/Interfaces"
	"github.com/nicolas-martin/scoreUpdate/Sports"
)

func main() {
	games := getGames()

	for _, game := range games {
		gameDate, _ := time.Parse("2006-01-02T15:04:05Z07:00", game.Start)

		if time.Now() >= gameDate {

			nhl := Sports.Nhl{Game: Interfaces.Game{URL: game.URL}}
			// go nhl.Loop()
			nhl.Loop()

		}

	}

}

func getGames() []game {
	//TODO: Put this somewhere else.
	db, err := sql.Open("mysql", "root:password@/ScoreBot")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	sqlQuery := "SELECT `Games`.`AwayId, `Games`.`Start, `Games`.`Finish, `Games`.`HomeScore, `Games`.`AwayScore, `Games`.`Status, `Games`.`homeId`, `Games`.`url` FROM `ScoreBot`.`Games`"

	row, err := db.Query(sqlQuery)
	if err != nil {
		panic(err)
	}

	var gameList []games

	for row.Next() {
		g := game{}

		err = row.Scan(&g.AwayID, &g.Start, &g.Finish, &g.HomeScore, &g.AwayScore, &g.Status, &g.HomeID, &g.URL)

		gameList = append(gameList, g)
	}

	return userList
}

type game struct {
	AwayID    int
	Start     string
	Finish    string
	HomeScore int
	AwayScore int
	Status    string
	HomeID    int
	URL       string
}
