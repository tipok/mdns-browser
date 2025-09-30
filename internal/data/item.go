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
}

func (i ListItem) Title() string {
	if strings.TrimSpace(i.Name) == "" {
		return i.Host
	}
	return i.Name
}
func (i ListItem) Description() string {
	return i.Info
}
func (i ListItem) FilterValue() string {
	return i.Title()
}
