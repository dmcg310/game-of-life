package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const defaultFilename = "gol-config.json"

type Config struct {
	Preset          string `json:"preset"`
	CellColor       string `json:"cell-color"`
	BackgroundColor string `json:"background-color"`
	ScaleFactor     int    `json:"scale-factor"`
	FPS             int    `json:"fps"`
}

func NewConfig(bytes []byte) *Config {
	if bytes == nil {
		return nil
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		msg := fmt.Sprintf("Cannot parse JSON in config file: '%s'. Continued with defaults.",
			defaultFilename)
		NewAppWarning(msg, "Please ensure that the JSON contains no syntactical errors.").
			ShowAppWarning()

		return nil
	}

	return &config
}

func NewConfigWithDefaults() *Config {
	return &Config{
		Preset:          "random",
		CellColor:       "white",
		BackgroundColor: "black",
		ScaleFactor:     1,
		FPS:             23,
	}
}

func ReadConfig() []byte {
	file, err := os.Open(defaultFilename)
	if err != nil {
		// if the file isnt found that is ok, just continue with defaults
		msg := fmt.Sprintf("Cannot open config file: '%s'. Continued with defaults.",
			defaultFilename)
		NewAppWarning(msg, "Make sure the file exists and is accessible by the program.").
			ShowAppWarning()

		return nil
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		msg := fmt.Sprintf("Cannot read config file: '%s'. Continued with defaults.",
			defaultFilename)
		NewAppWarning(msg, "Please try re-running the program.").
			ShowAppWarning()

		return nil
	}

	return bytes
}
