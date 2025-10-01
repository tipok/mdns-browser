package data

import "strings"

type ListItem struct {
	Name       string
	Host       string
	AddrV4     string
	AddrV6     string
	Port       int
	Info       string
	InfoFields []string
	MaxWidth   int
}

func truncateString(title string, maxWidth int) string {
	// Account for padding and borders in the title bar
	availableWidth := maxWidth - 10 // Conservative padding estimate
	if availableWidth <= 0 {
		return ""
	}
	if len(title) <= availableWidth {
		return title
	}
	if availableWidth <= 3 {
		return title[:availableWidth]
	}
	return title[:availableWidth-3] + "..."
}

func (i ListItem) Title() string {
	if strings.TrimSpace(i.Name) == "" {
		return truncateString(i.Host, i.MaxWidth)
	}
	return truncateString(i.Name, i.MaxWidth)
}
func (i ListItem) Description() string {
	return truncateString(i.Host, i.MaxWidth)
}
func (i ListItem) FilterValue() string {
	return i.Title()
}
