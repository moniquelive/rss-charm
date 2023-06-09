package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/net/html/charset"
)

type rssMarkdown string

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type readModel struct {
	viewport   viewport.Model
	title, url string
	loaded     bool
}

func newRead(title, url string) *readModel {
	vp := viewport.New(windowWidth, windowHeight-2)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return &readModel{
		viewport: vp,
		title:    title,
		url:      url,
		loaded:   false,
	}
}

func (rm *readModel) Init() tea.Cmd {
	return rm.downloadURL()
}

func (rm *readModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case rssMarkdown:
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(windowWidth),
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
			glamour.WithWordWrap(windowWidth),
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
			return models[choose], models[choose].Init()
		}
	}
	var cmd tea.Cmd
	rm.viewport, cmd = rm.viewport.Update(msg)
	return rm, cmd
}

func (rm *readModel) downloadURL() func() tea.Msg {
	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}
		res, err := c.Get(rm.url)
		if err != nil {
			return errMsg{err}
		}
		defer res.Body.Close()
		rm.loaded = true

		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return errMsg{err}
		}
		body := string(bytes)
		if !strings.HasPrefix(body, "<?xml version") && !strings.Contains(body, "2.0") {
			body = `<?xml version="1.0" encoding="iso-8859-1"?>` + body
		}
		var rss Rss
		decoder := xml.NewDecoder(strings.NewReader(body))
		decoder.CharsetReader = charset.NewReaderLabel
		if err := decoder.Decode(&rss); err != nil {
			return errMsg{err}
		}
		return rssMarkdown(
			fmt.Sprintf("# %s\n\n%s", strings.TrimSpace(rss.Channel.Description), formatItems(rss.Channel.Item)),
		)
	}
}

func formatItems(items Items) string {
	converter := md.NewConverter("", true, nil)
	formatted := ""
	for _, item := range items {
		contents := ""
		if description, _ := converter.ConvertString(strings.TrimSpace(item.Description)); description != "" {
			contents += fmt.Sprintf("- %s\n", description)
		}
		contents += strings.TrimSpace(item.Link)
		formatted += "## " + strings.TrimSpace(item.Title) + "\n\n"
		formatted += contents + "\n\n"
	}
	return formatted
}

func (rm *readModel) View() string {
	if rm.loaded {
		return rm.viewport.View() + rm.helpView()
	} else {
		return fmt.Sprintf("Loading %q...", rm.url)
	}
}

func (rm *readModel) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}
