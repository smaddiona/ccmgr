package cmd

import (
	"charm.land/bubbletea/v2"
	"github.com/smaddiona/ccmgr/internal/tui/create"
)

func runCreate() error {
	m := create.NewModel()
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}
