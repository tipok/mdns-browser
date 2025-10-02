package data

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
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

func wrapString(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	// Use visual width instead of byte length
	if runewidth.StringWidth(text) <= maxWidth {
		return []string{text}
	}

	var lines []string
	runes := []rune(text)

	for len(runes) > 0 {
		currentWidth := 0
		breakPoint := 0
		lastSpacePoint := -1

		// Find where to break based on visual width
		for i, r := range runes {
			charWidth := runewidth.RuneWidth(r)
			if currentWidth+charWidth > maxWidth {
				break
			}
			currentWidth += charWidth
			breakPoint = i + 1

			// Track the last space for better breaking
			if r == ' ' {
				lastSpacePoint = i
			}
		}

		// If we found a space and it's reasonable, break there
		if lastSpacePoint != -1 && lastSpacePoint > breakPoint/2 {
			breakPoint = lastSpacePoint
		}

		// Extract the line
		if breakPoint == 0 {
			breakPoint = 1 // Ensure we make progress
		}

		line := string(runes[:breakPoint])
		lines = append(lines, line)

		// Remove the processed part and any leading spaces
		runes = runes[breakPoint:]
		for len(runes) > 0 && runes[0] == ' ' {
			runes = runes[1:]
		}
	}

	return lines
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

// addWrappedValue adds a label-value pair with wrapping support for long values
func (i ListItem) addWrappedValue(details []string, labelStyle lipgloss.Style, valueStyle lipgloss.Style, label, value string) []string {
	if strings.TrimSpace(value) == "" {
		return details
	}

	// Account for label width and padding when calculating available width for value
	labelWidth := runewidth.StringWidth(label)
	availableWidth := i.MaxDetailsWidth - labelWidth - 4 // Conservative padding estimate

	wrappedLines := wrapString(value, availableWidth)

	// The first line includes the label
	if len(wrappedLines) > 0 {
		details = append(details, labelStyle.Render(label)+valueStyle.Render(wrappedLines[0]))

		// Later lines are indented to align with the value
		indent := strings.Repeat(" ", labelWidth)
		for _, line := range wrappedLines[1:] {
			details = append(details, indent+valueStyle.Render(line))
		}
	}

	return details
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

	// Service details with wrapping
	details = i.addWrappedValue(details, labelStyle, valueStyle, "Service Name: ", i.Name)
	details = i.addWrappedValue(details, labelStyle, valueStyle, "Host: ", i.Host)
	details = i.addWrappedValue(details, labelStyle, valueStyle, "IPv4 Address: ", i.AddrV4)
	details = i.addWrappedValue(details, labelStyle, valueStyle, "IPv6 Address: ", i.AddrV6)

	if i.Port > 0 {
		details = i.addWrappedValue(details, labelStyle, valueStyle, "Port: ", fmt.Sprintf("%d", i.Port))
	}

	// Additional information section
	if strings.TrimSpace(i.Info) != "" {
		details = append(details, "")
		details = append(details, sectionStyle.Render("📋 Additional Information"))

		// Wrap the info text
		wrappedInfo := wrapString(i.Info, i.MaxDetailsWidth-4) // Account for padding
		details = append(details, wrappedInfo...)
	}

	// Service fields section
	hasServiceFields := len(i.InfoFields) > 0 && len(i.InfoFields) != 1 || strings.TrimSpace(i.InfoFields[0]) != ""
	if hasServiceFields {
		details = append(details, "")
		details = append(details, sectionStyle.Render("🧰 Service Fields"))
		for _, field := range i.InfoFields {
			if strings.TrimSpace(field) != "" {
				// Wrap individual service fields
				wrappedField := wrapString(field, i.MaxDetailsWidth-6) // Account for bullet and padding
				if len(wrappedField) > 0 {
					details = append(details, bulletStyle.Render("• ")+wrappedField[0])
					for _, line := range wrappedField[1:] {
						details = append(details, "  "+line) // Indent continuation lines
					}
				}
			}
		}
	}

	return strings.Join(details, "\n")
}
