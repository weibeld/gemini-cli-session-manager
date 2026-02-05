package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"geminictl/internal/cache"
	"geminictl/internal/gemini"
	"geminictl/internal/scanner"
	"sort"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProjectStatus represents the resolution state of a project.
type ProjectStatus int

const (
	StatusScanning ProjectStatus = iota
	StatusValid
	StatusUnlocated
	StatusOrphaned
)

type Focus int

const (
	FocusProjects Focus = iota
	FocusSessions
)

type Mode int

const (
	ModeNav Mode = iota
	ModeInputPath
	ModeConfirmDelete
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

	tagStyle = lipgloss.NewStyle()

	strikethroughStyle = lipgloss.NewStyle().
				Strikethrough(true)

	promptStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1).
			Width(60)
)

// UI Helpers
func renderCursor(focused bool, current bool) string {
	if focused && current {
		return lipgloss.NewStyle().Foreground(special).Render("> ")
	}
	return "  "
}

func renderHash(id string) string {
	shortID := id
	if len(id) > 8 {
		shortID = id[:8]
	}
	return highlightStyle.Render(fmt.Sprintf("[%s]", shortID))
}

func getRowStyle(selected bool) lipgloss.Style {
	style := lipgloss.NewStyle()
	if selected {
		style = style.Foreground(special)
	}
	return style
}

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
	ID       string
	Path     string
	Status   ProjectStatus
	Sessions []scanner.Session
}

type Model struct {
	Projects      []projectView
	Cursor        int
	Selected      int
	SessionCursor int
	Focus         Focus
	Mode          Mode
	Width         int
	Height        int
	Err           error

	scanner   *scanner.Scanner
	cache     *cache.Cache
	spinner   spinner.Model
	textInput textinput.Model
	prompt    string
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
		projects = append(projects, deriveProjectView(p.ID, p.Sessions, c))
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(highlight)

	ti := textinput.New()
	ti.Placeholder = "Absolute path..."
	ti.Focus()

	m := Model{
		Projects:  projects,
		Selected:  0,
		Focus:     FocusProjects,
		Mode:      ModeNav,
		scanner:   sc,
		cache:     c,
		spinner:   s,
		textInput: ti,
	}
	m.sortProjects()
	return m
}

func deriveProjectView(id string, sessions []scanner.Session, c *cache.Cache) projectView {
	path, inCache := c.Get(id)

	var status ProjectStatus
	if !inCache {
		status = StatusScanning
		path = id
	} else if path == "" {
		status = StatusUnlocated
		path = id
	} else {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			status = StatusOrphaned
		} else {
			status = StatusValid
		}
	}

	return projectView{
		ID:       id,
		Path:     path,
		Status:   status,
		Sessions: sessions,
	}
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
		if p.Status == StatusScanning {
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

	if m.Mode != ModeNav {
		return m.updateInput(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "h", "left", "H":
			m.Focus = FocusProjects
		case "l", "right", "L":
			m.Focus = FocusSessions
		case "up", "k":
			if m.Focus == FocusProjects {
				if m.Cursor > 0 {
					m.Cursor--
					m.Selected = m.Cursor
					m.SessionCursor = 0
				}
			} else {
				if m.SessionCursor > 0 {
					m.SessionCursor--
				}
			}
		case "down", "j":
			if m.Focus == FocusProjects {
				if m.Cursor < len(m.Projects)-1 {
					m.Cursor++
					m.Selected = m.Cursor
					m.SessionCursor = 0
				}
			} else {
				p := m.Projects[m.Selected]
				if m.SessionCursor < len(p.Sessions)-1 {
					m.SessionCursor++
				}
			}
		case "enter", " ":
			if m.Focus == FocusProjects {
				m.Selected = m.Cursor
				m.SessionCursor = 0
			}
		case "c":
			if m.Focus == FocusProjects && len(m.Projects) > 0 {
				p := m.Projects[m.Cursor]
				m.Mode = ModeInputPath
				initialValue := p.Path
				if p.Status == StatusUnlocated || p.Status == StatusScanning {
					home, _ := os.UserHomeDir()
					initialValue = home
				}
				m.textInput.SetValue(initialValue)
				m.prompt = fmt.Sprintf("Enter new directory for project [%s]:", p.ID[:8])
				return m, textinput.Blink
			}
		case "d", "x":
			if m.Focus == FocusProjects && len(m.Projects) > 0 {
				p := m.Projects[m.Cursor]
				totalMessages := 0
				for _, s := range p.Sessions {
					totalMessages += s.MessageCount
				}
				m.Mode = ModeConfirmDelete
				m.prompt = fmt.Sprintf("Permanently delete project [%s] and its %d sessions (%d messages)?",
					p.ID[:8], len(p.Sessions), totalMessages)
			}
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
				m.Projects[i].Status = StatusValid

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
			if m.Projects[i].Status == StatusScanning {
				m.Projects[i].Status = StatusUnlocated
				m.cache.Set(m.Projects[i].ID, "")
				_ = m.cache.Save()
			}
		}
	}

	return m, nil
}

