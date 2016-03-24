package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Feed represents the nhl json object
type Feed struct {
	Key      int    `json:"gamePk"`
	Link     string `json:"link"`
	GameData struct {
		Date struct {
			StartDate string `json:"dateTime"`
			EndDate   string `json:"endDateTime"`
		} `json:"datetime"`
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
				Home team `json:"home"`
				Away team `json:"away"`
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

type team struct {
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

// Event to insert in the database
type Event struct {
	Type    string
	Media   string
	MatchID int
	Score   string
}

func main() {

	var prevFeed Feed

	for {

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

		// TODO: Made these 4 methods into interface that each sport will have to implement.
		// The previous state wasn't live and now it is -- Record game started event.
		if prevFeed.GameData.Status.DetailedState != "Live" && feed.GameData.Status.DetailedState == "Live" {
			//TODO: Need have this match present in a match table
			insertMessageToSend(db, Event{"GameStarted", "", feed.Key, "0-0"})
		}

		// Check for new goals -- Seperate for each team?
		if prevFeed.LiveData.LineScore.Teams.Home.Goals != feed.LiveData.LineScore.Teams.Home.Goals ||
			prevFeed.LiveData.LineScore.Teams.Away.Goals != feed.LiveData.LineScore.Teams.Away.Goals {

			score := fmt.Printf("%d - %d", feed.LiveData.LineScore.Teams.Home.Goals, feed.LiveData.LineScore.Teams.Away.Goals)
			insertMessageToSend(db, Event{"Goal", "", feed.Key, score})
		}

		if prevFeed.LiveData.LineScore.CurrentPeriod > 1 && prevFeed.LiveData.LineScore.CurrentPeriod != feed.LiveData.LineScore.CurrentPeriod {
			score := fmt.Printf("%d - %d", feed.LiveData.LineScore.Teams.Home.Goals, feed.LiveData.LineScore.Teams.Away.Goals)
			insertMessageToSend(db, Event{"End of period", "", feed.Key, score})
		}

		if prevFeed.GameData.Status.DetailedState == "Live" && feed.GameData.Status.DetailedState == "Final" {
			score := fmt.Printf("%d - %d", feed.LiveData.LineScore.Teams.Home.Goals, feed.LiveData.LineScore.Teams.Away.Goals)
			insertMessageToSend(db, Event{"GameEnded", "", feed.Key, score})
		}

		prevFeed = feed

	}

}

func insertMessageToSend(db *sql.DB, Event newEvent) {

	stmNewOutbox, err := db.Prepare("INSERT INTO `ScoreBot`.`Event` (`Type`,`Media`,`MatchId`,`Score`) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	defer stmNewOutbox.Close()
	_, err = stmNewOutbox.Exec(newEvent.Type, newEvent.Media, newEvent.MatchID, newEvent.Score)
	if err != nil {
		panic(err.Error())
	}
}

// func getUnReadMessages(db *sql.DB) (incomingMessage, error) {
// 	stmtOut, err := db.Prepare("SELECT id, author_id, thread_id, author_name, message, isRead, timeReceived FROM inbox WHERE isRead = ?")
// 	if err != nil {
// 		panic(err.Error())
// 	}
//
// 	defer stmtOut.Close()
//
// 	var message incomingMessage
//
// 	errNoRow := stmtOut.QueryRow(0).Scan(&message.ID, &message.AuthorID, &message.ThreadID, &message.AuthorName, &message.Message, &message.isRead)
// 	// if err != nil {
// 	// 	panic(err.Error())
// 	// }
//
// 	return message, errNoRow
// }
//
// func toggleIsRead(db *sql.DB, message incomingMessage) {
// 	stmUpdateInbox, err := db.Prepare("update inbox set isRead = 1 where id = ?;")
// 	if err != nil {
// 		panic(err.Error())
// 	}
//
// 	defer stmUpdateInbox.Close()
// 	_, err = stmUpdateInbox.Exec(message.ID)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// }
//
