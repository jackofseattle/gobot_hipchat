package listeners

import (
	"github.com/daneharrigan/hipchat"
	"github.com/jackofseattle/gobot_hipchat/lib"
)

type SimpleResponder struct {
	Robot *lib.Robot
}

func (r SimpleResponder) Test(input string) (bool, map[string]string) {
	return false, make(map[string]string)
}

func (r SimpleResponder) Handler(body string, user *hipchat.User, roomId string, params map[string]string) {
	r.Robot.Say(roomId, "I got your message")
}
