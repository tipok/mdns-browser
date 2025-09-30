package tui

import (
	"mdns-browser/internal/data"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ListOpts struct {
	Title string
	AddCh chan data.ListItem
}

// message carrying a new ListItem
type addItemMsg data.ListItem

// command that waits for the next ListItem from a channel
func listenForItems(ch <-chan data.ListItem) tea.Cmd {
	return func() tea.Msg {
		it := <-ch
		return addItemMsg(it)
	}
}

type model struct {
	list        list.Model
	addCh       chan data.ListItem
	spinnerTick tea.Cmd
}

func (m model) Init() tea.Cmd {
	return tea.Batch(listenForItems(m.addCh), m.spinnerTick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	case addItemMsg:
		listItem := data.ListItem(msg)
		currentItems := m.list.Items()
		idx := slices.IndexFunc(currentItems, func(it list.Item) bool {
			li, ok := it.(data.ListItem)
			if !ok {
				return false
			}
			return strings.EqualFold(li.Name, listItem.Name)
		})

		if idx == -1 {
			m.list.InsertItem(len(currentItems), listItem)
		}
		// keep listening
		return m, listenForItems(m.addCh)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func List(opts ListOpts) tea.Model {
	var items []list.Item
	listDelegate := list.NewDefaultDelegate()
	l := list.New(items, listDelegate, 0, 0)
	l.Styles.TitleBar.PaddingLeft(5)
	l.SetSpinner(spinner.MiniDot)
	tick := l.StartSpinner()
	m := model{
		list:        l,
		addCh:       opts.AddCh,
		spinnerTick: tick,
	}

	m.list.Title = opts.Title
	return m
}