func (m Model) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Mode = ModeNav
			m.Err = nil
			return m, nil
		case "enter":
			if m.Mode == ModeInputPath {
				newPath := m.textInput.Value()
				oldID := m.Projects[m.Cursor].ID

				// 1. Perform deep move
				newID, err := gemini.MoveProject(m.scanner.RootDir, oldID, newPath)
				if err != nil {
					m.Err = err
					m.prompt = fmt.Sprintf("Error migrating project: %v\n\nPress any key to continue...", err)
					return m, nil
				}

				// 2. Update Cache
				_ = m.cache.Delete(oldID)
				m.cache.Set(newID, newPath)
				_ = m.cache.Save()

				// 3. Refresh project in view
				m.Projects[m.Cursor] = deriveProjectView(newID, m.Projects[m.Cursor].Sessions, m.cache)
				m.sortProjects()
				m.Mode = ModeNav
				m.Err = nil
			}
			return m, nil
		case "y", "Y":
			if m.Mode == ModeConfirmDelete {
				id := m.Projects[m.Cursor].ID
				// 1. Delete from physical storage
				if err := gemini.DeleteProject(m.scanner.RootDir, id); err != nil {
					m.Err = err
					m.prompt = fmt.Sprintf("Error deleting project: %v\n\nPress any key to continue...", err)
					return m, nil
				}
				// 2. Delete from cache
				_ = m.cache.Delete(id)

				// 3. Remove from view
				m.Projects = append(m.Projects[:m.Cursor], m.Projects[m.Cursor+1:]...)
				if m.Cursor >= len(m.Projects) && len(m.Projects) > 0 {
					m.Cursor = len(m.Projects) - 1
				}
				if len(m.Projects) > 0 {
					m.Selected = m.Cursor
				} else {
					m.Selected = 0
				}
				m.Mode = ModeNav
			}
		case "n", "N":
			if m.Mode == ModeConfirmDelete {
				m.Mode = ModeNav
			}
		default:
			// If we had an error prompt, any key returns to nav or clears error
			if m.Err != nil {
				m.Mode = ModeNav
				m.Err = nil
				return m, nil
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) isScanningGlobal() bool {
	for _, p := range m.Projects {
		if p.Status == StatusScanning {
			return true
		}
	}
	return false
}

func (m Model) View() string {
	if len(m.Projects) == 0 {
		return "No projects found in ~/.gemini/tmp"
	}

	// Calculate widths: 50/50 split
	sidebarWidth := (m.Width / 2) - 2
	mainWidth := m.Width - sidebarWidth - 6
	paneHeight := m.Height - 6

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
		cursor := renderCursor(m.Focus == FocusProjects, m.Cursor == i)
		style := getRowStyle(m.Selected == i)
		idStr := renderHash(p.ID) + " "
		pathStr := collapseHome(p.Path)

		availableWidth := sidebarWidth - 6 - lipgloss.Width(idStr)
		if p.Status == StatusScanning {
			availableWidth -= 2
		} else if p.Status == StatusOrphaned {
			availableWidth -= 11
		} else if p.Status == StatusUnlocated {
			availableWidth -= 12
		}

		pathStr = truncateMiddle(pathStr, availableWidth)

		var row string
		switch p.Status {
		case StatusScanning:
			row = fmt.Sprintf("%s%s %s", idStr, style.Render(pathStr), m.spinner.View())
		case StatusValid:
			row = fmt.Sprintf("%s%s", idStr, style.Render(pathStr))
		case StatusOrphaned:
			row = fmt.Sprintf("%s%s %s", idStr, style.Inherit(strikethroughStyle).Render(pathStr), style.Render("[Orphaned]"))
		case StatusUnlocated:
			row = fmt.Sprintf("%s%s", idStr, style.Render("[Unlocated]"))
		}

		sidebar.WriteString(fmt.Sprintf("%s%s\n", cursor, row))
	}

	var main strings.Builder
	if m.Selected < len(m.Projects) {
		p := m.Projects[m.Selected]
		displayPath := collapseHome(p.Path)
		if p.Status == StatusUnlocated {
			displayID := p.ID
			if len(displayID) > 12 {
				displayID = displayID[:12] + "..."
			}
			displayPath = displayID
		}
		main.WriteString(titleStyle.Render(fmt.Sprintf("Sessions for %s", displayPath)) + "\n\n")

		if len(p.Sessions) == 0 {
			main.WriteString("No sessions found.")
		} else {
			for i, s := range p.Sessions {
				cursor := renderCursor(m.Focus == FocusSessions, m.SessionCursor == i)
				style := getRowStyle(m.Focus == FocusSessions && m.SessionCursor == i)

				idStr := renderHash(s.ID)
				content := fmt.Sprintf("%s %s | %s",
					idStr,
					style.Render(fmt.Sprintf("%d messages", s.MessageCount)),
					style.Render(formatRelativeTime(s.LastUpdate)))

				main.WriteString(fmt.Sprintf("%s%s\n", cursor, content))
			}
		}
	}

	view := lipgloss.JoinHorizontal(lipgloss.Top,
		listStyle.Width(sidebarWidth).Height(paneHeight).Render(sidebar.String()),
		detailsStyle.Width(mainWidth).Height(paneHeight).Render(main.String()),
	)

	if m.Mode != ModeNav {
		var modalContent string
		if m.Mode == ModeInputPath {
			modalContent = fmt.Sprintf("%s\n\n%s", m.prompt, m.textInput.View())
		} else {
			modalContent = fmt.Sprintf("%s\n\n(y/n)", m.prompt)
		}
		modal := promptStyle.Render(modalContent)

		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, modal)
	}

	return view
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
