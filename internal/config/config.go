package config

import (
	"encoding/json"
	"errors"
	"expense-tracker/internal/expense"
	"os"
)

var FILE string = "./save.json"

func GenerateExpenseTracker() (*expense.ExpenseTracker, error) {
	expenseTracker, err := loadExpenseTracker()
	if err != nil {
		return nil, err
	}
	return expenseTracker, nil
}

func loadExpenseTracker() (*expense.ExpenseTracker, error) {
	data, err := os.ReadFile(FILE)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || !fileExists(FILE) {
			return expense.NewExpenseTracker(1000)
		}
		return nil, err
	}

	var tracker expense.ExpenseTracker
	if err := json.Unmarshal(data, &tracker); err != nil {
		return nil, err
	}

	var maxID int64 = 0
	for id := range tracker.Expenses {
		if id > maxID {
			maxID = id
		}
	}
	if tracker.NextID < maxID {
		tracker.NextID = maxID + 1
	}

	return &tracker, nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || !os.IsNotExist(err)
}

func SaveState(tracker *expense.ExpenseTracker) error {
	data, err := json.MarshalIndent(tracker, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(FILE, data, 0644)
	return err
}
