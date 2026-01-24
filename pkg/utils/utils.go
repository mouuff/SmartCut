package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mouuff/SmartCut/pkg/types"
)

func ReadFromJson(path string, dataOut interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(data), dataOut); err != nil {
		return err
	}

	return nil
}

func GetConfigurationFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".smartcut.v2.json"), nil
}

func GetDefaultConfiguration() *types.SmartCutConfig {
	return &types.SmartCutConfig{
		HostUrl:        "default",
		MinRowsVisible: 7,
		PromptConfigs: []types.PromptConfig{
			{
				Model:        "mistral",
				Title:        "Rewrite for clarity (mistral)",
				SystemPrompt: "You are a rewriting assistant. Your task is to take the user's messages and rewrite them so they are clear, concise, and natural-sounding, while keeping the original meaning. Always use a human, conversational tone—avoid robotic or overly formal phrasing. Do not add extra information or change the intent",
				Prompt:       "please rewrite this: '{{input}}'",
				PropertyName: "rewritten_text",
			},
			{
				Model:        "mistral",
				Title:        "Fix grammar (mistral)",
				SystemPrompt: "You are a grammar editor.\nCorrect spelling, grammar, and punctuation in the text below.\nDo NOT change the meaning, tone, or structure.\nDo NOT rewrite or rephrase\nsentences unless required to fix grammar.\nReturn only the corrected text.",
				Prompt:       "{{input}}",
				PropertyName: "rewritten_text",
			},
		},
	}

}

func GetOrCreateConfiguration(configPath string) (*types.SmartCutConfig, error) {
	var err error

	if configPath == "" {
		configPath, err = GetConfigurationFilePath()
	}

	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {

		jsonData, err := json.MarshalIndent(GetDefaultConfiguration(), "", "  ")
		if err != nil {
			return nil, fmt.Errorf("could not marshal default config: %w", err)
		}

		err = os.WriteFile(configPath, jsonData, 0644)
		if err != nil {
			return nil, fmt.Errorf("could not create default config file: %w", err)
		}

		log.Printf("Created default config file at %s\n", configPath)
	}

	var config types.SmartCutConfig
	err = ReadFromJson(configPath, &config)
	if err != nil {
		return nil, fmt.Errorf("could not read config at %s: %w", configPath, err)
	}

	// Sets the config path
	config.ConfigPath = configPath

	return &config, nil
}

func OpenFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// macOS
		cmd = exec.Command("open", path)
	case "linux":
		// Linux
		cmd = exec.Command("xdg-open", path)
	case "windows":
		// Windows
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	// Make sure it inherits environment variables
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Start()
}
