package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.

//go:embed assets/icon/shipping-container.png
var iconStopped []byte

//go:embed assets/icon/shipping-container-fill.png
var iconRunning []byte

func main() {
	app := application.New(application.Options{
		Name: "isColimaRunning",
		Mac:  application.MacOptions{ActivationPolicy: application.ActivationPolicyAccessory},
	})

	tray := app.SystemTray.New()
	tray.SetTemplateIcon(iconStopped)

	menu := app.NewMenu()
	menu.Add("isColimaRunning").SetEnabled(false).SetBitmap(iconRunning)
	menu.AddSeparator()
	status := menu.Add("Status: checking…").SetEnabled(false)
	menu.AddSeparator()
	start := menu.Add("Start")
	stop := menu.Add("Stop")
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(*application.Context) { app.Quit() })
	tray.SetMenu(menu)

	render := func(label string, running, busy bool) {
		status.SetLabel("Status: " + label)
		start.SetEnabled(!running && !busy)
		stop.SetEnabled(running && !busy)
		if running {
			tray.SetTemplateIcon(iconRunning)
		} else {
			tray.SetTemplateIcon(iconStopped)
		}
		menu.Update()
	}
	refresh := func() {
		running := isColimaRunning()
		if running {
			render("running", true, false)
		} else {
			render("stopped", false, false)
		}
	}

	runAction := func(label string, action func() error, wasRunning bool) {
		render(label, wasRunning, true)
		go func() {
			if err := action(); err != nil {
				log.Printf("%s failed: %v", label, err)
			}
			refresh()
		}()
	}

	start.OnClick(func(*application.Context) { runAction("starting…", startColima, false) })
	stop.OnClick(func(*application.Context) { runAction("stopping…", stopColima, true) })

	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		refresh()
		go func() {
			for range time.Tick(5 * time.Second) {
				refresh()
			}
		}()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
