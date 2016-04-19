package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nicolas-martin/scoreUpdate/Interfaces"
	"github.com/nicolas-martin/scoreUpdate/Sports"
)

func main() {
	games := getGames()

	for _, game := range games {
		gameDate, _ := time.Parse("2006-01-02T15:04:05Z07:00", game.Start)

		if time.Now().Unix() >= gameDate.Unix() {

			nhl := Sports.Nhl{Game: Interfaces.Game{ID: game.GameID, URL: game.URL}}

			//TODO: Start a new thread here
			//TODO: Change the game status to "in progress" and modify the loop to avoid them
			nhl.Loop()

		}

	}

}

// TODO: move this in a sql package
func createDbConn() *sql.DB {
	db, err := sql.Open("mysql", "root:aiwojefoa39j2a9VVA3jj32fa3@cloudsql(sportsbot-1255:us-east1:sportsupdate)/ScoreBot")

	if err != nil {
		panic(err.Error())
	}

	return db
}

func getGames() []game {

	db := createDbConn()
	defer db.Close()

	sqlQuery := "SELECT `Games`.`GameId`, `Games`.`AwayId`, `Games`.`Start`, `Games`.`Finish`, `Games`.`HomeScore`, `Games`.`AwayScore`, `Games`.`Status`, `Games`.`homeId`, `Games`.`url` FROM `ScoreBot`.`Games`"

	row, err := db.Query(sqlQuery)
	if err != nil {
		panic(err)
	}

	var gameList []game

	for row.Next() {
		g := game{}

		err := row.Scan(&g.GameID, &g.AwayID, &g.Start, &g.Finish, &g.HomeScore, &g.AwayScore, &g.Status, &g.HomeID, &g.URL)

		if err != nil {
			fmt.Println(err)
		}

		gameList = append(gameList, g)
	}

	return gameList
}

type game struct {
	GameID    int
	AwayID    int
	Start     string
	Finish    string
	HomeScore int
	AwayScore int
	Status    string
	HomeID    int
	URL       string
}
