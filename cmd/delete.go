package cmd

import (
	"fmt"

	"charm.land/bubbletea/v2"
	"github.com/smaddiona/ccmgr/internal/config"
	"github.com/smaddiona/ccmgr/internal/tui/delete"
)

func runDelete(label string) error {
	pf, err := config.LoadProfiles()
	if err != nil {
		return fmt.Errorf("loading profiles: %w", err)
	}

	if len(pf.Profiles) == 0 {
		fmt.Println("No profiles configured.")
		return nil
	}

	var m tea.Model

	if label != "" {
		profile, idx := pf.FindByLabel(label)
		if idx < 0 {
			return fmt.Errorf("profile %q not found", label)
		}
		m = delete.NewConfirmModel(pf, profile)
	} else {
		m = delete.NewModel(pf)
	}

	p := tea.NewProgram(m)
	_, err = p.Run()
	return err
}
