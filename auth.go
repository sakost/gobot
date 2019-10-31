package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nikepan/govkbot"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	defaultClientSecret = "hHbZxrka2uZ6jB1inYsH"
	defaultClientId     = 2274003
	authLink            = "https://oauth.vk.com/token"
	configFilename      = "auth_settings.json"
)

type vkResponse struct {
	Token            string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type jsonConfig struct {
	Login string `json:"login"`
	Token string `json:"token"`
}

// get the token from config file(if exists)
// if not, return "", nil
func getCacheToken(login string) (token string, err error) {
	var (
		f         *os.File
		data      []byte
		cacheJson = jsonConfig{}
	)

	f, err = os.Open(configFilename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("config file doesn't exists %s", configFilename)
			return "", nil
		}
		return
	}
	defer f.Close()
	log.Printf("opened %s", configFilename)

	data, err = ioutil.ReadAll(f)
	if err != nil {
		return
	}
	if len(data) == 0 {
		log.Print("no data found in file")
		return "", nil
	}

	log.Print("parsing a data")
	err = json.Unmarshal(data, &cacheJson)
	if err != nil {
		return
	}

	if cacheJson.Login == login {
		token = cacheJson.Token
		log.Printf("found a token for %s", login)
	}
	return
}

// set the token and login in config file
// if force, set the token even if the token already set
func setCacheToken(login, token string, force bool) (n int, err error) {
	var (
		f    *os.File
		data []byte
	)

	f, err = os.Create(configFilename)
	if err != nil && !os.IsExist(err) {
		return
	} else if os.IsExist(err) {
		log.Printf("file %s already exists", configFilename)
		err = nil
	} else {
		log.Printf("created file %s", configFilename)
	}
	defer f.Close()

	f, err = os.OpenFile(configFilename, os.O_RDWR, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	log.Printf("opened %s", configFilename)

	data, err = ioutil.ReadAll(f)
	if err != nil {
		return
	}

	var cacheJson = jsonConfig{}

	if len(data) == 0 {
		log.Print("no data found in config file")
	} else {
		log.Print("parsing a file")
		err = json.Unmarshal(data, &cacheJson)
		if err != nil {
			return
		}
	}

	if cacheJson.Token != "" && cacheJson.Login != "" && cacheJson.Login == login {
		log.Print("login and token are already in file")
		if !force {
			return 0, nil
		}
	}

	// set data into file
	data, err = json.Marshal(jsonConfig{
		Login: login,
		Token: token,
	})
	if err != nil {
		return
	}

	log.Print("writing data to a file")
	n, err = f.Write(data)
	return
}

func GetToken(login, password, clientSecret string, clientId int) (token string, err error) {
	if login == "" || password == "" {
		err = errors.New("no password and/or login")
		return
	}
	log.Printf("requested token for %s login", login)

	token, err = getCacheToken(login)
	if err != nil {
		return
	}
	if token != "" {
		return
	}

	if clientSecret == "" {
		clientSecret = defaultClientSecret
	}
	log.Printf("using client secret \"%s\"", clientSecret)

	if clientId == -1 {
		clientId = defaultClientId
	}
	log.Printf("using client id %d", clientId)

	var q = url.Values{}
	q.Add("scope", "all")
	q.Add("client_id", fmt.Sprint(clientId))
	q.Add("client_secret", clientSecret)
	q.Add("username", login)
	q.Add("password", password)
	q.Add("2fa_supported", "1")
	q.Add("grant_type", "password")
	q.Add("lang", "ru")
	q.Add("v", "5.92") // todo hack
	/* trusted_hash: store.state.users.trustedHashes[login] */

	var (
		client *http.Client
		req    *http.Request
	)
	client = &http.Client{}

	req, err = http.NewRequest("GET", authLink, nil)
	if err != nil {
		return
	}

	req.URL.RawQuery = q.Encode()

	log.Print("making a request...")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var (
		vkResp = vkResponse{}
		data   []byte
	)

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	log.Print("parsing a request")
	err = json.Unmarshal(data, &vkResp)
	if err != nil {
		return
	}
	if govkbot.API.DEBUG {
		log.Printf("response body: \n%s", string(data))
	}

	if vkResp.Error != "" {
		err = errors.New(fmt.Sprintf("%s: %s", vkResp.Error, vkResp.ErrorDescription))
		return
	}

	token = vkResp.Token
	_, _err := setCacheToken(login, token, false)

	if _err != nil {
		log.Printf("error while caching occured: %s\ncontinue without caching token", _err.Error())
	}

	return
}
