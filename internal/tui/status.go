package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"geminictl/internal/registry"
	"geminictl/internal/scanner"
)

// Style definitions
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	warning   = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF5555"}

	listStyle = lipgloss.NewStyle().
		MarginRight(2).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle)

	detailsStyle = lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtle)

	highlightStyle = lipgloss.NewStyle().Foreground(highlight)

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(highlight)

	orphanStyle = lipgloss.NewStyle().
		Foreground(warning).
		Italic(true)
)

type projectView struct {
	ID       string
	Path     string
	IsOrphan bool
	Sessions []scanner.Session
	IsScanning bool
}

type Model struct {
	Projects []projectView
	Cursor   int
	Selected int
	Width    int
	Height   int
	Err      error
	
	scanner  *scanner.Scanner
	registry *registry.Registry
}

// Internal message to carry the channel along with the result
type resolutionPacket struct {
	res scanner.Resolution
	ch  <-chan scanner.Resolution
}

type ScanFinishedMsg struct{}

func NewModel(scanned []scanner.ProjectData, reg *registry.Registry, sc *scanner.Scanner) Model {
	var projects []projectView
	for _, p := range scanned {
		path, isOrphan, err := reg.GetProjectPath(p.ID)
		if err != nil {
			path = p.ID // Fallback to ID
		}
		projects = append(projects, projectView{
			ID:       p.ID,
			Path:     path,
			IsOrphan: isOrphan,
			Sessions: p.Sessions,
			IsScanning: path == p.ID, // Mark as scanning if unresolved
		})
	}

	return Model{
		Projects: projects,
		Selected: 0,
		scanner:  sc,
		registry: reg,
	}
}

func (m Model) Init() tea.Cmd {
	return m.startScanningCmd()
}

func (m Model) startScanningCmd() tea.Cmd {
	var unknownIDs []string
	for _, p := range m.Projects {
		if p.IsScanning {
			unknownIDs = append(unknownIDs, p.ID)
		}
	}

	if len(unknownIDs) == 0 {
		return nil
	}

	c := m.scanner.ResolveBackground(unknownIDs)

	return func() tea.Msg {
		return waitForResolution(c)
	}
}

func waitForResolution(c <-chan scanner.Resolution) tea.Msg {
	res, ok := <-c
	if !ok {
		return ScanFinishedMsg{}
	}
	return resolutionPacket{res, c}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Projects)-1 {
				m.Cursor++
			}
		case "enter", " ":
			m.Selected = m.Cursor
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case resolutionPacket:
		// Update model
		for i, p := range m.Projects {
			if p.ID == msg.res.Hash {
				m.Projects[i].Path = msg.res.Path
				m.Projects[i].IsOrphan = false
				m.Projects[i].IsScanning = false
				
				// Auto-save (Persistence Task)
				m.registry.AddProject(msg.res.Hash, msg.res.Path)
				_ = m.registry.Save()
				break
			}
		}
		// Continue listening
		return m, func() tea.Msg {
			return waitForResolution(msg.ch)
		}

	case ScanFinishedMsg:
		// Turn off scanning indicators for any that remain unresolved
		for i := range m.Projects {
			m.Projects[i].IsScanning = false
		}
	}

	return m, nil
}

func (m Model) View() string {
	if len(m.Projects) == 0 {
		return "No projects found in ~/.gemini/tmp"
	}

	var sidebar strings.Builder
	sidebar.WriteString(titleStyle.Render("Projects") + "\n\n")

	for i, p := range m.Projects {
		cursor := "  "
		if m.Cursor == i {
			cursor = "> "
		}

		style := lipgloss.NewStyle()
		if m.Selected == i {
			style = style.Foreground(special)
		}
		
		label := p.Path
		
		// Logic for display
		if p.IsScanning {
			label += " ..." // Indicator
		} else if p.IsOrphan || p.Path == p.ID {
			style = style.Inherit(orphanStyle)
			if len(label) > 12 && p.Path == p.ID {
				label = label[:12] + "..."
			}
			label += " [Orphan]"
		}

		sidebar.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(label)))
	}

	var main strings.Builder
	if m.Selected < len(m.Projects) {
		p := m.Projects[m.Selected]
		main.WriteString(titleStyle.Render(fmt.Sprintf("Sessions for %s", p.Path)) + "\n\n")

		if len(p.Sessions) == 0 {
			main.WriteString("No sessions found.")
		} else {
			for _, s := range p.Sessions {
				id := s.ID
				if len(id) > 8 {
					id = id[:8]
				}
				main.WriteString(fmt.Sprintf("• %s | %d messages | last: %s\n",
					highlightStyle.Render(id),
					s.MessageCount,
					formatRelativeTime(s.LastUpdate)))
			}
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		listStyle.Render(sidebar.String()),
		detailsStyle.Render(main.String()),
	)
}

func formatRelativeTime(t time.Time) string {
	duration := time.Since(t)
	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	default:
		return t.Format("2006-01-02")
	}
}
