package models

import (
	"fmt"
	"time"
)

type Profile struct {
	ID        string            `json:"id"`
	Label     string            `json:"label"`
	Preset    string            `json:"preset"`
	Env       map[string]string `json:"env,omitempty"`
	CreatedAt string            `json:"createdAt"`
	UpdatedAt string            `json:"updatedAt"`
}

type ProfilesFile struct {
	Version  string    `json:"version"`
	ActiveID string    `json:"activeId"`
	Profiles []Profile `json:"profiles"`
}

func NewProfilesFile() ProfilesFile {
	return ProfilesFile{
		Version:  "1",
		ActiveID: "",
		Profiles: []Profile{},
	}
}

func NewProfile(label, preset string, env map[string]string) Profile {
	now := time.Now().UTC().Format(time.RFC3339)
	return Profile{
		ID:        generateID(),
		Label:     label,
		Preset:    preset,
		Env:       env,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (pf *ProfilesFile) FindByID(id string) (*Profile, int) {
	for i := range pf.Profiles {
		if pf.Profiles[i].ID == id {
			return &pf.Profiles[i], i
		}
	}
	return nil, -1
}

func (pf *ProfilesFile) FindByLabel(label string) (*Profile, int) {
	for i := range pf.Profiles {
		if pf.Profiles[i].Label == label {
			return &pf.Profiles[i], i
		}
	}
	return nil, -1
}

func (pf *ProfilesFile) Active() *Profile {
	if pf.ActiveID == "" {
		return nil
	}
	p, _ := pf.FindByID(pf.ActiveID)
	return p
}

func (pf *ProfilesFile) Remove(id string) {
	for i, p := range pf.Profiles {
		if p.ID == id {
			pf.Profiles = append(pf.Profiles[:i], pf.Profiles[i+1:]...)
			if pf.ActiveID == id {
				pf.ActiveID = ""
			}
			return
		}
	}
}

func generateID() string {
	// Simple timestamp-based ID, sufficient for this use case
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
