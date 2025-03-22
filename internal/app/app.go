package app

import (
	"runtime"
	"runtime/debug"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
)

const (
	// App information
	AppName = "Algorand Rewards"
	AppID   = "com.calmdev.algorand-rewards"

	// Preference keys
	AddressKey     = "Address"
	GUIDKey        = "GUID"
	RewardsViewKey = "RewardsView"
	VersionKey     = "Version"
)

// CurrentApp returns the current instance of the App.
func CurrentApp() *App {
	return &App{
		App: fyne.CurrentApp(),
	}
}

// App represents the main application structure.
type App struct {
	fyne.App
}

// NewApp creates a new instance of the App.
func NewApp() *App {
	return &App{
		App: app.NewWithID(AppID),
	}
}

// VersionCheck checks if the version has changed and executes the handler if it has.
func (a *App) VersionCheck(handler func()) {
	if version := a.buildVersion(); version != a.Version() {
		a.setVersion(version)
		handler()
	}
}

// buildVersion returns the current build version of the app.
func (a *App) buildVersion() string {
	if b, ok := debug.ReadBuildInfo(); ok && len(b.Main.Version) > 0 {
		return b.Main.Version
	}
	return "unknown"
}

// Version returns the current version of the app.
func (a *App) Version() string {
	return a.Preferences().String(VersionKey)
}

// setVersion sets the version of the app.
func (a *App) setVersion(version string) {
	a.Preferences().SetString(VersionKey, version)
}

// Address returns the address associated with the app.
func (a *App) Address() string {
	return a.Preferences().String(AddressKey)
}

// SetAddress sets the address associated with the app.
func (a *App) SetAddress(value string) {
	a.Preferences().SetString(AddressKey, value)
}

// GUID returns the GUID associated with the app.
func (a *App) GUID() string {
	return a.Preferences().String(GUIDKey)
}

// SetGUID sets the GUID associated with the app.
func (a *App) SetGUID(value string) {
	a.Preferences().SetString(GUIDKey, value)
}

// RewardsView returns the RevardsView associated with the app.
func (a *App) RewardsView() string {
	return a.Preferences().String(RewardsViewKey)
}

// SetRewardsView sets the RewardsView associated with the app.
func (a *App) SetRewardsView(value string) {
	a.Preferences().SetString(RewardsViewKey, value)
}

// CacheFile returns the cache file for the given file name.
func (a *App) CacheFile(fileName string) (fyne.URI, error) {
	cacheFile, err := storage.Child(a.Storage().RootURI(), fileName)
	if err != nil {
		return nil, err
	}
	return cacheFile, nil
}

// ClearCacheFile clears the cache file for the given file names.
func (a *App) ClearCacheFile(fileNames ...string) error {
	for _, fileName := range fileNames {
		cacheFile, err := a.CacheFile(fileName)
		if err != nil {
			return err
		}
		err = storage.Delete(cacheFile)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsWindows returns true if the app is running on Windows.
func (a *App) IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsMacOS returns true if the app is running on macOS.
func (a *App) IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux returns true if the app is running on Linux.
func (a *App) IsLinux() bool {
	return runtime.GOOS == "linux"
}
