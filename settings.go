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
