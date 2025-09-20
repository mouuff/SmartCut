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
				Model:        "llama3.2",
				Title:        "Rewrite for clarity (llama3.2)",
				SystemPrompt: "You are a rewriting assistant. Your task is to take the user's messages and rewrite them so they are clear, concise, and natural-sounding, while keeping the original meaning. Always use a human, conversational tone—avoid robotic or overly formal phrasing. Do not add extra information or change the intent",
				Prompt:       "please rewrite this: '{{input}}'",
				PropertyName: "rewritten_text",
			},
			{
				Model:        "mistral",
				Title:        "Rewrite for clarity (mistral)",
				SystemPrompt: "Your role is to act as a text rewriter. For every message the user provides, rewrite it so it sounds clear, natural, and human-like. Keep the meaning intact, but improve readability and flow. Avoid jargon, stiffness, or AI-like phrasing. Be concise, friendly, and easy to understand.",
				Prompt:       "please rewrite this: '{{input}}'",
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
