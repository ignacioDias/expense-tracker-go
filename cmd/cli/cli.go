package cli

import (
	"errors"
	"expense-tracker/internal/config"
	"expense-tracker/internal/expense"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CLI struct {
	Args           []string
	ExpenseTracker *expense.ExpenseTracker
}

var ErrInvalidCommand error = errors.New("Invalid command, try with help for more info")

func NewCLI(args []string) (*CLI, error) {
	expenseTracker, err := config.GenerateExpenseTracker()
	if err != nil {
		return nil, err
	}
	return &CLI{
		Args:           args,
		ExpenseTracker: expenseTracker,
	}, nil
}

func (c *CLI) Start() error {
	if len(c.Args) == 0 {
		return ErrInvalidCommand
	}
	switch c.Args[0] {
	case "help":
		printHelp()
	case "summary":
		if err := c.handleSummaryCase(); err != nil {
			return err
		}
	case "list":
		if len(c.Args) == 1 {
			c.handleListCase()
		} else {
			if err := c.handleListCategoryCase(); err != nil {
				return err
			}
		}
	case "delete":
		if len(c.Args) != 3 {
			return ErrInvalidCommand
		}
		if err := c.handleDeleteCase(); err != nil {
			return err
		}
	case "add":
		if len(c.Args) != 5 && len(c.Args) != 7 {
			return ErrInvalidCommand
		}
		if err := c.handleAddCase(); err != nil {
			return err
		}

	case "update":
		if len(c.Args) < 4 {
			return ErrInvalidCommand
		}
		if err := c.handleUpdateCase(); err != nil {
			return err
		}
	case "update-budget":
		if len(c.Args) != 3 {
			return ErrInvalidCommand
		}
		if err := c.handleUpdateBudgetCase(); err != nil {
			return err
		}
		fmt.Println("Budget updated successfully")
	default:
		return ErrInvalidCommand
	}
	return nil
}

func printHelp() {
	fmt.Println("Expense Tracker - Manage your expenses")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  expense-tracker <command> [options]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  add --description <desc> --amount <amount> [--category <cat>]")
	fmt.Println("      Add a new expense")
	fmt.Println("")
	fmt.Println("  update --id <id> [--description <desc>] [--amount <amount>] [--category <cat>]")
	fmt.Println("      Update an expense (at least one field required)")
	fmt.Println("")
	fmt.Println("  delete --id <id>")
	fmt.Println("      Delete an expense by ID")
	fmt.Println("")
	fmt.Println("  list [--category <category>]")
	fmt.Println("      List all expenses or filter by category")
	fmt.Println("")
	fmt.Println("  summary [--month <1-12>]")
	fmt.Println("      Show total expenses (all time or by month)")
	fmt.Println("")
	fmt.Println("  update-budget --amount <amount>")
	fmt.Println("      Update monthly budget limit")
	fmt.Println("")
	fmt.Println("  help")
	fmt.Println("      Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  expense-tracker add --description \"Lunch\" --amount 15.50 --category Food")
	fmt.Println("  expense-tracker update --id 1 --amount 20.00")
	fmt.Println("  expense-tracker update --id 2 --description \"Dinner\" --category Food")
	fmt.Println("  expense-tracker list --category Food")
	fmt.Println("  expense-tracker summary --month 2")
	fmt.Println("  expense-tracker update-budget --amount 1500.00")
	fmt.Println("  expense-tracker delete --id 3")
}

func (c *CLI) handleAddCase() error {
	if c.Args[1] != "--description" || c.Args[3] != "--amount" {
		return ErrInvalidCommand
	}
	amount, err := strconv.ParseFloat(c.Args[4], 64)
	if err != nil {
		return ErrInvalidCommand
	}
	if len(c.Args) == 7 {
		if c.Args[5] != "--category" {
			return ErrInvalidCommand
		}
		c.checkBudget(amount)
		return c.ExpenseTracker.AddExpense(c.Args[2], &c.Args[6], amount)
	}
	c.checkBudget(amount)
	return c.ExpenseTracker.AddExpense(c.Args[2], nil, amount)
}

func (c *CLI) checkBudget(spent float64) {
	percentUsed := (spent / c.ExpenseTracker.MaxBudgetPerMonth) * 100
	if c.ExpenseTracker.CurrentBudget < 0 {
		fmt.Printf("\n⚠️  Warning: Over budget by $%.2f!\n", -c.ExpenseTracker.CurrentBudget)
	} else if percentUsed > 80 {
		fmt.Println("\n⚠️  Warning: Running low on budget!")
	}
}

func (c *CLI) handleDeleteCase() error {
	if c.Args[1] != "--id" {
		return ErrInvalidCommand
	}
	id, err := strconv.ParseInt(c.Args[2], 10, 64)
	if err != nil {
		return ErrInvalidCommand
	}
	if err := c.ExpenseTracker.DeleteExpense(id); err != nil {
		return err
	}
	return nil
}

func (c *CLI) handleListCase() {
	if len(c.ExpenseTracker.Expenses) == 0 {
		fmt.Println("No expenses found")
		return
	}
	fmt.Println("ID   Date         Description          Category      Amount")
	fmt.Println("================================================================")

	ids := make([]int64, 0, len(c.ExpenseTracker.Expenses))
	for id := range c.ExpenseTracker.Expenses {
		ids = append(ids, id)
	}
	// Sort IDs
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			if ids[i] > ids[j] {
				ids[i], ids[j] = ids[j], ids[i]
			}
		}
	}

	for _, id := range ids {
		expense := c.ExpenseTracker.Expenses[id]
		fmt.Printf("%-4d %-12s %-20s %-12s $%.2f\n",
			expense.ID,
			expense.Date.Format("2006-01-02"),
			expense.Description,
			expense.Category,
			expense.Amount)
	}
}
func (c *CLI) handleListCategoryCase() error {
	if len(c.Args) != 3 || c.Args[1] != "--category" {
		return ErrInvalidCommand
	}
	fmt.Printf("Category: %s\n", c.Args[2])
	fmt.Println("ID   Date         Description          Amount")
	fmt.Println("===============================================")

	ids := make([]int64, 0)
	for id, expense := range c.ExpenseTracker.Expenses {
		if expense.Category == strings.ToLower(c.Args[2]) {
			ids = append(ids, id)
		}
	}
	// Sort IDs
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			if ids[i] > ids[j] {
				ids[i], ids[j] = ids[j], ids[i]
			}
		}
	}

	for _, id := range ids {
		expense := c.ExpenseTracker.Expenses[id]
		fmt.Printf("%-4d %-12s %-20s $%.2f\n",
			expense.ID,
			expense.Date.Format("2006-01-02"),
			expense.Description,
			expense.Amount)
	}
	return nil
}
func (c *CLI) handleSummaryCase() error {
	if len(c.Args) == 1 {
		c.handleSummary()
		return nil
	} else if len(c.Args) == 3 {
		if err := c.handleSummaryPerMonth(); err != nil {
			return err
		}
	} else {
		return ErrInvalidCommand
	}
	return nil
}

