package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/spinner"
	"geminictl/internal/cache"
	"geminictl/internal/scanner"
	"sort"
)

// Style definitions
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	warning   = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF5555"}

	listStyle = lipgloss.NewStyle().
		MarginRight(1).
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

	orphanTagStyle = lipgloss.NewStyle().
			Foreground(warning).
			Bold(true)
)

func collapseHome(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if strings.HasPrefix(path, home) {
		return "~" + strings.TrimPrefix(path, home)
	}
	return path
}

func truncateMiddle(s string, max int) string {
	if len(s) <= max || max < 5 {
		return s
	}
	half := (max - 3) / 2
	return s[:half] + "..." + s[len(s)-half:]
}

type projectView struct {
	ID         string
	Path       string
	IsOrphan   bool
	Sessions   []scanner.Session
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
	cache    *cache.Cache
	spinner  spinner.Model
}

// Internal message to carry the channel along with the result
type resolutionPacket struct {
	res scanner.Resolution
	ch  <-chan scanner.Resolution
}

type ScanFinishedMsg struct{}

func NewModel(scanned []scanner.ProjectData, c *cache.Cache, sc *scanner.Scanner) Model {
	var projects []projectView
	for _, p := range scanned {
		path, ok := c.Get(p.ID)
		if !ok {
			path = p.ID // Initially just the ID
		}
		projects = append(projects, projectView{
			ID:         p.ID,
			Path:       path,
			IsOrphan:   false, // Initial state, will be updated in Phase 2
			Sessions:   p.Sessions,
			IsScanning: !ok, // Mark as scanning if not in cache
		})
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(highlight)

	m := Model{
		Projects: projects,
		Selected: 0,
		scanner:  sc,
		cache:    c,
		spinner:  s,
	}
	m.sortProjects()
	return m
}

func (m *Model) sortProjects() {
	var selectedID string
	if len(m.Projects) > 0 {
		selectedID = m.Projects[m.Selected].ID
	}
	var cursorID string
	if len(m.Projects) > 0 {
		cursorID = m.Projects[m.Cursor].ID
	}

	sort.Slice(m.Projects, func(i, j int) bool {
		return m.Projects[i].Path < m.Projects[j].Path
	})

	if selectedID != "" {
		for i, p := range m.Projects {
			if p.ID == selectedID {
				m.Selected = i
				break
			}
		}
	}
	if cursorID != "" {
		for i, p := range m.Projects {
			if p.ID == cursorID {
				m.Cursor = i
				break
			}
		}
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.startScanningCmd())
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
	var cmd tea.Cmd
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

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case resolutionPacket:
		for i, p := range m.Projects {
			if p.ID == msg.res.Hash {
				m.Projects[i].Path = msg.res.Path
				m.Projects[i].IsOrphan = false
				m.Projects[i].IsScanning = false

				m.cache.Set(msg.res.Hash, msg.res.Path)
				_ = m.cache.Save()
				break
			}
		}
		m.sortProjects()
		return m, func() tea.Msg {
			return waitForResolution(msg.ch)
		}

	case ScanFinishedMsg:
		for i := range m.Projects {
			if m.Projects[i].IsScanning {
				m.Projects[i].IsScanning = false
				// Save empty path for Unlocated (Persistence Task)
				m.cache.Set(m.Projects[i].ID, "")
				_ = m.cache.Save()
			}
		}
	}

	return m, nil
}

func (m Model) isScanningGlobal() bool {
	for _, p := range m.Projects {
		if p.IsScanning {
			return true
		}
	}
	return false
}

func (m Model) View() string {
	if len(m.Projects) == 0 {
		return "No projects found in ~/.gemini/tmp"
	}

	sidebarWidth := (m.Width / 2) - 2
	mainWidth := m.Width - sidebarWidth - 6

	var sidebar strings.Builder
	sidebar.WriteString(titleStyle.Render("Projects") + "\n")
	if m.isScanningGlobal() {
		text := "Resolving directories... " + m.spinner.View()
		padding := sidebarWidth - lipgloss.Width(text) - 4
		if padding < 0 {
			padding = 0
		}
		sidebar.WriteString(strings.Repeat(" ", padding) + lipgloss.NewStyle().Foreground(subtle).Render(text) + "\n")
	} else {
		sidebar.WriteString("\n")
	}

	for i, p := range m.Projects {
		cursor := "  "
		if m.Cursor == i {
			cursor = "> "
		}

		style := lipgloss.NewStyle()
		if m.Selected == i {
			style = style.Foreground(special)
		}

		displayPath := collapseHome(p.Path)
		if displayPath == p.ID && len(displayPath) > 12 {
			displayPath = displayPath[:12] + "..."
		}

		availableWidth := sidebarWidth - 6
		if p.IsScanning {
			availableWidth -= 2
		}
		if (p.IsOrphan || p.Path == p.ID) && !p.IsScanning {
			availableWidth -= 9
		}

		displayPath = truncateMiddle(displayPath, availableWidth)

		suffix := ""
		if p.IsScanning {
			suffix = " " + m.spinner.View()
		} else if p.IsOrphan || p.Path == p.ID {
			suffix = " " + orphanTagStyle.Render("[Orphan]") // Will update to [Unlocated] in Phase 2
		}

		sidebar.WriteString(fmt.Sprintf("%s%s%s\n", cursor, style.Render(displayPath), suffix))
	}

	var main strings.Builder
	if m.Selected < len(m.Projects) {
		p := m.Projects[m.Selected]
		main.WriteString(titleStyle.Render(fmt.Sprintf("Sessions for %s", collapseHome(p.Path))) + "\n\n")

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
		listStyle.Width(sidebarWidth).Render(sidebar.String()),
		detailsStyle.Width(mainWidth).Render(main.String()),
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
