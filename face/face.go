package face

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/shanghuiyang/go-speech/oauth"
)

const (
	baiduURL = "https://aip.baidubce.com/rest/2.0/face/v3/search"
)

// Face ...
type Face struct {
	auth *oauth.Oauth
}

type response struct {
	ErrorCode int     `json:"error_code"`
	ErrorMsg  string  `json:"error_msg"`
	LogID     int64   `json:"log_id"`
	Timestamp int64   `json:"timestamp"`
	Cached    int64   `json:"cached"`
	Result    *result `json:"result"`
}

// Result ...
type result struct {
	FaceToken string  `json:"face_token"`
	UserList  []*User `json:"user_list"`
}

// User ...
type User struct {
	GroupID  string  `json:"group_id"`
	UserID   string  `json:"user_id"`
	UserInfo string  `json:"user_info"`
	Score    float64 `json:"score"`
}

// New ...
func New(auth *oauth.Oauth) *Face {
	return &Face{
		auth: auth,
	}
}

// Recognize ...
func (f *Face) Recognize(imageFile, groupID string) ([]*User, error) {
	token, err := f.auth.GetToken()
	if err != nil {
		return nil, err
	}

	b64img, err := f.b64Image(imageFile)
	if err != nil {
		return nil, err
	}

	formData := url.Values{
		"access_token":  {token},
		"image":         {b64img},
		"image_type":    {"BASE64"},
		"group_id_list": {groupID},
	}
	resp, err := http.PostForm(baiduURL, formData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	if res.ErrorCode > 0 {
		return nil, fmt.Errorf("error_code: %v, error_msg: %v", res.ErrorCode, res.ErrorMsg)
	}
	if res.Result == nil {
		return nil, fmt.Errorf("not found")
	}
	return res.Result.UserList, nil
}

func (f *Face) b64Image(imageFile string) (string, error) {
	file, err := os.Open(imageFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	image, err := ioutil.ReadAll(file)
	if err != nil {
		return "nil", err
	}
	b64img := base64.StdEncoding.EncodeToString(image)
	return b64img, nil
}
