package Sports

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicolas-martin/scoreUpdate/Interfaces"
	// Blank import due to its use as a driver
	_ "github.com/go-sql-driver/mysql"
)

func (f *Feed) score() string {
	return fmt.Sprintf("%d - %d \r\n", f.LiveData.LineScore.Teams.Home.Goals, f.LiveData.LineScore.Teams.Away.Goals)
}

// Nhl represents an NHL game
type Nhl struct {
	Interfaces.Game
}

// Loop and check for new events
func (n *Nhl) Loop() {
	var prevFeed Feed

	for {

		resp, err := http.Get(n.URL)

		if err != nil {
			fmt.Println(err)
		}

		feed := new(Feed)
		err = json.NewDecoder(resp.Body).Decode(feed)

		fmt.Printf(feed.score())
		fmt.Println(feed)

		db, err := sql.Open("mysql", "root:password@/ScoreBot")
		if err != nil {
			panic(err.Error())
		}

		// First run
		if prevFeed.Key == 0 {
			prevFeed = *feed
		}

		defer db.Close()

		_ = "breakpoint"
		// TODO: Made these 4 methods into interface that each sport will have to implement.
		// The previous state wasn't live and now it is -- Record game started event.
		var strLive = "Live"
		if prevFeed.GameData.Status.DetailedState != strLive && feed.GameData.Status.DetailedState == strLive {
			//TODO: Need have this match present in a match table
			insertMessageToSend(db, Event{"GameStarted", "", feed.Key, "0-0"})
		}

		// Check for new goals -- Seperate for each team?
		if prevFeed.LiveData.LineScore.Teams.Home.Goals != feed.LiveData.LineScore.Teams.Home.Goals ||
			prevFeed.LiveData.LineScore.Teams.Away.Goals != feed.LiveData.LineScore.Teams.Away.Goals {

			insertMessageToSend(db, Event{"Goal", "", feed.Key, feed.score()})
		}

		if prevFeed.LiveData.LineScore.CurrentPeriod > 1 && prevFeed.LiveData.LineScore.CurrentPeriod != feed.LiveData.LineScore.CurrentPeriod {
			insertMessageToSend(db, Event{"End of period", "", feed.Key, feed.score()})
		}

		if prevFeed.GameData.Status.DetailedState == strLive && feed.GameData.Status.DetailedState == "Final" {
			insertMessageToSend(db, Event{"GameEnded", "", feed.Key, feed.score()})
			break
		}

		prevFeed = *feed
		time.Sleep(2 * time.Second)

	}
}

func insertMessageToSend(db *sql.DB, newEvent Event) {

	stmNewOutbox, err := db.Prepare("INSERT INTO `ScoreBot`.`Event` (`Type`,`Media`,`MatchId`,`Score`, `IsSent`) VALUES (?, ?, ?, ?, 0)")
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