func (c *CLI) handleSummaryPerMonth() error {
	if c.Args[1] != "--month" {
		return ErrInvalidCommand
	}
	monthNumber, err := strconv.ParseInt(c.Args[2], 10, 64)
	if err != nil {
		return ErrInvalidCommand
	}
	if monthNumber < 1 || monthNumber > 12 {
		return ErrInvalidCommand
	}
	month := time.Month(monthNumber)
	total := c.ExpenseTracker.SummaryPerMonth(month)
	fmt.Printf("Total expenses for %s: $%.2f\n", month, total)
	return nil
}

func (c *CLI) handleSummary() {
	total := c.ExpenseTracker.Summary()
	fmt.Printf("Total expenses: $%.2f\n", total)
}

func (c *CLI) handleUpdateCase() error {
	if c.Args[1] != "--id" {
		return ErrInvalidCommand
	}
	id, err := strconv.ParseInt(c.Args[2], 10, 64)
	if err != nil {
		return ErrInvalidCommand
	}

	if _, ok := c.ExpenseTracker.Expenses[id]; !ok {
		return errors.New("Expense ID not found")
	}

	updated := false
	for i := 3; i < len(c.Args); i += 2 {
		if i+1 >= len(c.Args) {
			return ErrInvalidCommand
		}
		flag := c.Args[i]
		value := c.Args[i+1]

		switch flag {
		case "--description":
			if err := c.ExpenseTracker.UpdateDescription(id, value); err != nil {
				return err
			}
			updated = true
		case "--amount":
			amount, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return errors.New("Invalid amount")
			}
			if err := c.ExpenseTracker.UpdateAmount(id, amount); err != nil {
				return err
			}
			updated = true
		case "--category":
			if err := c.ExpenseTracker.UpdateCategory(id, value); err != nil {
				return err
			}
			updated = true
		default:
			return ErrInvalidCommand
		}
	}

	if !updated {
		return errors.New("At least one field to update is required")
	}

	return nil
}

func (c *CLI) handleUpdateBudgetCase() error {
	if c.Args[1] != "--amount" {
		return ErrInvalidCommand
	}
	amount, err := strconv.ParseFloat(c.Args[2], 64)
	if err != nil {
		return errors.New("Invalid amount")
	}
	if err := c.ExpenseTracker.UpdateMaxBudgetPerMonth(amount); err != nil {
		return err
	}
	return nil
}
