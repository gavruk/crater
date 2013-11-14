package crater

// Settings is a
type Settings struct {
	// Path to the root folder with html files
	ViewsPath string
	// Path to the root folder with static content (js, css, images)
	StaticPath string
	// Extension of the View files or templates
	ViewExtension string
}

// DefaultSettings creates a settings object with default settings
func DefaultSettings() *Settings {
	settings := &Settings{
		ViewsPath:     ".",
		StaticPath:    ".",
		ViewExtension: "html",
	}
	return settings
}

// Update updates current settings with passed settings data
// If settings filed is empty, default value with be used for this field
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
