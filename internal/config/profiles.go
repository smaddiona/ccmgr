package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/smaddiona/ccmgr/internal/models"
)

func ProfilesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, ".ccmgr"), nil
}

func ProfilesPath() (string, error) {
	dir, err := ProfilesDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "profiles.json"), nil
}

func ensureProfilesDir() error {
	dir, err := ProfilesDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0700)
}

func LoadProfiles() (models.ProfilesFile, error) {
	path, err := ProfilesPath()
	if err != nil {
		return models.ProfilesFile{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return models.NewProfilesFile(), nil
		}
		return models.ProfilesFile{}, fmt.Errorf("reading profiles: %w", err)
	}

	var pf models.ProfilesFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return models.ProfilesFile{}, fmt.Errorf("parsing profiles: %w", err)
	}
	return pf, nil
}

func SaveProfiles(pf models.ProfilesFile) error {
	if err := ensureProfilesDir(); err != nil {
		return err
	}

	path, err := ProfilesPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling profiles: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing profiles: %w", err)
	}
	return nil
}
