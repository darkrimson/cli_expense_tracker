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

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateID := updateCmd.Int("id", 0, "Expense ID")
	updateDescription := updateCmd.String("description", "", "Expense description")
	updateAmount := updateCmd.Float64("amount", 0, "Expense amount")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID := deleteCmd.Int("id", 0, "Expense ID")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	summaryMonth := summaryCmd.Int("month", 0, "Month number 1-12")

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		expense := model.Expense{
			Description: *addDescription,
			Amount:      *addAmount,
			Date:        time.Now(),
		}
		savedExp, err := expenseService.AddExpense(expense.Amount, expense.Description)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("Expense added successfully (ID: %d)\n", savedExp.ID)

	case "update":
		updateCmd.Parse(os.Args[2:])
		err := expenseService.UpdateExpense(*updateID, *updateDescription, *updateAmount)
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
		fmt.Println("ID\tDate\t\tDescription\tAmount")

		for _, e := range expenses {
			fmt.Printf("%d\t%s\t%s\t$%.2f\n", e.ID, e.Date.Format("2006-01-02"), e.Description, e.Amount)
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
	}
}
