package app

import "github.com/ac2393921/docap/pkg/gui"

type App struct {
	Gui *gui.Gui
}

func NewApp() (*App, error) {
	var err error

	app := &App{}
	app.Gui, err = gui.NewGui()

	if err != nil {
		return app, err
	}

	return app, nil
}

func (app *App) Run() error {
	err := app.Gui.Run()
	return err
}
