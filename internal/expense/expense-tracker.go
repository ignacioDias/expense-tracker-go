package expense

import (
	"errors"
	"strings"
	"time"
)

type ExpenseTracker struct {
	NextID            int64             `json:"nextID"`
	Expenses          map[int64]Expense `json:"expenses"`
	MaxBudgetPerMonth float64           `json:"maxBudgetPerMonth"`
	CurrentBudget     float64           `json:"currentBudget"`
	CurrentMonth      time.Month        `json:"currentMonth"`
	CurrentYear       int               `json:"currentYear"`
}

type Expense struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Category    string    `json:"category"`
}

func NewExpense(id int64, descr string, amount float64, category *string) (*Expense, error) {
	if amount <= 0 {
		return nil, errors.New("Invalid amount")
	}
	if descr == "" {
		return nil, errors.New("Invalid description")
	}
	newCat := "others"
	if category != nil {
		newCat = strings.ToLower(*category)
	}

	return &Expense{
		ID:          id,
		Description: descr,
		Amount:      amount,
		Category:    newCat,
		Date:        time.Now(),
	}, nil
}

func NewExpenseTracker(maxBudgetPerMonth float64) (*ExpenseTracker, error) {
	if maxBudgetPerMonth <= 0 {
		return nil, errors.New("Invalid maxBudget")
	}
	return &ExpenseTracker{
		NextID:            1,
		Expenses:          make(map[int64]Expense),
		MaxBudgetPerMonth: maxBudgetPerMonth,
		CurrentBudget:     maxBudgetPerMonth,
		CurrentMonth:      time.Now().Month(),
		CurrentYear:       time.Now().Year(),
	}, nil
}

func (tracker *ExpenseTracker) AddExpense(description string, category *string, amount float64) error {
	expense, err := NewExpense(tracker.NextID, description, amount, category)
	if err != nil {
		return err
	}
	if expense.Date.Month() != tracker.CurrentMonth {
		tracker.CurrentMonth = expense.Date.Month()
		tracker.CurrentYear = expense.Date.Year()
		tracker.CurrentBudget = tracker.MaxBudgetPerMonth
	}
	tracker.CurrentBudget -= expense.Amount
	tracker.Expenses[tracker.NextID] = *expense
	tracker.NextID = tracker.NextID + 1
	return nil
}

func (tracker *ExpenseTracker) Summary() float64 {
	total := 0.0
	for _, expense := range tracker.Expenses {
		total += expense.Amount
	}
	return total
}
func (tracker *ExpenseTracker) SummaryPerMonth(month time.Month) float64 {
	total := 0.0
	for _, expense := range tracker.Expenses {
		if expense.Date.Month() == month && expense.Date.Year() == tracker.CurrentYear {
			total += expense.Amount
		}
	}
	return total
}

func (tracker *ExpenseTracker) UpdateDescription(id int64, description string) error {
	if _, ok := tracker.Expenses[id]; !ok {
		return errors.New("Invalid ID")
	}
	expense := tracker.Expenses[id]
	expense.Description = description
	tracker.Expenses[id] = expense
	return nil
}

func (tracker *ExpenseTracker) UpdateCategory(id int64, category string) error {
	if _, ok := tracker.Expenses[id]; !ok {
		return errors.New("Invalid ID")
	}
	expense := tracker.Expenses[id]
	if category == "" {
		expense.Category = "others"
	} else {
		expense.Category = strings.ToLower(category)
	}
	tracker.Expenses[id] = expense
	return nil
}

func (tracker *ExpenseTracker) UpdateMaxBudgetPerMonth(maxBudgetPerMonth float64) error {
	if maxBudgetPerMonth <= 0 {
		return errors.New("Invalid budget")
	}
	tracker.MaxBudgetPerMonth = maxBudgetPerMonth
	tracker.CurrentBudget = maxBudgetPerMonth
	return nil
}

func (tracker *ExpenseTracker) UpdateAmount(id int64, amount float64) error {
	if _, ok := tracker.Expenses[id]; !ok {
		return errors.New("Invalid ID")
	}
	if amount <= 0 {
		return errors.New("Invalid Amount")
	}
	expense := tracker.Expenses[id]
	oldAmount := expense.Amount
	expense.Amount = amount
	tracker.Expenses[id] = expense

	if expense.Date.Month() == tracker.CurrentMonth && expense.Date.Year() == tracker.CurrentYear {
		tracker.CurrentBudget += oldAmount - amount
	}
	return nil
}

func (tracker *ExpenseTracker) DeleteExpense(id int64) error {
	if _, ok := tracker.Expenses[id]; !ok {
		return errors.New("Invalid ID")
	}
	expense := tracker.Expenses[id]
	if expense.Date.Month() == tracker.CurrentMonth && expense.Date.Year() == tracker.CurrentYear {
		tracker.CurrentBudget += expense.Amount
	}
	delete(tracker.Expenses, id)
	return nil
}
func (tracker *ExpenseTracker) GetExpensesPerCategory(category string) []Expense {
	var expenses []Expense
	cat := strings.ToLower(category)
	for _, expense := range tracker.Expenses {
		if cat == expense.Category {
			expenses = append(expenses, expense)
		}
	}
	return expenses
}
