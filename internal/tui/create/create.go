package create

import (
	"fmt"

	"charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/smaddiona/ccmgr/internal/config"
	"github.com/smaddiona/ccmgr/internal/models"
	"github.com/smaddiona/ccmgr/internal/tui/shared"
)

type step int

const (
	stepPreset step = iota
	stepAPIKey
	stepModels
	stepLabel
	stepDone
)

type Model struct {
	current     step
	width       int
	height      int
	presetList  list.Model
	apiKey      textinput.Model
	opusInput   textinput.Model
	sonnetInput textinput.Model
	haikuInput  textinput.Model
	labelInput  textinput.Model
	modelField  int

	selectedPreset models.Preset
	err            error
}

type errMsg struct{ error }
type savedMsg struct{}

func NewModel() Model {
	items := make([]list.Item, len(models.BuiltInPresets))
	for i, p := range models.BuiltInPresets {
		items[i] = presetItem{
			name:        p.Name,
			needsAPIKey: p.NeedsAPIKey,
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#A6DA95")).Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#9399B2"))

	pl := list.New(items, delegate, 60, 10)
	pl.Title = "Select a Provider Preset"
	pl.Styles.Title = shared.Title

	ak := textinput.New()
	ak.Placeholder = "Enter your API key"
	ak.EchoMode = textinput.EchoPassword
	ak.EchoCharacter = '*'
	ak.Focus()

	opus := textinput.New()
	opus.Placeholder = "Opus model name"
	opus.Focus()

	sonnet := textinput.New()
	sonnet.Placeholder = "Sonnet model name"

	haiku := textinput.New()
	haiku.Placeholder = "Haiku model name"

	label := textinput.New()
	label.Placeholder = "e.g. Z.AI Production"
	label.CharLimit = 50

	return Model{
		current:     stepPreset,
		presetList:  pl,
		apiKey:      ak,
		opusInput:   opus,
		sonnetInput: sonnet,
		haikuInput:  haiku,
		labelInput:  label,
		modelField:  0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.presetList.SetSize(msg.Width, min(msg.Height-4, 10))
		return m, nil

	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEsc {
			if m.current == stepPreset {
				return m, tea.Quit
			}
			m.current--
			return m, nil
		}

	case errMsg:
		m.err = msg.error
		return m, nil

	case savedMsg:
		m.current = stepDone
		return m, tea.Quit
	}

	switch m.current {
	case stepPreset:
		return m.updatePreset(msg)
	case stepAPIKey:
		return m.updateAPIKey(msg)
	case stepModels:
		return m.updateModels(msg)
	case stepLabel:
		return m.updateLabel(msg)
	}

	return m, nil
}

func (m Model) updatePreset(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEnter {
			i, ok := m.presetList.SelectedItem().(presetItem)
			if !ok {
				return m, nil
			}
			preset, _ := models.PresetByName(i.name)
			m.selectedPreset = preset

			if !preset.NeedsAPIKey && !preset.NeedsModels {
				m.current = stepLabel
				m.labelInput.Focus()
				return m, nil
			}

			if preset.NeedsAPIKey {
				m.current = stepAPIKey
				m.apiKey.Focus()
				return m, nil
			}

			m.current = stepModels
			m.applyPresetDefaults()
			m.opusInput.Focus()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.presetList, cmd = m.presetList.Update(msg)
	return m, cmd
}

func (m *Model) applyPresetDefaults() {
	if m.selectedPreset.ModelDefaults.Opus != "" {
		m.opusInput.SetValue(m.selectedPreset.ModelDefaults.Opus)
	}
	if m.selectedPreset.ModelDefaults.Sonnet != "" {
		m.sonnetInput.SetValue(m.selectedPreset.ModelDefaults.Sonnet)
	}
	if m.selectedPreset.ModelDefaults.Haiku != "" {
		m.haikuInput.SetValue(m.selectedPreset.ModelDefaults.Haiku)
	}
}

func (m Model) updateAPIKey(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEnter {
			if m.apiKey.Value() == "" {
				m.err = fmt.Errorf("API key is required")
				return m, nil
			}
			m.err = nil
			m.apiKey.Blur()

			if m.selectedPreset.NeedsModels {
				m.current = stepModels
				m.applyPresetDefaults()
				m.opusInput.Focus()
				return m, nil
			}

			m.current = stepLabel
			m.labelInput.Focus()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.apiKey, cmd = m.apiKey.Update(msg)
	return m, cmd
}

func (m Model) updateModels(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyEnter:
			inputs := []*textinput.Model{&m.opusInput, &m.sonnetInput, &m.haikuInput}
			if m.modelField < 2 {
				inputs[m.modelField].Blur()
				m.modelField++
				inputs[m.modelField].Focus()
				return m, nil
			}
			m.haikuInput.Blur()
			m.current = stepLabel
			m.labelInput.Focus()
			return m, nil
		case tea.KeyTab:
			inputs := []*textinput.Model{&m.opusInput, &m.sonnetInput, &m.haikuInput}
			inputs[m.modelField].Blur()
			m.modelField = (m.modelField + 1) % 3
			inputs[m.modelField].Focus()
			return m, nil
		}
	}

	var cmd tea.Cmd
	switch m.modelField {
	case 0:
		m.opusInput, cmd = m.opusInput.Update(msg)
	case 1:
		m.sonnetInput, cmd = m.sonnetInput.Update(msg)
	case 2:
		m.haikuInput, cmd = m.haikuInput.Update(msg)
	}
	return m, cmd
}

func (m Model) updateLabel(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEnter {
			if m.labelInput.Value() == "" {
				m.err = fmt.Errorf("label is required")
				return m, nil
			}
			m.err = nil
			return m, saveProfile(m.selectedPreset, m.apiKey.Value(), m.opusInput.Value(), m.sonnetInput.Value(), m.haikuInput.Value(), m.labelInput.Value())
		}
	}

	var cmd tea.Cmd
	m.labelInput, cmd = m.labelInput.Update(msg)
	return m, cmd
}

func saveProfile(preset models.Preset, apiKey, opus, sonnet, haiku, label string) tea.Cmd {
	return func() tea.Msg {
		env := make(map[string]string)
		for k, v := range preset.DefaultEnv {
			env[k] = v
		}

		if preset.NeedsAPIKey {
			env["ANTHROPIC_AUTH_TOKEN"] = apiKey
		}
		if preset.NeedsModels {
			if opus != "" {
				env["ANTHROPIC_DEFAULT_OPUS_MODEL"] = opus
			}
			if sonnet != "" {
				env["ANTHROPIC_DEFAULT_SONNET_MODEL"] = sonnet
			}
			if haiku != "" {
				env["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = haiku
			}
		}

		profile := models.NewProfile(label, preset.Name, env)

		pf, err := config.LoadProfiles()
		if err != nil {
			return errMsg{err}
		}

		if _, idx := pf.FindByLabel(label); idx >= 0 {
			return errMsg{fmt.Errorf("a profile with label %q already exists", label)}
		}

		pf.Profiles = append(pf.Profiles, profile)
		if err := config.SaveProfiles(pf); err != nil {
			return errMsg{err}
		}

		return savedMsg{}
	}
}

func (m Model) View() tea.View {
	var s string

	switch m.current {
	case stepPreset:
		s = m.viewPreset()
	case stepAPIKey:
		s = m.viewAPIKey()
	case stepModels:
		s = m.viewModels()
	case stepLabel:
		s = m.viewLabel()
	case stepDone:
		s = shared.Success.Render("Profile saved!")
	}

	if m.err != nil {
		s += "\n" + shared.Error.Render("Error: "+m.err.Error())
	}

	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func (m Model) viewPreset() string {
	return shared.Title.Render("Create New Profile — Step 1: Select Preset") + "\n" +
		m.presetList.View() + "\n" +
		shared.Dim.Render("↑/↓ navigate • enter select • esc cancel")
}

func (m Model) viewAPIKey() string {
	return shared.Title.Render("Create New Profile — Step 2: API Key") + "\n" +
		shared.Label.Render("Provider: "+m.selectedPreset.Name) + "\n\n" +
		shared.Prompt.Render("API Key: ") + m.apiKey.View() + "\n\n" +
		shared.Dim.Render("enter continue • esc back")
}

func (m Model) viewModels() string {
	cursor := func(field int) string {
		if m.modelField == field {
			return shared.ActiveItem.Render(">")
		}
		return " "
	}

	return shared.Title.Render("Create New Profile — Step 3: Model Names") + "\n" +
		shared.Label.Render("Provider: "+m.selectedPreset.Name) + "\n\n" +
		cursor(0) + " Opus:   " + m.opusInput.View() + "\n" +
		cursor(1) + " Sonnet: " + m.sonnetInput.View() + "\n" +
		cursor(2) + " Haiku:  " + m.haikuInput.View() + "\n\n" +
		shared.Dim.Render("tab switch field • enter continue • esc back")
}

func (m Model) viewLabel() string {
	s := shared.Title.Render("Create New Profile — Step 4: Label & Confirm") + "\n\n" +
		shared.Prompt.Render("Profile label: ") + m.labelInput.View() + "\n\n"

	s += shared.Label.Render("Summary:") + "\n"
	s += fmt.Sprintf("  Preset: %s\n", m.selectedPreset.Name)
	if m.selectedPreset.NeedsAPIKey {
		key := m.apiKey.Value()
		if len(key) > 8 {
			key = key[:4] + "..." + key[len(key)-4:]
		}
		s += fmt.Sprintf("  API Key: %s\n", key)
	}
	if m.selectedPreset.NeedsModels {
		s += fmt.Sprintf("  Opus: %s\n", m.opusInput.Value())
		s += fmt.Sprintf("  Sonnet: %s\n", m.sonnetInput.Value())
		s += fmt.Sprintf("  Haiku: %s\n", m.haikuInput.Value())
	}

	s += "\n" + shared.Dim.Render("enter save • esc back")
	return s
}

type presetItem struct {
	name        string
	needsAPIKey bool
}

func (p presetItem) FilterValue() string { return p.name }
func (p presetItem) Title() string       { return p.name }
func (p presetItem) Description() string {
	if p.needsAPIKey {
		return "External provider — requires API key"
	}
	return "Default Anthropic — no env vars"
}
