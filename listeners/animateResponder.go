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
	GsearchResultClass  string
	width               string
	height              string
	imageId             string
	tbWidth             string
	tbHeight            string
	Result              string `json:"unescapedUrl"`
	url                 string
	visibleUrl          string
	title               string
	titleNoFormatting   string
	originalContextUrl  string
	content             string
	contentNoFormatting string
	tbUrl               string
}

type gResponseData struct {
	Results []gImageResult `json:"results"`
	cursor  interface{}
}

type gAjaxReply struct {
	Data            interface{} `json:"responseData"`
	responseDetails interface{}
	responseStatus  interface{}
}

type AnimateResponder struct {
	Robot *lib.Robot
}

func (r AnimateResponder) Test(input string) (bool, map[string]string) {

	cmp := lib.NamedRegexp{regexp.MustCompile(`^(?P<type>animate|image)\s(me\s)?(?P<query>.+)`)}
	res := cmp.FindStringSubmatchMap(input)
	return len(res) > 0, res
}

func (r AnimateResponder) Handler(body string, user *hipchat.User, roomId string, params map[string]string) {
	query := make(url.Values)
	query.Set("v", "1.0")
	query.Set("rsz", "8")
	query.Set("safe", "active")

	q, _ := params["query"]
	query.Set("q", q)

	if t, _ := params["type"]; t == "animate" {
		query.Set("imgtype", "animated")
	}

	r.Robot.Say(roomId, r.getImage(query))
}

func (r AnimateResponder) getImage(query url.Values) string {
	u, _ := url.Parse(googleImagesEndpoint)

	u.RawQuery = query.Encode()
	log.Println(u)

	response, err := http.Get(u.String())

	if err != nil {
		log.Fatal(err)
		return ""
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
			return ""
		}

		var dat gAjaxReply

		json.Unmarshal(contents, &dat)

		if dat.Data == nil {
			log.Printf("Failed on data, %+v - %+v \n %+v", dat.Data, &dat.Data, dat)
			return ""
		}

		k := dat.Data.(map[string]interface{})

		res, _ := k["results"]

		n := rand.Intn(len(res.([]interface{})) - 1)

		item := res.([]interface{})[n].(map[string]interface{})
		imageUrl := item["unescapedUrl"]

		log.Printf("contents: %+v", imageUrl)

		return imageUrl.(string)
	}

}
