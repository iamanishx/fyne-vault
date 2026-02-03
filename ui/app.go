package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"

	"passvault-fyne/internal/clipboard"
	"passvault-fyne/internal/database"
	"passvault-fyne/internal/state"
)

type App struct {
	FyneApp    fyne.App
	MainWindow fyne.Window
	DB         *database.DB
	State      *state.SecureState
	Clipboard  *clipboard.Manager
	Secrets    []database.SecretEntry
	Filtered   []database.SecretEntry
	List       *widget.List
	Search     *widget.Entry
}

func NewApp() (*App, error) {
	a := app.NewWithID("com.passvault.app")
	w := a.NewWindow("PassVault")

	db, err := database.NewDB()
	if err != nil {
		return nil, err
	}

	if err := db.InitSchema(); err != nil {
		return nil, err
	}

	s := state.NewSecureState()
	c := clipboard.NewManager(w.Clipboard())

	passVaultApp := &App{
		FyneApp:    a,
		MainWindow: w,
		DB:         db,
		State:      s,
		Clipboard:  c,
	}

	a.Settings().SetTheme(&TransparentTheme{})
	w.SetPadded(false)

	return passVaultApp, nil
}

func (a *App) Start() {
	a.MainWindow.Resize(fyne.NewSize(1024, 768))
	a.MainWindow.CenterOnScreen()
	a.MainWindow.Show()
	a.MainWindow.RequestFocus()
	a.ShowMainUI()
	a.FyneApp.Run()
}

func (a *App) applyWindowTweaks() {}
