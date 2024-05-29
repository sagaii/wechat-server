package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type accessTokenStore struct {
	AccessToken       string
	Mutex             sync.RWMutex
	ExpirationSeconds int
}

type response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

var s accessTokenStore

func InitAccessTokenStore() {
	go func() {
		for {
			RefreshAccessToken()
			s.Mutex.RLock()
			sleepDuration := Max(s.ExpirationSeconds, 60)
			s.Mutex.RUnlock()
			time.Sleep(time.Duration(sleepDuration) * time.Second)
		}
	}()
}

func RefreshAccessToken() {
	client := http.Client{
		Timeout: 25 * time.Second, // Increased timeout to handle potential network delays
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", WeChatAppID, WeChatAppSecret), nil)
	if err != nil {
		SysError(err.Error())
		return
	}

	for attempts := 0; attempts < 3; attempts++ {
		responseData, err := client.Do(req)
		if err != nil {
			SysError(fmt.Sprintf("attempt %d: failed to refresh access token for URL %s: %v", attempts+1, req.URL, err))
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		if responseData != nil && responseData.Body != nil {
			defer responseData.Body.Close()
			var res response
			if err = json.NewDecoder(responseData.Body).Decode(&res); err == nil {
				if res.ErrCode != 0 {
					SysError(fmt.Sprintf("attempt %d: failed to refresh access token:  %s: %v", attempts+1, req.URL, err))
					SysError("access token request failed with errcode: " + strconv.Itoa(res.ErrCode) + ", errmsg: " + res.ErrMsg)
					time.Sleep(5 * time.Second) // Wait before retrying
					continue
				}
				s.Mutex.Lock()
				s.AccessToken = res.AccessToken
				s.ExpirationSeconds = res.ExpiresIn
				s.Mutex.Unlock()
				SysLog("========>access token refreshed:" + res.AccessToken)
				break // Success, exit retry loop
			} else {
				SysError(fmt.Sprintf("attempt %d: failed to decode response: %v", attempts+1, err))
			}
		} else {
			SysError(fmt.Sprintf("attempt %d: responseData or responseData.Body is nil", attempts+1))
		}
		time.Sleep(5 * time.Second) // Wait before retrying
	}
}

func GetAccessTokenAndExpirationSeconds() (string, int) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	return s.AccessToken, s.ExpirationSeconds
}

func GetAccessToken() string {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	return s.AccessToken
}
