package lib

import (
	"errors"
	"fmt"
	"github.com/daneharrigan/hipchat"
	"strings"
)

type Listener interface {
	// A test to see if the handler for this instance should be called.
	Test(string) bool

	// Inputs are body, user object, roomid
	Handler(string, hipchat.User, string)
}

type Robot struct {
	listeners []Listener

	Name        string
	MentionName string
	Alias       string

	client *hipchat.Client
}

func (robot *Robot) Connect(jabberId string, password string) error {
	client, err := hipchat.NewClient(jabberId, password, "gobot")
	robot.client = client

	go client.KeepAlive()
	client.Status("chat")
	robot.JoinAllAvailableRooms()
	fmt.Printf("Connected to all available rooms")
	return err
}

// Joins all discoverable rooms on the connected server.
// NOTE: for testing this only joins a single room.
func (robot *Robot) JoinAllAvailableRooms() {
	for _, room := range robot.client.Rooms() {
		if room.Name == "BotTesting" {
			robot.client.Join(room.Id, "Go Bot")
			fmt.Println("Joined Room")
		}
	}
}

// Adds a listener to the robot.
func (robot *Robot) ListenFor(l Listener) {
	robot.listeners = append(robot.listeners, l)
}

func (robot *Robot) Say(roomId string, message string) {
	robot.client.Say(roomId, robot.Name, message)
}

// The standard listening loop. Gathers messages and meta data as they are received and calls all listeners.
// This is a blocking loop.
func (robot *Robot) StartListening() {
	for msg := range robot.client.Messages() {

		user, err := robot.GetUserObject(msg.From)

		if err != nil {
			fmt.Printf("Error: %s", err)
		} else {

			fmt.Printf("Message: (to: %s)\n", msg.To)
			fmt.Printf("  From: %s \n", user.Name)
			fmt.Printf("  Body: %s\n\n", msg.Body)
		}
	}
}

// The 'from' field for jabber messages comes back as a url with the name as the last path.
func (robot *Robot) GetUserObject(from string) (*hipchat.User, error) {
	split := strings.Split(from, "/")
	if len(split) != 2 {
		return new(hipchat.User), errors.New("Unable to parse a name from " + from)
	}
	for _, user := range robot.client.Users() {
		if user.Name == split[1] {
			return user, nil
		}
	}
	return new(hipchat.User), errors.New("Unable to find user named " + from)
}
