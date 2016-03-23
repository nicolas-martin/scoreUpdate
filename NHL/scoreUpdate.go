package main

import (
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

	// db, err := sql.Open("mysql", "root:password@/ScoreBot")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()
	//
	// for {
	//
	// }

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
// func insertMessageToSend(db *sql.DB, message incomingMessage, returnMessage string) {
// 	stmNewOutbox, err := db.Prepare("INSERT INTO `outbox` (`author_id`, `thread_id`, `author_name`, `message`) values (?, ?, ?, ?)")
// 	if err != nil {
// 		panic(err.Error())
// 	}
//
// 	defer stmNewOutbox.Close()
// 	fmt.Printf("returnMessage = %s \r\n", returnMessage)
// 	_, err = stmNewOutbox.Exec(message.AuthorID, message.ThreadID, message.AuthorName, returnMessage)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// }
