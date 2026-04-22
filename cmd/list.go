package cmd

import (
	"fmt"
	"os"

	"github.com/smaddiona/ccmgr/internal/config"
)

func runList() error {
	pf, err := config.LoadProfiles()
	if err != nil {
		return fmt.Errorf("loading profiles: %w", err)
	}

	if len(pf.Profiles) == 0 {
		fmt.Println("No profiles configured. Run 'ccmgr create' to add one.")
		return nil
	}

	for _, p := range pf.Profiles {
		marker := "  "
		if p.ID == pf.ActiveID {
			marker = "* "
		}

		activeLabel := ""
		if p.ID == pf.ActiveID {
			activeLabel = "  active"
		}

		fmt.Printf("%s%-20s [%s]%s\n", marker, p.Label, p.Preset, activeLabel)
	}

	return nil
}

func checkProfiles() {
	if _, err := config.LoadProfiles(); err != nil {
		fatal("Error loading profiles", err)
	}
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fatal("Error getting home directory", err)
	}
	return home
}
