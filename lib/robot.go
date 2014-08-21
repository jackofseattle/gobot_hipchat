package lib

import (
	"github.com/daneharrigan/hipchat"
	"log"
	"strings"
	"time"
)

type Listener interface {
	// A test to see if the handler for this instance should be called.
	Test(string) (bool, map[string]string)

	// Inputs are body, user object, roomid
	Handler(string, *hipchat.User, string, map[string]string)
}

type Robot struct {
	listeners []Listener

	Name        string
	MentionName string
	Alias       string
	UserList    map[string]*hipchat.User

	client        *hipchat.Client
	initialLoaded bool
}

func (robot *Robot) Connect(jabberId string, password string) error {
	client, err := hipchat.NewClient(jabberId, password, "gobot")
	robot.client = client

	client.Status("chat")
	robot.JoinAllAvailableRooms()
	log.Printf("Connected to all available rooms")
	robot.initialLoaded = false

	robot.UserList = make(map[string]*hipchat.User)
	go robot.CollectUserObjects()

	go client.KeepAlive()
	return err
}

// Joins all discoverable rooms on the connected server.
// NOTE: for testing this only joins a single room.
func (robot *Robot) JoinAllAvailableRooms() {
	for _, room := range robot.client.Rooms() {
		if room.Name == "BotTesting" {
			robot.client.Join(room.Id, "Go Bot")
			log.Println("Joined Room")
		}
	}
}

// Adds a new listener to the robot.
func (robot *Robot) Listen(l Listener) {
	robot.listeners = append(robot.listeners, l)
}

func (robot *Robot) Say(roomId string, message string) {
	robot.client.Say(roomId, robot.Name, message)
}

// The standard listening loop. Gathers messages and meta data as they are received and calls all interested listeners.
// This is a blocking loop.
//
// The loop is skipped for the first 3 seconds after calling. This allows for all of the history items to be cleared
// before we start handling messages. Further, the loop is skipped if the user list hasn't been loaded yet.
// Messages from the bot or an undefined user are ignored.
func (robot *Robot) StartListening() {
	go robot.deferMessageReception()

	for msg := range robot.client.Messages() {
		if len(robot.UserList) == 0 || !robot.initialLoaded {
			continue
		}

		userName := robot.getUserName(msg.From)
		if userName == "" {
			log.Printf("Message received without user name (From: %s)", msg.From)
			continue
		}

		user, ok := robot.UserList[userName]

		if !ok {
			log.Printf("No user object found for %s (From: %s)", userName, msg.From)
			continue
		}

		if user.Id == robot.client.Id {
			log.Printf("message from robot, discarding")
			continue
		}

		testBody, ok := robot.getBotMessage(msg.Body)
		if !ok {
			continue
		}

		for _, l := range robot.listeners {
			if ok, params := l.Test(testBody); ok {
				go l.Handler(testBody, user, msg.From, params)
			}
		}
	}
}

// Collectes all the users on the server for later lookup.
// We do this here because fetching the users can be a timeconsuming process.
func (robot *Robot) CollectUserObjects() {
	for _, user := range robot.client.Users() {
		robot.UserList[user.Name] = user
	}
}

func (robot *Robot) getUserName(messageUrl string) string {
	split := strings.Split(messageUrl, "/")

	if len(split) != 2 {
		return ""
	}
	return split[1]
}

func (robot *Robot) getBotMessage(body string) (string, bool) {
	botAliasLen := len(robot.Alias)

	if botAliasLen >= len(body) {
		return "", false
	}

	if strings.ToLower(body[:botAliasLen]) != strings.ToLower(robot.Alias) {
		return "", false
	}

	return strings.TrimSpace(body[botAliasLen+1:]), true
}

// Hipchat loves to send along the entire history of the room as a series of rapid messages.
// Rather than trying to figure which of the messages is a history message and which is new, we'll just ignore all
// messages for the first 3 seconds.
func (robot *Robot) deferMessageReception() {
	time.Sleep(time.Second * 3)
	robot.initialLoaded = true
	log.Printf("now receiving messages \n")
}
