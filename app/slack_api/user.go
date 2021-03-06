package slackAPI

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"projects/slack-api/config"
)

type SlackUser struct {
	Ok   bool `json:"ok"`
	User struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Profile struct {
			DisplayName string `json:"display_name"`
			Email       string `json:"email"`
		} `json:"profile"`
	} `json:"user"`
}

func GetUserInfo(userID string) *SlackUser {
	request, err := http.NewRequest(http.MethodGet, config.Config.SlackURL+"/users.info", nil)
	if err != nil {
		log.Println("failed to create a new request:", err)
		return &SlackUser{}
	}

	params := request.URL.Query()
	params.Set("token", config.Config.Token)
	params.Set("user", userID)
	request.URL.RawQuery = params.Encode()

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Println("failed to request for get user info: ", err)
		return &SlackUser{}
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("read request body failed: ", err)
		return &SlackUser{}
	}

	var user SlackUser
	if err = json.Unmarshal(body, &user); err != nil {
		log.Println("json unmarshal slack user failed: ", err)
		return &SlackUser{}
	}

	return &user
}

func (user *SlackUser) lookupUserByEmail() {
	request, err := http.NewRequest(http.MethodGet, config.Config.SlackURL+"/users.lookupByEmail", nil)
	if err != nil {
		log.Println("failed to create a new request:", err)
		return
	}

	params := request.URL.Query()
	params.Set("token", config.Config.Token)
	params.Set("email", user.User.Profile.Email)
	request.URL.RawQuery = params.Encode()

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Println("failed to request for get user info: ", err)
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("read request body failed: ", err)
		return
	}

	if err = json.Unmarshal(body, user); err != nil {
		log.Println("json unmarshal slack user failed: ", err)
		return
	}
}

func (user *SlackUser) verify() bool {
	requestURL := "https://api-reach.herokuapp.com/sessions/update"
	json, err := json.Marshal(user)
	if err != nil {
		log.Println("fail to marshal json: ", err)
		return false
	}

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(json))

	if err != nil {
		log.Println("fail to create request")
		return false
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		log.Println("fail to new request: ", err)
		return false
	}

	return response.StatusCode == 200
}
