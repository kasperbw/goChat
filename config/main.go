package config

import "io/ioutil"
import "encoding/json"
import "time"

type Config struct {
	Server  ServerConfig  `json:"server"`
	Session SessionConfig `json:"session"`
	Auth    AuthConfig    `json:"auth"`
}

//서버 설정 정보
type ServerConfig struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type SessionConfig struct {
	AppKey   string `json:"app_key"`
	Secret   string `json:"secret"`
	UserKey  string `json:"user_key"`
	Duration int    `json:"duration_sec"`
	NextKey  string `json:"next_page_key"`
}

type AuthConfig struct {
	SecurityKey string     `json:"security_key"`
	Google      GoogleAuth `json:"google"`
}

type GoogleAuth struct {
	ClientID string `json:client_id`
	Secret   string `json:"secret"`
}

var instance Config

//MustLoad 파일 읽기
func MustLoad(file string) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(bytes, &instance); err != nil {
		panic(err)
	}

	//session 유지시간을 변경해준다
	instance.Session.Duration *= int(time.Second)
}

func Server() *ServerConfig {
	return &instance.Server
}

func Session() *SessionConfig {
	return &instance.Session
}

func Auth() *AuthConfig {
	return &instance.Auth
}
