package main

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
)

const NON_CONTENT_HEIGHT = 7

// Based on https://github.com/charmbracelet/bubbletea/blob/master/tutorials/basics/README.md

type model struct {
	options  []*Branch
	cursor   int
	offset   int
	height   int
	selected map[int]struct{}
}

func initialModel(options []*Branch) model {
	return model{
		options:  options,
		height:   math.MaxInt,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.resetNavigation() // reset cursor and offset on resize

	case tea.KeyMsg:
		switch msg.String() {

		case "d":
			return m, tea.Quit

		case "ctrl+c", "q":
			m.selected = make(map[int]struct{}) // clear selection
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				if m.cursor == m.offset {
					m.offset--
				}
			}
			if m.cursor > 1 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
				if (m.cursor - m.offset) >= m.maxListSize() {
					m.offset++
				}
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "What branches do you want to delete?\n\n"

	hiddenBefore := "\n"
	if m.offset > 0 {
		hiddenBefore = fmt.Sprintf("(↑ %d more item(s) hidden)\n", m.offset)
	}
	s += hiddenBefore

	upperBound := int(math.Min(float64(len(m.options)), float64(m.offset+m.maxListSize())))

	for i, choice := range m.options[m.offset:upperBound] {
		index := i + m.offset

		cursor := " "
		if m.cursor == index {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[index]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.print())
	}

	hiddenAfter := "\n"
	if upperBound < len(m.options) {
		hiddenAfter = fmt.Sprintf("(↓ %d more item(s) hidden)\n", len(m.options)-upperBound)
	}
	s += hiddenAfter

	s += "\n↑/k move up; ↓/j mode down; enter/space select; d delete; q quit.\n"

	return s
}

func (m model) maxListSize() int {
	return m.height - NON_CONTENT_HEIGHT
}

func (m *model) resetNavigation() {
	m.offset = 0
	m.cursor = 1
}
