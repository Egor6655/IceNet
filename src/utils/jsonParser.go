package utils

import (
	"embed"
	"encoding/json"
)

//go:embed config/config.json
var fs embed.FS

type Config struct {
	Links []string //`json:"Links"`
}

type Response struct {
	Target     string
	Times      int
	Typemethod string
	Cmd        string `json:"Cmd"`
	Mirrors    []string
}

func ParseTarget(body string) Response {
	var data Response
	json.Unmarshal([]byte(body), &data)

	return data

}

func ParseUrls() Config {
	var config Config

	data, _ := fs.ReadFile(`config/config.json`)
	json.Unmarshal([]byte(data), &config)

	return config

}
