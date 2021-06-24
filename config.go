package main

import (
	"encoding/json"
	"io/ioutil"
)

type DbConfig struct {
	DbType   string `json:"dbType"`
	DbSource string `json:"dbSource"`
}

type Config struct {
	ServerUrl string   `json:"serverUrl"`
	DbConfig  DbConfig `json:"dbConfig"`
}

var GlobalConfig Config

type JsonStruct struct {
}

func (j *JsonStruct) Load(configFile string, config interface{}) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}
	return nil
}
