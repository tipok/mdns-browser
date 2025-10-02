package tui

import (
	"mdns-browser/internal/data"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	list         list.Model
	vp           viewport.Model
	help         help.Model
	addCh        chan data.ListItem
	spinnerTick  tea.Cmd
	listWidth    int
	vpWidth      int
	focusedView  int  // 0 = list, 1 = viewport
	showFullHelp bool // Whether to show full help or short help
}

// keyMap defines key bindings for our TUI
type keyMap struct {
	// Common keys
	Quit       key.Binding
	Tab        key.Binding
	HelpToggle key.Binding

	// List-specific keys
	Up    key.Binding
	Down  key.Binding
	Slash key.Binding

	// Viewport-specific keys
	ScrollUp   key.Binding
	ScrollDown key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	GoToTop    key.Binding
	GoToBottom key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	var bindings []key.Binding
	bindings = append(bindings, k.Quit, k.Tab, k.HelpToggle)

	// Add context-specific keys
	if len(k.Up.Keys()) > 0 {
		bindings = append(bindings, k.Up, k.Down)
	}
	if len(k.ScrollUp.Keys()) > 0 {
		bindings = append(bindings, k.ScrollUp, k.ScrollDown)
	}

	return bindings
}

// FullHelp returns keybindings for the expanded help view
func (k keyMap) FullHelp() [][]key.Binding {
	commonKeys := []key.Binding{k.Quit, k.Tab, k.HelpToggle}

	// List-specific keys
	if len(k.Up.Keys()) > 0 {
		return [][]key.Binding{
			commonKeys,
			{k.Up, k.Down, k.Slash},
		}
	}

	// Viewport-specific keys
	if len(k.ScrollUp.Keys()) > 0 {
		return [][]key.Binding{
			commonKeys,
			{k.ScrollUp, k.ScrollDown, k.PageUp, k.PageDown},
			{k.GoToTop, k.GoToBottom},
		}
	}

	return [][]key.Binding{commonKeys}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
	HelpToggle: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Slash: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter list"),
	),
	ScrollUp: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "scroll up"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "scroll down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("ctrl+k", "pgup"),
		key.WithHelp("ctrl+k/pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("ctrl+j", "pgdown"),
		key.WithHelp("ctrl+j/pgdn", "page down"),
	),
	GoToTop: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "go to top"),
	),
	GoToBottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "go to bottom"),
	),
}

// contextualKeyMap creates a keyMap based on the current focused view
func (m model) contextualKeyMap() keyMap {
	if m.focusedView == 0 { // list focused
		return keyMap{
			Quit:       keys.Quit,
			Tab:        keys.Tab,
			HelpToggle: keys.HelpToggle,
			Up:         keys.Up,
			Down:       keys.Down,
			Slash:      keys.Slash,
		}
	} else { // viewport focused
		return keyMap{
			Quit:       keys.Quit,
			Tab:        keys.Tab,
			HelpToggle: keys.HelpToggle,
			ScrollUp:   keys.ScrollUp,
			ScrollDown: keys.ScrollDown,
			PageUp:     keys.PageUp,
			PageDown:   keys.PageDown,
			GoToTop:    keys.GoToTop,
			GoToBottom: keys.GoToBottom,
		}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(listenForItems(m.addCh), m.spinnerTick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		switch k {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			// Switch between list and viewport focus
			m.focusedView = (m.focusedView + 1) % 2
			return m, nil
		case "?":
			// Toggle help display
			m.showFullHelp = !m.showFullHelp
			m.help.ShowAll = m.showFullHelp
			return m, nil
		case "j", "down":
			if m.focusedView == 1 { // viewport focused
				m.vp.ScrollDown(1)
				return m, nil
			}
			// Otherwise handle in the list below
		case "k", "up":
			if m.focusedView == 1 { // viewport focused
				m.vp.ScrollUp(1)
				return m, nil
			}
			// Otherwise handle in the list below
		case "ctrl+j", "pgdown":
			if m.focusedView == 1 {
				m.vp.HalfPageDown()
				return m, nil
			}
		case "ctrl+k", "pgup":
			if m.focusedView == 1 {
				m.vp.HalfPageUp()
				return m, nil
			}
		case "g":
			if m.focusedView == 1 {
				m.vp.GotoTop()
				return m, nil
			}
		case "G":
			if m.focusedView == 1 {
				m.vp.GotoBottom()
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		totalWidth := msg.Width - h

		// Reserve space for help at the bottom
		helpHeight := 3
		availableHeight := msg.Height - v - helpHeight

		m.listWidth = totalWidth * 2 / 3
		m.vpWidth = totalWidth / 3
		m.list.SetSize(m.listWidth, availableHeight)
		m.vp.Width = m.vpWidth
		m.vp.Height = availableHeight
		m.help.Width = msg.Width
		for _, item := range m.list.Items() {
			li, ok := item.(data.ListItem)
			if !ok {
				continue
			}
			li.MaxListWidth = m.listWidth
			li.MaxDetailsWidth = m.vpWidth
		}
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	case addItemMsg:
		listItem := data.ListItem(msg)
		listItem.MaxListWidth = m.listWidth
		listItem.MaxDetailsWidth = m.vpWidth
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

	// Only update the list if it's focused
	if m.focusedView == 0 {
		oldIndex := m.list.Index()
		m.list, cmd = m.list.Update(msg)

		// Update viewport content when list selection changes
		if m.list.Index() != oldIndex && len(m.list.Items()) > 0 {
			if selectedItem := m.list.SelectedItem(); selectedItem != nil {
				if listItem, ok := selectedItem.(data.ListItem); ok {
					m.vp.SetContent(listItem.Details())
					m.vp.GotoTop() // Reset the scroll position when switching items
				}
			}
		}
	}

	return m, cmd
}

func (m model) View() string {
	// Style focused and unfocused views differently
	listStyle := lipgloss.NewStyle().Width(m.listWidth)
	vpStyle := lipgloss.NewStyle().Width(m.vpWidth)

	if m.focusedView == 0 { // list focused
		listStyle = listStyle.BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4"))
		vpStyle = vpStyle.BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666"))
	} else { // viewport focused
		listStyle = listStyle.BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666"))
		vpStyle = vpStyle.BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4"))
	}

	listView := listStyle.Render(m.list.View())
	vpView := vpStyle.Render(m.vp.View())
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, listView, vpView)

	// Create a contextual help view
	contextualKeys := m.contextualKeyMap()
	helpView := m.help.View(contextualKeys)

	return lipgloss.JoinVertical(lipgloss.Left, mainView, helpView)
}

func Tui(opts ListOpts) tea.Model {
	var items []list.Item
	listDelegate := list.NewDefaultDelegate()
	l := list.New(items, listDelegate, 0, 0)
	l.Styles.TitleBar.PaddingLeft(5)
	l.SetSpinner(spinner.MiniDot)

	// Disable the built-in help for the list since we'll handle it ourselves
	l.SetShowHelp(false)

	tick := l.StartSpinner()
	vp := viewport.New(0, 0)
	vp.SetContent(`No service selected.`)

	h := help.New()
	h.ShowAll = true // Start with full help to show more keys

	m := model{
		list:         l,
		addCh:        opts.AddCh,
		spinnerTick:  tick,
		vp:           vp,
		help:         h,
		focusedView:  0,    // Start with list focused
		showFullHelp: true, // Start with full help
	}

	m.list.Title = opts.Title
	return m
}
