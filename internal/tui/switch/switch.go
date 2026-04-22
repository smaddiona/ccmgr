package switcher

import (
	"fmt"
	"time"

	"charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
	"github.com/smaddiona/ccmgr/internal/config"
	"github.com/smaddiona/ccmgr/internal/models"
	"github.com/smaddiona/ccmgr/internal/tui/shared"
)

type Model struct {
	list       list.Model
	profiles   models.ProfilesFile
	selectedID string
	applied    bool
	err        error
	quitting   bool
}

type appliedMsg struct {
	id string
}

type errMsg struct{ error }

func NewModel(pf models.ProfilesFile) Model {
	items := make([]list.Item, len(pf.Profiles))
	for i, p := range pf.Profiles {
		items[i] = profileItem{id: p.ID, label: p.Label, preset: p.Preset, active: p.ID == pf.ActiveID}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#A6DA95")).Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#9399B2"))

	l := list.New(items, delegate, 60, 10)
	l.Title = "Select Configuration Profile"
	l.Styles.Title = shared.Title

	return Model{
		list:     l,
		profiles: pf,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, min(msg.Height-4, 10))
		return m, nil

	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEsc {
			m.quitting = true
			return m, tea.Quit
		}
		if msg.Code == tea.KeyEnter && !m.applied {
			i, ok := m.list.SelectedItem().(profileItem)
			if !ok {
				return m, nil
			}
			m.selectedID = i.id
			return m, applyProfile(i.id, m.profiles)
		}

	case appliedMsg:
		m.applied = true
		m.profiles.ActiveID = msg.id
		return m, tea.Tick(1500*time.Millisecond, func(t time.Time) tea.Msg {
			return tea.Quit()
		})

	case errMsg:
		m.err = msg.error
		return m, nil
	}

	if m.applied {
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func applyProfile(id string, pf models.ProfilesFile) tea.Cmd {
	return func() tea.Msg {
		settingsPath, err := config.ClaudeSettingsPath()
		if err != nil {
			return errMsg{err}
		}

		profile, _ := pf.FindByID(id)
		if profile == nil {
			return errMsg{fmt.Errorf("profile not found")}
		}

		env := profile.Env
		if len(env) == 0 {
			env = nil
		}

		if err := config.ApplyEnvToSettings(settingsPath, env); err != nil {
			return errMsg{err}
		}

		pf.ActiveID = id
		if err := config.SaveProfiles(pf); err != nil {
			return errMsg{err}
		}

		return appliedMsg{id: id}
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

	if m.applied {
		profile, _ := m.profiles.FindByID(m.profiles.ActiveID)
		msg := fmt.Sprintf("Switched to %q", profile.Label)
		v.SetContent(shared.Success.Render(msg))
		return v
	}

	v.SetContent(m.list.View() + "\n" + shared.Dim.Render("↑/↓ navigate • enter select • esc cancel"))
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
