package crater

// Settings for you app
type Settings struct {
	ViewsPath     string
	StaticPath    string
	ViewExtension string
}

func DefaultSettings() *Settings {
	settings := &Settings{
		ViewsPath:     ".",
		StaticPath:    ".",
		ViewExtension: "html",
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
	if newSettings.StaticPath == "" {
		settings.StaticPath = "."
	} else {
		settings.StaticPath = newSettings.StaticPath
	}
	if newSettings.ViewExtension == "" {
		settings.ViewExtension = "html"
	} else {
		settings.ViewExtension = newSettings.ViewExtension
	}
}
