package cmd

import (
	"fmt"

	"charm.land/bubbletea/v2"
	"github.com/smaddiona/ccmgr/internal/config"
	switcher "github.com/smaddiona/ccmgr/internal/tui/switch"
)

func runSwitch() error {
	pf, err := config.LoadProfiles()
	if err != nil {
		return fmt.Errorf("loading profiles: %w", err)
	}

	if len(pf.Profiles) == 0 {
		fmt.Println("No profiles configured. Run 'ccmgr create' to add one.")
		return nil
	}

	m := switcher.NewModel(pf)
	p := tea.NewProgram(m)
	_, err = p.Run()
	return err
}
