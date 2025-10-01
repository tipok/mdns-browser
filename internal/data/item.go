package data

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ListItem struct {
	Name            string
	Host            string
	AddrV4          string
	AddrV6          string
	Port            int
	Info            string
	InfoFields      []string
	MaxListWidth    int
	MaxDetailsWidth int
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
		return truncateString(i.Host, i.MaxListWidth)
	}
	return truncateString(i.Name, i.MaxListWidth)
}
func (i ListItem) Description() string {
	return truncateString(i.Host, i.MaxListWidth)
}
func (i ListItem) FilterValue() string {
	return i.Title()
}

// Details are used for the details view which is showing all the
//
//	properties of the item as a styled string using lipgloss
func (i ListItem) Details() string {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#04B575"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Padding(0, 1)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B6B")).
		MarginTop(1).
		MarginBottom(1)

	bulletStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	var details []string

	// Title
	details = append(details, titleStyle.Render("🔍 Service Details"))

	// Service details
	if strings.TrimSpace(i.Name) != "" {
		details = append(details, labelStyle.Render("Service Name: ")+valueStyle.Render(i.Name))
	}

	if strings.TrimSpace(i.Host) != "" {
		details = append(details, labelStyle.Render("Host: ")+valueStyle.Render(i.Host))
	}

	if strings.TrimSpace(i.AddrV4) != "" {
		details = append(details, labelStyle.Render("IPv4 Address: ")+valueStyle.Render(i.AddrV4))
	}

	if strings.TrimSpace(i.AddrV6) != "" {
		details = append(details, labelStyle.Render("IPv6 Address: ")+valueStyle.Render(i.AddrV6))
	}

	if i.Port > 0 {
		details = append(details, labelStyle.Render("Port: ")+valueStyle.Render(fmt.Sprintf("%d", i.Port)))
	}

	// Additional information section
	if strings.TrimSpace(i.Info) != "" {
		details = append(details, "")
		details = append(details, sectionStyle.Render("📋 Additional Information"))
		details = append(details, i.Info)
	}

	// Service fields section
	if len(i.InfoFields) > 0 {
		details = append(details, "")
		details = append(details, sectionStyle.Render("🏷️  Service Fields"))
		for _, field := range i.InfoFields {
			if strings.TrimSpace(field) != "" {
				details = append(details, bulletStyle.Render("• ")+field)
			}
		}
	}

	return strings.Join(details, "\n")
}
