package crater

// Settings for you app
type Settings struct {
	ViewsPath       string
	StaticFilesPath string
	ViewExtension   string
}

func DefaultSettings() *Settings {
	settings := &Settings{
		ViewsPath:       ".",
		StaticFilesPath: ".",
		ViewExtension:   "html",
	}
	return settings
}

func (settings *Settings) Update(newSettings *Settings) {
	if settings == nil {
		return
	}
	if newSettings.ViewsPath == "" {
		settings.ViewsPath = "."
	} else {
		settings.ViewsPath = newSettings.ViewsPath
	}
	if newSettings.StaticFilesPath == "" {
		settings.StaticFilesPath = "."
	} else {
		settings.StaticFilesPath = newSettings.StaticFilesPath
	}
	if newSettings.ViewExtension == "" {
		settings.ViewExtension = "html"
	} else {
		settings.ViewExtension = newSettings.ViewExtension
	}
}
