package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Hosting  HostingConfig  `json:"hosting"`
	Database DatabaseConfig `json:"database"`
	Auth     AuthConfig     `json:"auth"`
	Fs       StorerConfig   `json:"fs"`
}

type DatabaseConfig struct {
	Url    URLConfig `json:"url"`
	User   string    `json:"user"`
	Pwd    string    `json:"pwd"`
	DBName string    `json:"dbname"`
}

type StorerConfig struct {
	Path       string `json:"path"`
	MutexCount int    `json:"mutex-count"`
}

type AuthConfig struct {
	JwtSecret          string `json:"jwt-secret"`
	TokenExpireMinutes int    `json:"token-expire-minutes"`
}

type HostingConfig struct {
	Backend  URLConfig `json:"back-end"`
	Frontend URLConfig `json:"front-end"`
}

type URLConfig struct {
	Port   int    `json:"port"`
	Host   string `json:"host"`
	UseSSL bool   `json:"use-ssl"` // determines whether prefix should be http or https
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

// Url creates an actual url string from the struct fields
func (url URLConfig) Url() string {
	var prefix string
	if url.UseSSL {
		prefix = "https"
	} else {
		prefix = "http"
	}
	return fmt.Sprintf("%s://%s:%d", prefix, url.Host, url.Port)
}
