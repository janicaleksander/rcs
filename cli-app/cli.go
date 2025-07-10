package cli_app

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"github.com/rivo/tview"
	"reflect"
	"time"
)

type Application struct {
	app   *tview.Application
	pages *tview.Pages
}

type CLI struct {
	serverPID   *actor.PID
	application Application
}

func NewCLI() actor.Producer {
	return func() actor.Receiver {
		return &CLI{}
	}
}

func (c *CLI) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
	case actor.Initialized:
	case actor.Stopped:
	case *Proto.NeededServerConfiguration:
		c.serverPID = actor.NewPID(msg.ServerPID.Address, msg.ServerPID.Id)
	case *Proto.StartCLI:
		// creating tview application
		c.application.app = tview.NewApplication()
		c.application.pages = tview.NewPages()

		//prep forms etc
		loginForm := c.createLoginForm(ctx)
		c.application.pages.AddPage("login", loginForm, true, true)

		//handle pages and run
		if err := c.application.app.SetRoot(c.application.pages, true).Run(); err != nil {
			Server.Logger.Error("CLI can't start")
		}
	default:
		Server.Logger.Warn("Server got unknown message", "Type:", reflect.TypeOf(msg).String())
	}
}

func (c *CLI) pagesHandler() {}

func (c *CLI) createLoginForm(ctx *actor.Context) *tview.Form {
	loginForm := tview.NewForm().
		AddInputField("Email", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil)

	loginForm.AddButton("Login", func() {
		emailField := loginForm.GetFormItemByLabel("Email").(*tview.InputField)
		email := emailField.GetText()

		passwordField := loginForm.GetFormItemByLabel("Password").(*tview.InputField)
		password := passwordField.GetText()

		loadingView := tview.NewTextView().
			SetText("Logging in...").
			SetTextAlign(tview.AlignCenter)
		c.application.pages.AddAndSwitchToPage("loading", loadingView, true)
		go func() {
			resp := ctx.Request(c.serverPID, &Proto.LoginUser{
				Email:    email,
				Password: password,
			}, time.Second*2)

			val, err := resp.Result()
			var resultText string
			if err != nil {
				Server.Logger.Error("Can't do the request!", "err: ", err)
				resultText = "Login failed due to request error" + err.Error()
			} else if reflect.TypeOf(val) == reflect.TypeOf(&Proto.Accept{}) {
				resultText = "Login successful!"
			} else {
				resultText = "Login failed!"
			}

			c.application.app.QueueUpdateDraw(func() {
				textView := tview.NewTextView().
					SetText(resultText).
					SetTextAlign(tview.AlignCenter)

				c.application.pages.AddAndSwitchToPage("result", textView, true)
			})
		}()

	})
	loginForm.SetBorder(true).SetTitle("Login")
	return loginForm
}
