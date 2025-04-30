package main

import (
    "fmt"
	"math/rand"
	"time"
    "os"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
)

type Food struct{
	name string
	calorieValue int
	key string
}

type Model struct {
    catWeight int
    fedCount  int
    question  string
    err       error
	catName   string
	foods []Food
	selectedFood int
	exercising  bool   
    exerciseKey string
    exerciseCount int 
	exerciseStep int 
    starved bool
}

func initialModel() Model {
	rand.Seed(time.Now().UnixNano())

	initialWeight := rand.Intn(5) + 1 
	catNames := []string{"Whiskers", "Mittens", "Fluffy", "Shadow", "Ginger", "Simba", "Luna", "Oliver", "Bella", "Charlie"}
	catName := catNames[rand.Intn(len(catNames))]

	foods := []Food{
		{name: "Fish", calorieValue: 2, key: "a"},
		{name: "Chicken", calorieValue: 3, key: "s"},
		{name: "Beef", calorieValue: 4, key: "d"},
		{name: "Keso extra protein", calorieValue: -2, key: "f"},
		{name: "Treats", calorieValue: 2, key: "g"},
		{name: "Catnip", calorieValue: 1, key: "h"},
	}

    return Model{
        catWeight: initialWeight, 
        fedCount: 0,
        question: "Do you want to feed " + catName + "? (press key)",
		catName: catName,
		foods: foods,
		exercising:  false,
        exerciseKey: "jk",
        exerciseCount: 0,
		exerciseStep: 0,
        starved: false,
    }
}

func (m Model) Init() tea.Cmd {
    return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case msg.String() == "ctrl+c" || msg.String() == "q":
            return m, tea.Quit
	case m.exercising:
		expectedKey := ""
		if m.exerciseStep == 0 {
			expectedKey = "j"
		} else {
			expectedKey = "k"
		}

		if msg.String() == expectedKey {
			m.exerciseCount++
			m.exerciseStep = 1 - m.exerciseStep

			if m.exerciseCount >= 10 { 
				m.catWeight -= 5 
				if m.catWeight < 5 {
					m.catWeight = 5 
				}
				m.exercising = false
				m.exerciseCount = 0
				m.exerciseStep = 0 
				m.question = m.catName + " feels much better! What do you want to feed " + m.catName + "? (Press the key)"
				m.err = nil
			} else {
				m.question = fmt.Sprintf("%s is exercising! Press %s (%d/%d)", m.catName, expectedKey, m.exerciseCount, 10)
			}
			return m, nil
		} else {
			m.exercising = false
			m.exerciseCount = 0
			m.exerciseStep = 0 
			m.question = "Incorrect exercise key! " + m.catName + " stopped exercising. What do you want to feed " + m.catName + "? (Press the key)"
			return m, nil
		}

    case m.starved:
        m.catWeight <= 0 {
            m.catWeight = 0
            m.starved = true;
        }
        m.question = "Your cat has starved to death! Press 'q' to quit."
        return m, nil

        case msg.String() == "e" && m.catWeight >= 20:
            m.exercising = true
            m.exerciseCount = 0
            m.question = fmt.Sprintf("%s is exercising! Press %s (%d/%d)", m.catName, m.exerciseKey, m.exerciseCount, 10)
            m.err = nil
            return m, nil

        default:  
            for _, food := range m.foods {
                if msg.String() == food.key {
                    if m.catWeight >= 20 { 
                        m.question = m.catName + " is too fat to feed! Press 'e' to exercise."
                        return m, nil
                    }
                    m.catWeight += food.calorieValue
                    m.fedCount++
                    m.question = "What do you want to feed " + m.catName + " again? (Press the key)"
                    m.err = nil
                    return m, nil
                }
            }

            m.question = "Invalid key. What do you want to feed " + m.catName + "? Press the key)"
        }
    }
    return m, nil
}

func (m Model) View() string {
    var s strings.Builder

	s.WriteString(fmt.Sprintf("Cat Name: %s\n", m.catName))
	s.WriteString(fmt.Sprintf("Cat Weight: %d\n", m.catWeight))
    s.WriteString(fmt.Sprintf("Times Fed: %d\n\n", m.fedCount))
    if m.err != nil {
        s.WriteString(fmt.Sprintf("Error: %s\n\n", m.err))
    }

	s.WriteString(getCatArt(m.catWeight, m.exerciseStep, m.exercising))
    s.WriteString("\n\n")

    s.WriteString("Food Options:\n")
    for _, food := range m.foods {
        s.WriteString(fmt.Sprintf("%s: %s (%d kcal)\n", food.key, food.name, food.calorieValue))
    }

    s.WriteString("\n" + m.question + "\n")

    return s.String()
}

func getCatArt(weight int, exerciseStep int, exercising bool) string {
	if exercising {
        if exerciseStep == 0 {
           
            return `     
       / \__/ \
      (  0.0  )
        >===<
     <</_____\

        `
        } else {
            
            return `
       / \__/ \
      (  0.0  )
        >===<
       /_____\>>

        `
        }
    }

    switch {
    case weight < 8:
        return `
	/\_/\ 
	=^.^=


    `
    case weight < 12:
        return `
	/\_/\  
       ( o.o )
	> ^ <

    `
    case weight < 16:
        return `
        /\_/\
       ( O.O )
        >   <
        /   \
    `
    case weight < 20:
        return `
        / \_/\
       ( >.<  )
        >---<
       /     \
    `
    default:
        return `     
       / \__/ \
      (  U.U  )
        >===<
       /_____\
        `
    }
}

func main() {
    p := tea.NewProgram(initialModel(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}