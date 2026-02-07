package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"geminictl/internal/cache"
	"geminictl/internal/gemini"
	"geminictl/internal/scanner"
	"sort"

	"github.com/charmbracelet/bubbles/spinner"
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
	ModeMove
	ModeDelete
	ModeInspect
	ModeOpen
	ModeDeleteSession
	ModeMoveSession
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

	strikethroughStyle = lipgloss.NewStyle().
				Strikethrough(true)
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

	scanner *scanner.Scanner
	cache   *cache.Cache
	spinner spinner.Model
	modal   Modal
}

// Internal message to carry the channel along with the result
type resolutionPacket struct {
	res scanner.Resolution
	ch  <-chan scanner.Resolution
}

type ScanFinishedMsg struct{}

type SessionOpenedMsg struct {
	Err error
}

func NewModel(scanned []scanner.ProjectData, c *cache.Cache, sc *scanner.Scanner) *Model {
	var projects []projectView
	for _, p := range scanned {
		projects = append(projects, deriveProjectView(p.ID, p.Sessions, c))
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(highlight)

	m := &Model{
		Projects: projects,
		Selected: 0,
		Focus:    FocusProjects,
		Mode:     ModeNav,
		scanner:  sc,
		cache:    c,
		spinner:  s,
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

// syncState updates the projects list from fresh scanner data while preserving selection.
func (m *Model) syncState(scanned []scanner.ProjectData) {
	// 1. Capture current selection by ID
	var selectedProjectID string
	var selectedSessionID string
	
	if len(m.Projects) > 0 && m.Selected < len(m.Projects) {
		p := m.Projects[m.Selected]
		selectedProjectID = p.ID
		if len(p.Sessions) > 0 && m.SessionCursor < len(p.Sessions) {
			selectedSessionID = p.Sessions[m.SessionCursor].ID
		}
	}

	// 2. Build new state
	var projects []projectView
	for _, p := range scanned {
		projects = append(projects, deriveProjectView(p.ID, p.Sessions, m.cache))
	}
	m.Projects = projects
	m.sortProjects()

	// 3. Restore Project selection
	if selectedProjectID != "" {
		for i, p := range m.Projects {
			if p.ID == selectedProjectID {
				m.Selected = i
				m.Cursor = i // Keep cursor synced with selected
				
				// 4. Restore Session selection within the project
				if selectedSessionID != "" {
					for j, s := range p.Sessions {
						if s.ID == selectedSessionID {
							m.SessionCursor = j
							break
						}
					}
				}
				break
			}
		}
	}

	// 5. Final safety clamps
	if len(m.Projects) == 0 {
		m.Selected = 0
		m.Cursor = 0
		m.SessionCursor = 0
		return
	}
	
	m.Selected = min(m.Selected, len(m.Projects)-1)
	m.Cursor = min(m.Cursor, len(m.Projects)-1)
	
	p := m.Projects[m.Selected]
	if len(p.Sessions) > 0 {
		m.SessionCursor = min(m.SessionCursor, len(p.Sessions)-1)
	} else {
		m.SessionCursor = 0
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.startScanningCmd())
}

func (m *Model) startScanningCmd() tea.Cmd {
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

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// 1. Handle Modal Result
	if res, ok := msg.(ModalResult); ok {
		return m.handleModalResult(res)
	}

	// 2. Global Messages (Window Size, etc.)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// If modal is active, it needs the resize too
		if m.modal != nil {
			m.modal, cmd = m.modal.Update(msg)
			return m, cmd
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		if m.modal == nil {
			return m, cmd
		}
	}

	// 3. Handle Active Modal Update
	if m.modal != nil {
		m.modal, cmd = m.modal.Update(msg)
		return m, cmd
	}

	// 4. Main Navigation
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
		case "enter":
			if m.Focus == FocusProjects {
				m.Selected = m.Cursor
				m.SessionCursor = 0
			} else {
				// Open Session in Gemini CLI
				p := m.Projects[m.Selected]
				if len(p.Sessions) > 0 {
					if p.Status != StatusValid {
						m.modal = ErrorModal{
							Title: "Cannot Open Session",
							Err:   fmt.Errorf("project directory is unlocated or scanning. Please 'Move' ([m]) the project to a valid directory first."),
						}
						return m, m.modal.Init()
					}
					s := p.Sessions[m.SessionCursor]
					m.modal = ConfirmModal{
						Title:  "Open Session",
						Prompt: fmt.Sprintf("Open session [%s] in Gemini CLI?", s.ID[:8]),
					}
					m.Mode = ModeOpen
					return m, m.modal.Init()
				}
			}
		case " ":
			if m.Focus == FocusProjects {
				m.Selected = m.Cursor
				m.SessionCursor = 0
			} else {
				// Inspect Session
				p := m.Projects[m.Selected]
				if len(p.Sessions) > 0 {
					sID := p.Sessions[m.SessionCursor].ID
					fullSession, err := gemini.GetSession(m.scanner.RootDir, p.ID, sID)
					if err != nil {
						m.Err = err
					} else {
						m.Mode = ModeInspect
						m.modal = NewInspectModal(fullSession)
						// We need to trigger a resize for the modal to initialize viewport
						return m, func() tea.Msg { 
							return tea.WindowSizeMsg{Width: m.Width, Height: m.Height} 
						}
					}
				}
			}
		case "m":
			if len(m.Projects) == 0 {
				break
			}
			p := m.Projects[m.Selected]
			if m.Focus == FocusProjects {
				m.Mode = ModeMove
				startDir := p.Path
				if p.Status == StatusUnlocated || p.Status == StatusScanning {
					startDir, _ = os.UserHomeDir()
				}
				m.modal = NewTextInputModal(fmt.Sprintf("Move [%s] to:", p.ID[:8]), startDir, "Absolute path...")
				return m, m.modal.Init()
			} else {
				// Move Session
				if len(p.Sessions) == 0 {
					break
				}
				s := p.Sessions[m.SessionCursor]
				var options []ListOption
				for _, other := range m.Projects {
					if other.ID == p.ID {
						continue
					}
					label := other.Path
					if other.Status == StatusUnlocated {
						label = fmt.Sprintf("[%s] Unlocated", other.ID[:8])
					}
					options = append(options, ListOption{ID: other.ID, Label: label})
				}
				if len(options) == 0 {
					m.Err = fmt.Errorf("no other projects to move session to")
					break
				}
				m.modal = ListSelectorModal{
					Title:   fmt.Sprintf("Move Session [%s] to:", s.ID[:8]),
					Options: options,
				}
				m.Mode = ModeMoveSession
				return m, m.modal.Init()
			}
		case "d":
			if len(m.Projects) == 0 {
				break
			}
			p := m.Projects[m.Selected]
			if m.Focus == FocusProjects {
				totalMessages := 0
				for _, s := range p.Sessions {
					totalMessages += s.MessageCount
				}
				m.modal = ConfirmModal{
					Title:  "Confirm Deletion",
					Prompt: fmt.Sprintf("Permanently delete project [%s] and its %d sessions (%d messages)?", p.ID[:8], len(p.Sessions), totalMessages),
				}
				m.Mode = ModeDelete
				return m, m.modal.Init()
			} else {
				// Delete Session
				if len(p.Sessions) == 0 {
					break
				}
				s := p.Sessions[m.SessionCursor]
				m.modal = ConfirmModal{
					Title:  "Delete Session",
					Prompt: fmt.Sprintf("Permanently delete session [%s] (%d messages)?", s.ID[:8], s.MessageCount),
				}
				m.Mode = ModeDeleteSession
				return m, m.modal.Init()
			}
		}

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
	case SessionOpenedMsg:
		if msg.Err != nil {
			m.modal = ErrorModal{
				Title: "Open Failed",
				Err:   msg.Err,
			}
		}
		// Always refresh state after return, as new messages might have been added
		if updated, err := m.scanner.Scan(); err == nil {
			m.syncState(updated)
		}
	}

	return m, nil
}

func (m *Model) handleModalResult(res ModalResult) (tea.Model, tea.Cmd) {
	m.modal = nil

	if res.Canceled {
		m.Mode = ModeNav
		return m, nil
	}

	switch m.Mode {
	case ModeDelete:
		if res.Value.(bool) {
			p := m.Projects[m.Selected]
			if err := gemini.DeleteProject(m.scanner.RootDir, p.ID); err != nil {
				m.Err = err
			} else {
				_ = m.cache.Delete(p.ID)
				m.Projects = append(m.Projects[:m.Selected], m.Projects[m.Selected+1:]...)
				if m.Selected >= len(m.Projects) && len(m.Projects) > 0 {
					m.Selected = len(m.Projects) - 1
				}
				m.Cursor = m.Selected
			}
		}
	case ModeMove:
		newPath := res.Value.(string)
		oldID := m.Projects[m.Selected].ID
		newID, err := gemini.MoveProject(m.scanner.RootDir, oldID, newPath)
		if err != nil {
			m.Err = err
		} else {
			_ = m.cache.Delete(oldID)
			m.cache.Set(newID, newPath)
			_ = m.cache.Save()
			m.Projects[m.Selected] = deriveProjectView(newID, m.Projects[m.Selected].Sessions, m.cache)
			m.sortProjects()
		}
	case ModeOpen:
		if res.Value.(bool) {
			p := m.Projects[m.Selected]
			s := p.Sessions[m.SessionCursor]
			
			// Wrap in a shell to clear the screen and show a loading message
			script := fmt.Sprintf("clear && echo 'Launching Gemini CLI for session [%s]...' && gemini --resume %s && clear", s.ID[:8], s.ID)
			c := exec.Command("sh", "-c", script)
			c.Dir = p.Path
			
			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				return SessionOpenedMsg{Err: err}
			})
		}
	case ModeDeleteSession:
		if res.Value.(bool) {
			p := &m.Projects[m.Selected]
			s := p.Sessions[m.SessionCursor]
			if err := gemini.DeleteSession(m.scanner.RootDir, p.ID, s.ID); err != nil {
				m.Err = err
			} else {
				// Refresh the entire state to be safe and simple
				if updated, err := m.scanner.Scan(); err == nil {
					m.syncState(updated)
				}
			}
		}
	case ModeMoveSession:
		targetProjectID := res.Value.(string)
		p := &m.Projects[m.Selected]
		s := p.Sessions[m.SessionCursor]
		if err := gemini.MoveSession(m.scanner.RootDir, p.ID, targetProjectID, s.ID); err != nil {
			m.Err = err
		} else {
			// Refresh the entire state
			if updated, err := m.scanner.Scan(); err == nil {
				m.syncState(updated)
			}
		}
	}

	m.Mode = ModeNav
	return m, nil
}

func (m *Model) isScanningGlobal() bool {
	for _, p := range m.Projects {
		if p.Status == StatusScanning {
			return true
		}
	}
	return false
}

func (m *Model) View() string {
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

	if m.modal != nil {
		return m.modal.View(m.Width, m.Height)
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
