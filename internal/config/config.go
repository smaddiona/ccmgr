package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ClaudeSettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, ".claude", "settings.json"), nil
}

func ReadClaudeSettings(path string) (map[string]json.RawMessage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]json.RawMessage), nil
		}
		return nil, fmt.Errorf("reading settings: %w", err)
	}
	result := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parsing settings: %w", err)
	}
	return result, nil
}

func WriteClaudeSettings(path string, data map[string]json.RawMessage) error {
	// Backup existing file before overwriting
	if _, err := os.Stat(path); err == nil {
		backup := path + ".bak"
		existing, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading existing settings for backup: %w", err)
		}
		if err := os.WriteFile(backup, existing, 0600); err != nil {
			return fmt.Errorf("creating backup: %w", err)
		}
	}

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling settings: %w", err)
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0644)
}

func ApplyEnvToSettings(settingsPath string, env map[string]string) error {
	settings, err := ReadClaudeSettings(settingsPath)
	if err != nil {
		return err
	}

	if len(env) > 0 {
		envBytes, err := json.Marshal(env)
		if err != nil {
			return fmt.Errorf("marshaling env: %w", err)
		}
		settings["env"] = envBytes
	} else {
		delete(settings, "env")
	}

	return WriteClaudeSettings(settingsPath, settings)
}
