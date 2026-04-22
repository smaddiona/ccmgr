package models

type Preset struct {
	Name         string
	BaseURL      string
	DefaultEnv   map[string]string
	NeedsAPIKey  bool
	NeedsModels  bool
	ModelDefaults ModelDefaults
}

type ModelDefaults struct {
	Opus   string
	Sonnet string
	Haiku  string
}

var BuiltInPresets = []Preset{
	{
		Name:        "Default",
		NeedsAPIKey: false,
		NeedsModels: false,
	},
	{
		Name:    "Z.AI",
		BaseURL: "https://api.z.ai/api/anthropic",
		DefaultEnv: map[string]string{
			"ANTHROPIC_BASE_URL": "https://api.z.ai/api/anthropic",
			"API_TIMEOUT_MS":     "3000000",
		},
		NeedsAPIKey: true,
		NeedsModels: true,
		ModelDefaults: ModelDefaults{
			Opus:   "glm-5.1",
			Sonnet: "glm-4.7",
			Haiku:  "glm-4.5-air",
		},
	},
}

func PresetByName(name string) (Preset, bool) {
	for _, p := range BuiltInPresets {
		if p.Name == name {
			return p, true
		}
	}
	return Preset{}, false
}
