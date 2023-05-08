package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html/charset"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const width = 78

type rssMarkdown string

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type readModel struct {
	viewport   viewport.Model
	err        error
	title, url string
}

func newRead(title, url string) *readModel {
	vp := viewport.New(width, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return &readModel{
		viewport: vp,
		title:    title,
		url:      url,
	}
}

func (rm readModel) Init() tea.Cmd {
	return nil
}

func (rm readModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case rssMarkdown:
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(width),
		)
		if err != nil {
			return rm, nil
		}

		str, err := renderer.Render(string(msg))
		if err != nil {
			return rm, nil
		}
		rm.viewport.SetContent(str)

		return rm, nil
	case errMsg:
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(width),
		)
		if err != nil {
			return rm, nil
		}
		str, err := renderer.Render(msg.Error())
		if err != nil {
			return rm, nil
		}
		rm.viewport.SetContent(str)
		return rm, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			models[read] = rm
			return models[choose], nil
		}
	}
	if msg == nil {
		return rm, rm.downloadURL()
	}
	var cmd tea.Cmd
	rm.viewport, cmd = rm.viewport.Update(msg)
	return rm, cmd
}

func (rm readModel) downloadURL() func() tea.Msg {
	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}
		res, err := c.Get(rm.url)
		if err != nil {
			return errMsg{err}
		}
		defer res.Body.Close()

		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return errMsg{err}
		}
		body := `<?xml version="1.0" encoding="iso-8859-1"?>\n` + string(bytes)
		var rss Rss
		decoder := xml.NewDecoder(strings.NewReader(body))
		decoder.CharsetReader = charset.NewReaderLabel
		if err := decoder.Decode(&rss); err != nil {
			return errMsg{err}
		}
		return rssMarkdown(fmt.Sprintf("# %s\n\n%s", rss.Channel.Description, formatItems(rss.Channel.Item)))
	}
}

func formatItems(items Items) string {
	retval := ""
	for _, item := range items {
		retval += "## " + item.Title + "\n\n"
		retval += item.Description + "\n\n"
	}
	return retval
}

func (rm readModel) View() string {
	return rm.viewport.View() + rm.helpView()
}

func (rm readModel) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}
