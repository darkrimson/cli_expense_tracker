package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"cli_expense_tracker/model"
	"cli_expense_tracker/service"
	"cli_expense_tracker/storage"
)

func printHelp() {
	fmt.Println("Usage: expense-tracker <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  add       Add a new expense")
	fmt.Println("    --description  string   Expense description")
	fmt.Println("    --amount       float    Expense amount")
	fmt.Println("    --category     string   Expense category")
	fmt.Println()
	fmt.Println("  update    Update an existing expense")
	fmt.Println("    --id           int      Expense ID")
	fmt.Println("    --description  string   Expense description")
	fmt.Println("    --amount       float    Expense amount")
	fmt.Println("    --category     string   Expense category")
	fmt.Println()
	fmt.Println("  delete    Delete an expense by ID")
	fmt.Println("    --id           int      Expense ID")
	fmt.Println()
	fmt.Println("  list      List all expenses")
	fmt.Println()
	fmt.Println("  summary   Show total expenses")
	fmt.Println("    --month        int      Month number (1-12)")
	fmt.Println()
	fmt.Println("  help      Show this help message")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'add', 'update', 'delete', 'list' or 'summary' subcommands")
		return
	}

	repo := storage.NewJSONRepository("expenses.json")
	expenseService := service.NewExpenseService(repo)

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDescription := addCmd.String("description", "", "Expense description")
	addAmount := addCmd.Float64("amount", 0, "Expense amount")
	addCategory := addCmd.String("category", "", "Expense category")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateID := updateCmd.Int("id", 0, "Expense ID")
	updateDescription := updateCmd.String("description", "", "Expense description")
	updateAmount := updateCmd.Float64("amount", 0, "Expense amount")
	updateCategory := updateCmd.String("category", "", "Expense category")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID := deleteCmd.Int("id", 0, "Expense ID")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	summaryMonth := summaryCmd.Int("month", 0, "Month number 1-12")

	switch os.Args[1] {

	case "help":
		printHelp()

	case "add":
		addCmd.Parse(os.Args[2:])
		expense := model.Expense{
			Description: *addDescription,
			Category:    *addCategory,
			Amount:      *addAmount,
			Date:        time.Now(),
		}
		savedExp, err := expenseService.AddExpense(expense.Amount, expense.Description, *addCategory)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("Expense added successfully (ID: %d)\n", savedExp.ID)

	case "update":
		updateCmd.Parse(os.Args[2:])
		err := expenseService.UpdateExpense(*updateID, *updateDescription, *updateCategory, *updateAmount)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Expense updated successfully")

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		err := expenseService.DeleteExpense(*deleteID)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Expense deleted successfully")

	case "list":
		listCmd.Parse(os.Args[2:])
		expenses, err := expenseService.ListExpenses()

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("ID\tDate\t\tDescription\tCategory\tAmount")

		for _, e := range expenses {
			fmt.Printf("%d\t%s\t%s\t\t%s\t\t$%.2f\n", e.ID, e.Date.Format("2006-01-02"), e.Description, e.Category, e.Amount)
		}

	case "summary":
		summaryCmd.Parse(os.Args[2:])
		if *summaryMonth != 0 {
			total, err := expenseService.GetSummaryByMonth(*summaryMonth)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Printf("Total expenses for month %d: $%.2f\n", *summaryMonth, total)
		} else {
			total, err := expenseService.GetSummary()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Printf("Total expenses: $%.2f\n", total)
		}

	default:
		fmt.Println("Unknown command. Expected 'add', 'update', 'delete', 'list' or 'summary'")
		printHelp()
	}
}
