package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mouuff/SmartCuts/pkg/types"
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
		return "", fmt.Errorf("could not determine user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".smartcut.json"), nil
}

func GetDefaultConfiguration() *types.SmartCutConfig {
	return &types.SmartCutConfig{
		Model:          "llama3.2",
		MinRowsVisible: 7,
		PromptConfigs: []*types.PromptConfig{
			{
				Index:          0,
				Title:          "Rewrite formally",
				PropertyName:   "result_text",
				PromptTemplate: "Rewrite this formally: '{{input}}'",
			},
			{
				Index:          1,
				Title:          "Rewrite for clarity",
				PropertyName:   "result_text",
				PromptTemplate: "Rewrite this for clarity: '{{input}}'",
			},
		},
	}

}

func GetOrCreateConfiguration() (*types.SmartCutConfig, error) {
	configPath, err := GetConfigurationFilePath()
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
		fmt.Printf("Created default config file at %s\n", configPath)
	}

	var config types.SmartCutConfig
	err = ReadFromJson(configPath, &config)
	if err != nil {
		return nil, fmt.Errorf("could not read config at: %w", err)
	}

	return &config, nil
}
