package listeners

import (
	"encoding/json"
	"github.com/daneharrigan/hipchat"
	"github.com/jackofseattle/gobot_hipchat/lib"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
)

const googleImagesEndpoint = "http://ajax.googleapis.com/ajax/services/search/images"

type gImageResult struct {
	Result string `json:"unescapedUrl"`
}

type gResponseData struct {
	Results []gImageResult `json:"results"`
}

type gAjaxReply struct {
	Data gResponseData `json:"responseData"`
}

// AnimateResponder - A responder to handle requests for 'animate' and 'image'. Fetches a random result from the google images
// api and returns the url to the browser.
type AnimateResponder struct {
	Robot *lib.Robot
}

// Test checks to see that the incoming requests matches the animate/image (me) <query> pattern
func (r AnimateResponder) Test(input string) (bool, map[string]string) {
	cmp := lib.NamedRegexp{regexp.MustCompile(`^(?P<type>animate|image)\s(me\s)?(?P<query>.+)`)}
	res := cmp.FindStringSubmatchMap(input)
	return len(res) > 0, res
}

// Handler fetches an image from the google images api. With an additional query parameter for animated requests.
func (r AnimateResponder) Handler(body string, user *hipchat.User, roomID string, params map[string]string) {
	query := make(url.Values)
	query.Set("v", "1.0")
	query.Set("rsz", "8")
	query.Set("safe", "active")

	q, _ := params["query"]
	query.Set("q", q)

	if t, _ := params["type"]; t == "animate" {
		query.Set("imgtype", "animated")
	}

	r.Robot.Say(roomID, r.getImage(query))
}

func (r AnimateResponder) getImage(query url.Values) string {
	u, _ := url.Parse(googleImagesEndpoint)

	u.RawQuery = query.Encode()
	log.Println(u)

	response, err := http.Get(u.String())

	if err != nil {
		log.Fatal(err)
		return ""
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	var dat gAjaxReply

	json.Unmarshal(contents, &dat)

	if &dat.Data == nil {
		log.Printf("Failed on data, %+v - %+v \n %+v", dat.Data, &dat.Data, dat)
		return ""
	}

	item := dat.Data.Results[rand.Intn(len(dat.Data.Results)-1)]
	return item.Result
}
