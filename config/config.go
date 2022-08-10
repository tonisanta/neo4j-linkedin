package config

import (
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"io/ioutil"
)

type Config struct {
	Uri      string `json:"NEO4J_URI"`
	Username string `json:"NEO4J_USERNAME"`
	Password string `json:"NEO4J_PASSWORD"`
}

func ReadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := Config{}
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func NewDriver(settings *Config) (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(
		settings.Uri,
		neo4j.BasicAuth(settings.Username, settings.Password, ""),
	)

	if err != nil {
		return nil, err
	}
	err = driver.VerifyConnectivity()
	if err != nil {
		return nil, err
	}
	return driver, nil
}
