package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	choose = iota
	read
)

var models []tea.Model

type item struct {
	title, url string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.url }
func (i item) FilterValue() string { return i.title }

var feeds = []list.Item{
	item{title: "Home UOL", url: "http://rss.home.uol.com.br/index.xml"},
	item{title: "UOL Notícias", url: "http://rss.uol.com.br/feed/noticias.xml"},
	item{title: "UOL Tecnologia", url: "http://tecnologia.uol.com.br/ultnot/index.xml"},
	item{title: "UOL Economia", url: "http://rss.uol.com.br/feed/economia.xml"},
	item{title: "UOL Esporte", url: "http://rss.esporte.uol.com.br/ultimas/index.xml"},
	item{title: "UOL Esporte - Basquete", url: "http://esporte.uol.com.br/basquete/ultimas/index.xml"},
	item{title: "UOL Esporte - Futebol", url: "http://esporte.uol.com.br/futebol/ultimas/index.xml"},
	item{title: "UOL Esporte - Tênis", url: "http://esporte.uol.com.br/tenis/ultimas/index.xml"},
	item{title: "UOL Jogos", url: "http://rss.uol.com.br/feed/jogos.xml"},
	item{title: "UOL Cinema", url: "http://cinema.uol.com.br/ultnot/index.xml"},
	item{title: "UOL Música", url: "http://musica.uol.com.br/ultnot/index.xml"},
	item{title: "UOL Vestibular", url: "http://rss.uol.com.br/feed/vestibular.xml"},
	item{title: "UOL Carros", url: "http://rss.carros.uol.com.br/ultnot/index.xml"},
	item{title: "UOL Eleições 2022", url: "https://rss.uol.com.br/feed/eleicoes2022.xml"},
	item{title: "Newsletter - Para Começar o Dia", url: "https://rss.uol.com.br/feed/comecar-o-dia.xml"},

	item{title: "Atlético-MG", url: "http://rss.esporte.uol.com.br/futebol/clubes/atleticomg.xml"},
	item{title: "Atlético-PR", url: "http://rss.esporte.uol.com.br/futebol/clubes/atleticopr.xml"},
	item{title: "Bahia", url: "http://esporte.uol.com.br/futebol/clubes/bahia.xml"},
	item{title: "Botafogo", url: "http://esporte.uol.com.br/futebol/clubes/botafogo.xml"},
	item{title: "Ceará", url: "http://rss.esporte.uol.com.br/futebol/clubes/ceara.xml"},
	item{title: "Corinthians", url: "http://rss.esporte.uol.com.br/futebol/clubes/corinthians.xml"},
	item{title: "Coritiba", url: "http://rss.esporte.uol.com.br/futebol/clubes/coritiba.xml"},
	item{title: "Cruzeiro", url: "http://rss.esporte.uol.com.br/futebol/clubes/cruzeiro.xml"},
	item{title: "Flamengo", url: "http://rss.esporte.uol.com.br/futebol/clubes/flamengo.xml"},
	item{title: "Fluminense", url: "http://rss.esporte.uol.com.br/futebol/clubes/fluminense.xml"},
	item{title: "Fortaleza", url: "http://rss.esporte.uol.com.br/futebol/clubes/fortaleza.xml"},
	item{title: "Grêmio", url: "http://rss.esporte.uol.com.br/futebol/clubes/gremio.xml"},
	item{title: "Internacional", url: "http://rss.esporte.uol.com.br/futebol/clubes/internacional.xml"},
	item{title: "Náutico", url: "http://rss.esporte.uol.com.br/futebol/clubes/nautico.xml"},
	item{title: "Palmeiras", url: "http://rss.esporte.uol.com.br/futebol/clubes/palmeiras.xml"},
	item{title: "Santa Cruz", url: "http://rss.esporte.uol.com.br/futebol/clubes/santacruz.xml"},
	item{title: "Santos", url: "http://rss.esporte.uol.com.br/futebol/clubes/santos.xml"},
	item{title: "São Paulo", url: "http://rss.esporte.uol.com.br/futebol/clubes/saopaulo.xml"},
	item{title: "Sport", url: "http://rss.esporte.uol.com.br/futebol/clubes/sport.xml"},
	item{title: "Vasco", url: "http://rss.esporte.uol.com.br/futebol/clubes/vasco.xml"},
	item{title: "Vitória", url: "http://esporte.uol.com.br/futebol/clubes/vitoria.xml"},
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i := m.list.SelectedItem().(item)
			readM := newRead(i.title, i.url)
			models[choose] = m
			models[read] = readM
			return models[read].Update(nil)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	m := model{list: list.New(feeds, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "UOL RSS Feeds"

	models = []tea.Model{&m, nil}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
