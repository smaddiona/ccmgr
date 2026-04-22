package delete

import (
	"fmt"

	"charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
	"github.com/smaddiona/ccmgr/internal/config"
	"github.com/smaddiona/ccmgr/internal/models"
	"github.com/smaddiona/ccmgr/internal/tui/shared"
)

type step int

const (
	stepSelect step = iota
	stepConfirm
	stepDone
)

type Model struct {
	current  step
	list     list.Model
	profiles models.ProfilesFile
	selected *models.Profile
	err      error
	quitting bool
}

type errMsg struct{ error }
type doneMsg struct{}

func NewModel(pf models.ProfilesFile) Model {
	items := make([]list.Item, len(pf.Profiles))
	for i, p := range pf.Profiles {
		items[i] = profileItem{id: p.ID, label: p.Label, preset: p.Preset, active: p.ID == pf.ActiveID}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#ED8796")).Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#9399B2"))

	l := list.New(items, delegate, 60, 10)
	l.Title = "Delete Profile"
	l.Styles.Title = shared.Title

	return Model{
		current:  stepSelect,
		list:     l,
		profiles: pf,
	}
}

func NewConfirmModel(pf models.ProfilesFile, profile *models.Profile) Model {
	return Model{
		current:  stepConfirm,
		profiles: pf,
		selected: profile,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.list.Title != "" {
			m.list.SetSize(msg.Width, min(msg.Height-4, 10))
		}
		return m, nil

	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEsc {
			if m.current == stepConfirm && len(m.profiles.Profiles) > 1 {
				m.current = stepSelect
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg.error
		return m, nil

	case doneMsg:
		m.current = stepDone
		return m, nil
	}

	switch m.current {
	case stepSelect:
		return m.updateSelect(msg)
	case stepConfirm:
		return m.updateConfirm(msg)
	}

	return m, nil
}

func (m Model) updateSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEnter {
			i, ok := m.list.SelectedItem().(profileItem)
			if !ok {
				return m, nil
			}
			profile, _ := m.profiles.FindByID(i.id)
			if profile != nil {
				m.selected = profile
				m.current = stepConfirm
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Text {
		case "y", "Y":
			if m.selected == nil {
				return m, nil
			}
			return m, deleteProfile(m.selected.ID, m.profiles)
		case "n", "N":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func deleteProfile(id string, pf models.ProfilesFile) tea.Cmd {
	return func() tea.Msg {
		wasActive := pf.ActiveID == id
		pf.Remove(id)

		if err := config.SaveProfiles(pf); err != nil {
			return errMsg{err}
		}

		if wasActive {
			settingsPath, err := config.ClaudeSettingsPath()
			if err == nil {
				_ = config.ApplyEnvToSettings(settingsPath, nil)
			}
		}

		return doneMsg{}
	}
}

func (m Model) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true

	if m.quitting {
		return v
	}

	if m.err != nil {
		v.SetContent(shared.Error.Render("Error: " + m.err.Error()))
		return v
	}

	switch m.current {
	case stepSelect:
		v.SetContent(m.list.View() + "\n" + shared.Dim.Render("↑/↓ navigate • enter select • esc cancel"))
	case stepConfirm:
		v.SetContent(
			shared.Title.Render("Confirm Deletion") + "\n\n" +
				fmt.Sprintf("Delete profile %q [%s]?", m.selected.Label, m.selected.Preset) + "\n\n" +
				shared.Dim.Render("y confirm • n cancel"),
		)
	case stepDone:
		v.SetContent(shared.Success.Render("Profile deleted."))
	}

	return v
}

type profileItem struct {
	id     string
	label  string
	preset string
	active bool
}

func (p profileItem) FilterValue() string { return p.label }
func (p profileItem) Title() string {
	if p.active {
		return "* " + p.label
	}
	return "  " + p.label
}
func (p profileItem) Description() string {
	if p.active {
		return fmt.Sprintf("[%s] — active", p.preset)
	}
	return fmt.Sprintf("[%s]", p.preset)
}
