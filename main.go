package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nicolas-martin/scoreUpdate/Interfaces"
	"github.com/nicolas-martin/scoreUpdate/Sports"
)

func main() {
	parseSchedule()
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

func parseSchedule() {

	//TODO: Select the entire season
	resp, err := http.Get("https://statsapi.web.nhl.com/api/v1/schedule?startDate=2016-04-24&endDate=2016-05-21")

	if err != nil {
		fmt.Println(err)
	}

	if resp.Body == nil {
		panic("ahhh")
	}

	var tSchedule schedule

	dec := json.NewDecoder(resp.Body)

	err2 := dec.Decode(&tSchedule)

	if err2 != nil {
		fmt.Println(err)
	}

	insertMessageToSchedule(&tSchedule)

}

func insertMessageToSchedule(pSchedule *schedule) {
	db := createDbConn()
	for i := 0; i < len(pSchedule.Dates); i++ {

		for j := 0; j < len(pSchedule.Dates[i].Games); j++ {

			// stmNewOutbox, err := db.Prepare("INSERT INTO `ScoreBot`.`Event` (`Type`,`Media`,`MatchId`,`Score`, `IsSent`) VALUES (?, ?, ?, ?, 0)")
			stmNewOutbox, err := db.Prepare("INSERT INTO `ScoreBot`.`Games` (`AwayId`, `Start`,`Finish`,`HomeScore`,`AwayScore`,`Status`,`homeId`, `url`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				// panic(err.Error())
				fmt.Println(err)
			}

			game := pSchedule.Dates[i].Games[j]

			defer stmNewOutbox.Close()
			gameDate, _ := time.Parse("2006-01-02T15:04:05Z07:00", game.GameDate)

			_, err = stmNewOutbox.Exec(game.Teams.Away.Team.ID, gameDate, gameDate, 0, 0, game.Status.DetailedState, game.Teams.Home.Team.ID, game.Link)
			if err != nil {
				panic(err.Error())
			}
		}

	}
}

// TODO: move this in a sql package
func createDbConn() *sql.DB {
    //104.196.10.95
	db, err := sql.Open("mysql", "root:aiwojefoa39j2a9VVA3jj32fa3@tcp(104.196.10.95:3306)/ScoreBot")
	// db, err := sql.Open("mysql", "root:password@/ScoreBot")

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

type schedule struct {
	Copyright string `json:"copyright"`
	Dates     []struct {
		Date       string         `json:"date"`
		Games      []scheduleGame `json:"games"`
		TotalItems int            `json:"totalItems"`
	} `json:"dates"`
	TotalItems int `json:"totalItems"`
	Wait       int `json:"wait"`
}

type scheduleGame struct {
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
