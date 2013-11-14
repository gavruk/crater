package crater

import (
	"testing"
)

func TestDefaultSettings(t *testing.T) {
	settings := DefaultSettings()

	if settings.ViewsPath != "." {
		t.Error("ViewsPath was not set correctly")
	}
	if settings.StaticFilesPath != "." {
		t.Error("StaticFilesPath was not set correctly")
	}
	if settings.ViewExtension != "html" {
		t.Error("ViewExtension was not set correctly")
	}
}

func TestUpdate(t *testing.T) {
	settings := DefaultSettings()

	newSettings := &Settings{
		ViewsPath:       "./folder",
		StaticFilesPath: "./folder",
		ViewExtension:   "tmpl",
	}

	settings.Update(newSettings)

	if settings.ViewsPath != newSettings.ViewsPath {
		t.Error("ViewsPath was not set correctly")
	}
	if settings.StaticFilesPath != newSettings.StaticFilesPath {
		t.Error("StaticFilesPath was not set correctly")
	}
	if settings.ViewExtension != newSettings.ViewExtension {
		t.Error("ViewExtension was not set correctly")
	}
}
