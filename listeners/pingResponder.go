package listeners

import (
	"github.com/daneharrigan/hipchat"
	"github.com/jackofseattle/gobot_hipchat/lib"
	"regexp"
)

// PingResponder is a simple responder that replies to 'ping' with 'pong'
type PingResponder struct {
	Robot *lib.Robot
}

// Test checks the incoming string to see that it starts with ping
func (r PingResponder) Test(input string) (bool, map[string]string) {
	match, _ := regexp.MatchString(`^ping\s*`, input)
	return match, make(map[string]string)
}

// Handler responds to the room with 'pong'
func (r PingResponder) Handler(body string, user *hipchat.User, roomID string, params map[string]string) {
	r.Robot.Say(roomID, "pong")
}
