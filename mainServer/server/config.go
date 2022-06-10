package server

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Hosting  hostingConfig  `json:"hosting"`
	Database databaseConfig `json:"database"`
	Auth     authConfig     `json:"auth"`
	Git      gitConfig      `json:"git"`
}

type databaseConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
	Name string `json:"name"`
}

type gitConfig struct {
	Path string `json:"path"`
}

type authConfig struct {
	JwtSecret          string `json:"jwt-secret"`
	TokenExpireMinutes int    `json:"token-expire-minutes"`
}

type hostingConfig struct {
	BackendURL  string `json:"back-end-url"`
	FrontendURL string `json:"front-end-url"`
}

// ReadConfig reads a JSON configuration file for the system into a Config struct
// The JSON file can include fields that are not specified in the struct, if those fields are only used in the frontend
func ReadConfig(path string) Config {
	// open file
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	cfgFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	// parse json to struct
	decoder := json.NewDecoder(cfgFile)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	// close file
	err = cfgFile.Close()
	if err != nil {
		panic(err)
	}

	return config
}
