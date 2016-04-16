package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//TODO: Make crawler to get all the game ID, teams and datetime
	// nhl := Sports.Nhl{Game: Interfaces.Game{URL: "https://statsapi.web.nhl.com/api/v1/game/2015021078/feed/live"}}
	// nhl.Loop()

	//https://statsapi.web.nhl.com/api/v1/schedule?startDate=2016-04-16&endDate=2016-04-21
	parseSchedule()
}

func parseSchedule() {

	//TODO: Put this somewhere else.
	db, err := sql.Open("mysql", "root:password@/ScoreBot")
	if err != nil {
		panic(err.Error())
	}

	//TODO: Select the entire season
	resp, err := http.Get("https://statsapi.web.nhl.com/api/v1/schedule?startDate=2016-04-16&endDate=2016-04-21")

	if err != nil {
		fmt.Println(err)
	}

	schedule := new(schedule)
	err = json.NewDecoder(resp.Body).Decode(schedule)

	if err != nil {
		fmt.Println(err)
	}

	insertMessageToSchedule(db, schedule)

}

func insertMessageToSchedule(db *sql.DB, schedule *schedule) {

	for i := 0; i < len(schedule.Dates); i++ {

		for j := 0; j < len(schedule.Dates[i].Games); j++ {

			// stmNewOutbox, err := db.Prepare("INSERT INTO `ScoreBot`.`Event` (`Type`,`Media`,`MatchId`,`Score`, `IsSent`) VALUES (?, ?, ?, ?, 0)")
			stmNewOutbox, err := db.Prepare("INSERT INTO `ScoreBot`.`Games` (`AwayId`, `Start`,`Finish`,`HomeScore`,`AwayScore`,`Status`,`homeId`, `url`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				panic(err.Error())
			}

			game := schedule.Dates[i].Games[j]

			_ = "breakpoint"

			defer stmNewOutbox.Close()
			gameDate, _ := time.Parse("2006-01-02T15:04:05Z07:00", game.GameDate)

			_, err = stmNewOutbox.Exec(game.Teams.Away.Team.ID, gameDate, gameDate, 0, 0, game.Status.DetailedState, game.Teams.Home.Team.ID, game.Link)
			if err != nil {
				panic(err.Error())
			}
		}

	}
}

type schedule struct {
	Copyright string `json:"copyright"`
	Dates     []struct {
		Date       string `json:"date"`
		Games      []game `json:"games"`
		TotalItems int    `json:"totalItems"`
	} `json:"dates"`
	TotalItems int `json:"totalItems"`
	Wait       int `json:"wait"`
}

type game struct {
	Content struct {
		Link string `json:"link"`
	} `json:"content"`
	GameDate string `json:"gameDate"`
	GamePk   int    `json:"gamePk"`
	GameType string `json:"gameType"`
	Link     string `json:"link"`
	Season   string `json:"season"`
	Status   struct {
		AbstractGameState string `json:"abstractGameState"`
		CodedGameState    string `json:"codedGameState"`
		DetailedState     string `json:"detailedState"`
		StatusCode        string `json:"statusCode"`
	} `json:"status"`
	Teams struct {
		Away struct {
			LeagueRecord struct {
				Losses int    `json:"losses"`
				Type   string `json:"type"`
				Wins   int    `json:"wins"`
			} `json:"leagueRecord"`
			Score int `json:"score"`
			Team  struct {
				ID   int    `json:"id"`
				Link string `json:"link"`
				Name string `json:"name"`
			} `json:"team"`
		} `json:"away"`
		Home struct {
			LeagueRecord struct {
				Losses int    `json:"losses"`
				Type   string `json:"type"`
				Wins   int    `json:"wins"`
			} `json:"leagueRecord"`
			Score int `json:"score"`
			Team  struct {
				ID   int    `json:"id"`
				Link string `json:"link"`
				Name string `json:"name"`
			} `json:"team"`
		} `json:"home"`
	} `json:"teams"`
	Venue struct {
		Name string `json:"name"`
	} `json:"venue"`
}
