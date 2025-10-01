package tui

import (
	"mdns-browser/internal/data"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
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
	vp          viewport.Model
	addCh       chan data.ListItem
	spinnerTick tea.Cmd
	listWidth   int
	vpWidth     int
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
		totalWidth := msg.Width - h
		m.listWidth = totalWidth * 2 / 3
		m.vpWidth = totalWidth / 3
		m.list.SetSize(m.listWidth, msg.Height-v)
		m.vp.Width = m.vpWidth
		m.vp.Height = msg.Height - v
		for _, item := range m.list.Items() {
			li, ok := item.(data.ListItem)
			if !ok {
				continue
			}
			li.MaxWidth = m.listWidth
		}
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	case addItemMsg:
		listItem := data.ListItem(msg)
		listItem.MaxWidth = m.listWidth
		currentItems := m.list.Items()
		idx := slices.IndexFunc(currentItems, func(it list.Item) bool {
			li, ok := it.(data.ListItem)
			if !ok {
				return false
			}
			return strings.EqualFold(li.Name, listItem.Name)
		})
		numberOfItems := len(currentItems)
		if idx == -1 {
			if numberOfItems == 0 {
				m.vp.SetContent(listItem.Details())
			}
			m.list.InsertItem(numberOfItems, listItem)
		}
		// keep listening
		return m, listenForItems(m.addCh)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	listView := lipgloss.NewStyle().Width(m.listWidth).Render(m.list.View())
	vpView := lipgloss.NewStyle().Width(m.vpWidth).Render(m.vp.View())
	return lipgloss.JoinHorizontal(lipgloss.Top, listView, vpView)
}

func Tui(opts ListOpts) tea.Model {
	var items []list.Item
	listDelegate := list.NewDefaultDelegate()
	l := list.New(items, listDelegate, 0, 0)
	l.Styles.TitleBar.PaddingLeft(5)
	l.SetSpinner(spinner.MiniDot)
	tick := l.StartSpinner()
	vp := viewport.New(0, 0)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)
	m := model{
		list:        l,
		addCh:       opts.AddCh,
		spinnerTick: tick,
		vp:          vp,
	}

	m.list.Title = opts.Title
	return m
}
