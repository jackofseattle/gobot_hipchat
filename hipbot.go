package main

import (
	"fmt"
	"github.com/jackofseattle/gobot_hipchat/lib"
	"os"
)

func main() {
	user := os.Getenv("GOBOT_JABBER_USER_ID")
	password := os.Getenv("GOBOT_JABBER_PASSWORD")

	robot := new(lib.Robot)

	robot.Name = "Go Bot"
	robot.Alias = "GoBot"

	err := robot.Connect(user, password)

	if err != nil {
		fmt.Printf("Error connecting: %s", err)
		return
	}

	robot.JoinAllAvailableRooms()

	fmt.Println("Beginning Listen Loop")
	robot.StartListening()
}
