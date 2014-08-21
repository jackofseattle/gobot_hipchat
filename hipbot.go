package main

import (
	"github.com/jackofseattle/gobot_hipchat/lib"
	"github.com/jackofseattle/gobot_hipchat/listeners"
	"log"
	"os"
)

func main() {

	user := os.Getenv("GOBOT_JABBER_USER_ID")
	password := os.Getenv("GOBOT_JABBER_PASSWORD")

	robot := lib.Robot{Name: "Go Bot", Alias: "GoBot"}

	err := robot.Connect(user, password)

	if err != nil {
		log.Fatalf("Error connecting: %s", err)
		return
	}

	registerListeners(&robot)

	log.Println("Beginning Listen Loop")

	robot.StartListening()
}

func registerListeners(robot *lib.Robot) {
	robot.Listen(listeners.PingResponder{robot})
	robot.Listen(listeners.AnimateResponder{robot})
}
