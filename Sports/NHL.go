package Sports

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		//https://statsapi.web.nhl.com/api/v1/game/2015021078/feed/live"}}
		apiURL := fmt.Sprintf("https://statsapi.web.nhl.com/%s", n.URL)
		resp, err := http.Get(apiURL)

		if err != nil {
			fmt.Println(err)
		}

		var feed = new(Feed)

		err = json.NewDecoder(resp.Body).Decode(&feed)

		if err != nil {
			fmt.Println(err)
		}

		db := createDbConn()
		// First run
		if prevFeed.Key == 0 {
			prevFeed = *feed
		}

		defer db.Close()
		var newEvent Event

		// The previous state wasn't live and now it is -- Record game started event.
		var strLive = "Live"
		if prevFeed.GameData.Status.DetailedState != strLive && feed.GameData.Status.DetailedState == strLive {
			newEvent = Event{"GameStarted", feed.Key}
		}

		// Check for new goals -- Seperate for each team?
		if prevFeed.LiveData.LineScore.Teams.Home.Goals != feed.LiveData.LineScore.Teams.Home.Goals ||
			prevFeed.LiveData.LineScore.Teams.Away.Goals != feed.LiveData.LineScore.Teams.Away.Goals {
			description := fmt.Sprintf("Goal %s", feed.score())
			newEvent = Event{description, feed.Key}
		}

		if prevFeed.LiveData.LineScore.CurrentPeriod > 1 && prevFeed.LiveData.LineScore.CurrentPeriod != feed.LiveData.LineScore.CurrentPeriod {
			description := fmt.Sprintf("End of period %s", feed.score())
			newEvent = Event{description, feed.Key}
		}

		if prevFeed.GameData.Status.DetailedState == strLive && feed.GameData.Status.DetailedState == "Final" {
			description := fmt.Sprintf("GameEnded %s", feed.score())
			insertMessageToSend(db, Event{description, feed.Key})
			sendRequestToChat(&newEvent, feed.LiveData.LineScore.Teams.Home.InTeam.ID, feed.LiveData.LineScore.Teams.Away.InTeam.ID)
			break
		}

		if newEvent.GameID != 0 {

			insertMessageToSend(db, newEvent)
			sendRequestToChat(&newEvent, feed.LiveData.LineScore.Teams.Home.InTeam.ID, feed.LiveData.LineScore.Teams.Away.InTeam.ID)
		}

		defer db.Close()

		prevFeed = *feed
		time.Sleep(2 * time.Second)

	}
}

func insertMessageToSend(db *sql.DB, newEvent Event) {

	stmNewOutbox, err := db.Prepare("INSERT INTO `ScoreBot`.`Event` (`Description`,`GameId`) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}

	defer stmNewOutbox.Close()
	_, err = stmNewOutbox.Exec(newEvent.Description, newEvent.GameID)

	if err != nil {
		panic(err.Error())
	}
}

func getAllUsersForTeam(teamID int) []user {
	url := fmt.Sprintf("https://sportsbot-1255.appspot.com/User/%v", teamID)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var users []user

	if err := json.Unmarshal(body, &users); err != nil {
		fmt.Println(err)
	}

	return users

}

//TODO: use the ChatBot package
func sendRequestToChat(incomingEvent *Event, homeID int, awayID int) {
	var allUsers []user

	//Get all users that are subscribed to a team
	homeUsers := getAllUsersForTeam(homeID)
	awayUsers := getAllUsersForTeam(awayID)

	for _, user := range awayUsers {
		allUsers = append(allUsers, user)
	}

	for _, user := range homeUsers {
		allUsers = append(allUsers, user)
	}

	// TODO: Switch depending on platform
	url := "https://66jezutq1k.execute-api.us-east-1.amazonaws.com/production/v1/update/kik"

	newChatPost := chatPost{allUsers, incomingEvent.Description, "text"}

	buf, _ := json.Marshal(newChatPost)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
}

//TODO: Use the same struct package as the API
type user struct {
	Username string
	Platform string
	Phone    string
	Country  string
	Joined   string
	ChatID   string
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

// TODO: Move this to sql package
func createDbConn() *sql.DB {
	// db, err := sql.Open("mysql", "root:aiwojefoa39j2a9VVA3jj32fa3@cloudsql(sportsbot-1255:us-east1:sportsupdate)/ScoreBot")
	db, err := sql.Open("mysql", "root:password@/ScoreBot")

	if err != nil {
		panic(err.Error())
	}

	return db
}

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
	Description string
	GameID      int
}

type chatPost struct {
	Users []user
	Body  string
	Type  string
}
