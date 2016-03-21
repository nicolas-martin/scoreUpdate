package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Item post value
type Feed struct {
	Key      int    `json:"gamePk"`
	Link     string `json:"link"`
	GameData struct {
		Status struct {
			GameState     string `json:"abstractGameState"`
			DetailedState string `json:"detailedState"`
		} `json:"status"`
	} `json:"gameData"`
	LiveData struct {
		LineScore struct {
			CurrentPeriod         int    `json:"currentPeriod"`
			CurrentPeriodTimeLeft string `json:"currentPeriodTimeRemaining"`
			Teams                 struct {
				Home Team `json:"home"`
				Away Team `json:"away"`
			} `json:"teams"`
		} `json:"linescore"`
		BoxScore struct {
			Team struct {
			} `json:"teams"`
		} `json:"boxscore"`
		Plays struct {
			ScoringPlays []int `json:"scoringPlays"`
		} `json:"plays"`
	} `json:"liveData"`
}

type Team struct {
	// GoaliePulled bool `json:"goaliePulled"`
	Goals int `json:"goals"`
	// NumSkaters   int  `json:"numSkaters"`
	// PowerPlay    bool `json:"powerPlay"`
	ShotsOnGoal int `json:"shotsOnGoal"`
	InTeam      struct {
		ID   int    `json:"id"`
		Link string `json:"link"`
		Name string `json:"name"`
	} `json:"team"`
}

func main() {
	resp, err := http.Get("https://statsapi.web.nhl.com/api/v1/game/2015021078/feed/live")

	if err != nil {
		fmt.Println(err)
	}

	feed := new(Feed)
	err = json.NewDecoder(resp.Body).Decode(feed)
	fmt.Printf("%d - %d \r\n", feed.LiveData.LineScore.Teams.Home.Goals, feed.LiveData.LineScore.Teams.Away.Goals)
	fmt.Println(feed)

	db, err := sql.Open("mysql", "root:password@/ScoreBot")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	for {

	}

}
