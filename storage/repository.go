package storage

import (
	"cli_expense_tracker/model"
)

type ExpenseRepository interface {
	Save(exp model.Expense) (model.Expense, error)
	Update(id int, exp model.Expense) error
	Delete(id int) error
	GetAll() ([]model.Expense, error)
	GetByID(id int) (model.Expense, error)
}
