package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/smaddiona/ccmgr/internal/models"
)

func TestReadSettingsPreservesFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	original := map[string]interface{}{
		"model":    "opus[1m]",
		"env":      map[string]string{"OLD_KEY": "old_val"},
		"plugins":  []string{"a", "b"},
		"newField": true,
	}
	data, _ := json.MarshalIndent(original, "", "  ")
	os.WriteFile(path, data, 0644)

	settings, err := ReadClaudeSettings(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	newEnv := map[string]string{"ANTHROPIC_AUTH_TOKEN": "new-key"}
	envBytes, _ := json.Marshal(newEnv)
	settings["env"] = envBytes

	if err := WriteClaudeSettings(path, settings); err != nil {
		t.Fatalf("write: %v", err)
	}

	result, _ := os.ReadFile(path)
	var parsed map[string]interface{}
	json.Unmarshal(result, &parsed)

	if parsed["model"] != "opus[1m]" {
		t.Errorf("model not preserved: %v", parsed["model"])
	}
	if parsed["newField"] != true {
		t.Errorf("newField not preserved: %v", parsed["newField"])
	}
	envMap, ok := parsed["env"].(map[string]interface{})
	if !ok {
		t.Fatal("env is not a map")
	}
	if envMap["ANTHROPIC_AUTH_TOKEN"] != "new-key" {
		t.Errorf("env not updated: %v", envMap)
	}
}

func TestApplyEnvRemovesKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	original := map[string]interface{}{
		"model": "opus",
		"env":   map[string]string{"KEY": "val"},
	}
	data, _ := json.MarshalIndent(original, "", "  ")
	os.WriteFile(path, data, 0644)

	if err := ApplyEnvToSettings(path, nil); err != nil {
		t.Fatalf("apply: %v", err)
	}

	result, _ := os.ReadFile(path)
	var parsed map[string]interface{}
	json.Unmarshal(result, &parsed)

	if _, exists := parsed["env"]; exists {
		t.Error("env key should have been removed")
	}
	if parsed["model"] != "opus" {
		t.Error("model not preserved")
	}
}

func TestApplyEnvSetsValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	os.WriteFile(path, []byte(`{"model": "opus"}`), 0644)

	env := map[string]string{
		"ANTHROPIC_AUTH_TOKEN": "test-key",
		"ANTHROPIC_BASE_URL":   "https://api.z.ai/api/anthropic",
	}
	if err := ApplyEnvToSettings(path, env); err != nil {
		t.Fatalf("apply: %v", err)
	}

	settings, _ := ReadClaudeSettings(path)
	var envMap map[string]string
	json.Unmarshal(settings["env"], &envMap)
	if envMap["ANTHROPIC_AUTH_TOKEN"] != "test-key" {
		t.Errorf("token not set: %v", envMap)
	}
}

func TestProfileCRUD(t *testing.T) {
	pf := models.NewProfilesFile()
	p1 := models.NewProfile("Z.AI Prod", "Z.AI", map[string]string{"KEY": "val"})
	p2 := models.NewProfile("Default", "Default", nil)

	pf.Profiles = append(pf.Profiles, p1, p2)

	if len(pf.Profiles) != 2 {
		t.Fatalf("expected 2 profiles")
	}

	found, idx := pf.FindByLabel("Z.AI Prod")
	if idx != 0 || found.ID != p1.ID {
		t.Error("FindByLabel failed")
	}

	pf.ActiveID = p1.ID
	active := pf.Active()
	if active.Label != "Z.AI Prod" {
		t.Error("Active() failed")
	}

	pf.Remove(p1.ID)
	if len(pf.Profiles) != 1 {
		t.Error("Remove failed")
	}
	if pf.ActiveID != "" {
		t.Error("ActiveID should be cleared on remove")
	}

	_, idx = pf.FindByLabel("Z.AI Prod")
	if idx >= 0 {
		t.Error("should not find removed profile")
	}
}

func TestBackupCreated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	os.WriteFile(path, []byte(`{"model":"old"}`), 0644)

	settings, _ := ReadClaudeSettings(path)
	WriteClaudeSettings(path, settings)

	if _, err := os.Stat(path + ".bak"); err != nil {
		t.Error("backup file not created")
	}
}
