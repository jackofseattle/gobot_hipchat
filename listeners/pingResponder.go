package listeners

import (
	"github.com/daneharrigan/hipchat"
	"github.com/jackofseattle/gobot_hipchat/lib"
	"regexp"
)

type PingResponder struct {
	Robot *lib.Robot
}

func (r PingResponder) Test(input string) (bool, map[string]string) {
	match, _ := regexp.MatchString(`^ping\s*`, input)
	return match, make(map[string]string)
}

func (r PingResponder) Handler(body string, user *hipchat.User, roomId string, params map[string]string) {
	r.Robot.Say(roomId, "pong")
}
