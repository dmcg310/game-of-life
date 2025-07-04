package config

import (
	"encoding/json"
	"fmt"
	"github.com/dmcg310/game-of-life/internal/errors"
	"io"
	"os"
)

const defaultFilename = "gol-config.json"

type Config struct {
	Preset          string `json:"preset"`
	CellColor       string `json:"cell-color"`
	BackgroundColor string `json:"background-color"`
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
		errors.NewAppWarning(msg, "Please ensure that the JSON contains no syntactical errors.").
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
		FPS:             23,
	}
}

func ReadConfig() []byte {
	var file *os.File

	file, err := os.Open(defaultFilename)
	if err != nil {
		// if not in current directory check the configuration directory
		configDir, err := os.UserConfigDir()
		if err != nil {
			errors.NewAppWarning("Cannot get the current configuration directory. Continued with defaults.",
				"Please try re-running the program.").ShowAppWarning()

			return nil
		}

		path := fmt.Sprintf("%s/gol/%s", configDir, defaultFilename)
		file, err = os.Open(path)
		if err != nil {
			msg := fmt.Sprintf("Cannot open config file: '%s'. Continued with defaults.",
				defaultFilename)
			errors.NewAppWarning(msg, "Make sure the file exists and is accessible by the program.").
				ShowAppWarning()

			return nil
		}
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		msg := fmt.Sprintf("Cannot read config file: '%s'. Continued with defaults.",
			defaultFilename)
		errors.NewAppWarning(msg, "Please try re-running the program.").
			ShowAppWarning()

		return nil
	}

	return bytes
}
