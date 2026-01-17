package main

import (
	"fmt"
	"log"
	"os"
	"task-tracker/internal/storage"
	"task-tracker/internal/tasks"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#001fff")).
		Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return m.NewStatusMessage(statusMessageStyle.Render("You chose " + title))

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(statusMessageStyle.Render("Deleted " + title))
			}
		}

		return nil
	}

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{keys.choose, keys.remove}
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{{keys.choose, keys.remove}}
	}

	return d
}

type listKeyMap struct {
	insertItem key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
	}
}

type model struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
}

func newModel(service *tasks.Service) model {
	delegateKeys := newDelegateKeyMap()
	listKeys := newListKeyMap()
	var items []list.Item

	listTasks, err := service.ListTasks()
	if err != nil {
		return model{}
	}

	if len(listTasks) == 0 {

	} else {
		for _, task := range listTasks {
			stat := "Status: " + task.Status + "\n Created at: " + task.CreatedAt.Format("02.01.2006 15:04") + "\n Updated at: " + task.UpdatedAt.Format("02.01.2006 15:04")
			items = append(items, item{title: task.Description, description: stat})
		}
	}

	delegate := newItemDelegate(delegateKeys)
	groceryList := list.New(items, delegate, 0, 0)
	groceryList.Title = "Мои задачи"
	groceryList.Styles.Title = titleStyle
	groceryList.SetShowHelp(true)

	groceryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{listKeys.insertItem}
	}

	return model{
		list:         groceryList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		if key.Matches(msg, m.keys.insertItem) {
			return m, m.list.NewStatusMessage(statusMessageStyle.Render("Нажми a для добавления (заглушка)"))
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

func main() {
	jsonStorage, err := storage.NewJSONStorage("data.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	p := tea.NewProgram(newModel(tasks.NewService(jsonStorage)), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
