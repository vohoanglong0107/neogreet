package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app     *tview.Application
	config  *Config
	greeter *Greeter
}

func NewUI(config *Config, greeter *Greeter) *UI {
	ui := &UI{app: tview.NewApplication(), config: config, greeter: greeter}
	return ui
}

func (ui *UI) Draw() {
	ui.greeter.OnDone = func() {
		ui.app.Stop()
	}
	if err := ui.app.SetRoot(ui.DrawContainer(), true).Run(); err != nil {
		panic(err)
	}
}

func (ui *UI) DrawContainer() *tview.Grid {
	logo := ui.config.getLogo()
	info := ui.config.getInfo(&SystemInfo{})
	logoWidth, logoHeigth := getDim(logo)
	infoWidth, infoHeigth := getDim(info)
	var firstRowHeight int
	if logoHeigth > infoHeigth {
		firstRowHeight = logoHeigth
	} else {
		firstRowHeight = infoHeigth
	}

	container := tview.NewGrid()
	container.SetColumns(-1, logoWidth, infoWidth, -1)
	container.SetRows(-1, firstRowHeight, 5, 1, -1)
	container.SetGap(0, 4)

	container.AddItem(
		tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(ui.DrawText(logo), logoHeigth, 0, true).
			AddItem(nil, 0, 1, false), 1, 1, 1, 1, 0, 0, false)
	container.AddItem(
		tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(ui.DrawText(info), infoHeigth, 0, true).
			AddItem(nil, 0, 1, false), 1, 2, 1, 1, 0, 0, false)
	container.AddItem(
		tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false).
			AddItem(ui.DrawLoginForm(), 30, 0, true).
			AddItem(nil, 0, 1, false),
		2, 1, 1, 2, 0, 0, true)
	container.AddItem(
		tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false).
			AddItem(ui.DrawNotice(), 30, 0, true).
			AddItem(nil, 0, 1, false),
		3, 1, 1, 2, 0, 0, false,
	)

	return container
}

func (ui *UI) DrawText(text string) *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			ui.app.Draw()
		})
	textView.SetBackgroundColor(tcell.ColorDefault)

	fmt.Fprintf(textView, "%s", text)
	return textView
}

func (ui *UI) DrawLoginForm() *tview.Form {
	loginForm := tview.NewForm()
	loginForm.SetFieldBackgroundColor(tcell.ColorDefault)
	loginForm.SetFieldTextColor(tcell.ColorLightCyan)
	loginForm.SetBackgroundColor(tcell.ColorDefault)
	loginForm.SetItemPadding(0)
	loginForm.SetBorderPadding(0, 0, 0, 0)

	loginForm.AddFormItem(ui.AddInputField("Username: ", false, ui.greeter.CreateSession))
	ui.greeter.OnRequestInput = func(label string, secret bool) {
		loginForm.AddFormItem(ui.AddInputField(label, secret, ui.greeter.HandleInput))
	}
	ui.greeter.OnError = func() {
		loginForm.Clear(true)
		loginForm.AddFormItem(ui.AddInputField("Username: ", false, ui.greeter.CreateSession))
	}
	return loginForm
}

func (ui *UI) DrawNotice() *tview.TextView {
	noticeBox := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			ui.app.Draw()
		})
	noticeBox.SetBackgroundColor(tcell.ColorDefault)
	noticeBox.SetTextColor(tcell.ColorRed)
	noticeBox.SetTextAlign(tview.AlignCenter)
	ui.greeter.OnNotice = func(notice string) {
		noticeBox.Clear()
		fmt.Fprintf(noticeBox, "%s", notice)
	}
	return noticeBox
}

func (ui *UI) AddInputField(label string, secret bool, handle func(input string)) tview.FormItem {
	input := ""
	inputField := tview.NewInputField()
	inputField.SetLabel(label)
	inputField.SetFieldWidth(40)
	inputField.SetLabelStyle(tcell.StyleDefault.Bold(true))
	inputField.SetFieldStyle(tcell.StyleDefault.Bold(true))
	if secret {
		inputField.SetMaskCharacter('*')
	}
	inputField.SetChangedFunc(func(text string) { input = text })
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyEnter {
			return event
		}
		handle(input)
		return event
	})
	return inputField
}

func getDim(logo string) (int, int) {
	logoLines := strings.Split(logo, "\n")
	heigth := len(logoLines)
	width := math.MinInt
	colorTag, _ := regexp.Compile(`\[.*\]`)
	for _, line := range logoLines {
		rawLine := []rune(string(colorTag.ReplaceAll([]byte(line), []byte(""))))
		// fmt.Println(rawLine, len([]rune(rawLine)))
		if len(rawLine) > width {
			width = len(rawLine)
		}
	}
	return width, heigth
}
