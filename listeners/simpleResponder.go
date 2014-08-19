package listeners

import (
	"github.com/daneharrigan/hipchat"
	"github.com/jackofseattle/gobot_hipchat/lib"
)

type SimpleResponder struct {
	robot *lib.Robot
}

func (r *SimpleResponder) Test(input string) bool {
	return true
}

func (r *SimpleResponder) Handler(body string, user *hipchat.User, roomId string) {
	r.robot.Say(roomId, "I got your message")
}
