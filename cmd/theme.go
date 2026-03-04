package cmd

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
)

func ApplyColorScheme() fang.ColorSchemeFunc {
	return func (lipgloss.LightDarkFunc) fang.ColorScheme {
		return EchoTheme()
	}
}

type ColorScheme struct {
	Base           lipgloss.Style
	Title          lipgloss.Style
	Description    lipgloss.Style
	Codeblock      lipgloss.Style
	Program        lipgloss.Style
	DimmedArgument lipgloss.Style
	Comment        lipgloss.Style
	Flag           lipgloss.Style
	FlagDefault    lipgloss.Style
	Command        lipgloss.Style
	QuotedString   lipgloss.Style
	Argument       lipgloss.Style
	Help           lipgloss.Style
	Dash           lipgloss.Style
	ErrorHeader    lipgloss.Style
	ErrorDetails   lipgloss.Style
}

func EchoTheme() fang.ColorScheme {
	// Core palette (Echo-inspired)
	const (
		bgDark   = "#0F172A" // deep navy
		cyan     = "#22D3EE" // bright cyan accent
		blue     = "#38BDF8"
		green    = "#4ADE80"
		magenta  = "#C084FC"
		white    = "#F8FAFC"
		gray     = "#94A3B8"
		darkGray = "#64748B"
		red      = "#EF4444"
	)

	return fang.ColorScheme{
		Base: lipgloss.Color(white),
		Title: lipgloss.Color(cyan),
		Command: lipgloss.Color(cyan),
	}
}
