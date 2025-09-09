package service

import (
	"cli_expense_tracker/model"
	"cli_expense_tracker/storage"
	"errors"
	"time"
)

type ExpenseService struct {
	repo storage.ExpenseRepository
}

func NewExpenseService(r storage.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: r}
}

func (s *ExpenseService) AddExpense(amount float64, description, category string) (model.Expense, error) {
	if amount <= 0 {
		return model.Expense{}, errors.New("amount must be greater than zero")
	}

	exp := model.Expense{
		Description: description,
		Category:    category,
		Amount:      amount,
		Date:        time.Now(),
	}

	return s.repo.Save(exp)
}

func (s *ExpenseService) DeleteExpense(id int) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *ExpenseService) GetSummary() (float64, error) {
	expenses, err := s.repo.GetAll()
	if err != nil {
		return 0, err
	}

	var total float64
	for _, e := range expenses {
		total += e.Amount
	}
	return total, nil
}

func (s *ExpenseService) ListExpenses() ([]model.Expense, error) {
	expenses, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (s *ExpenseService) GetSummaryByMonth(month int) (float64, error) {
	expenses, err := s.repo.GetAll()
	if err != nil {
		return 0, err
	}

	var total float64
	for _, e := range expenses {
		if e.Date.Month() == time.Month(month) {
			total += e.Amount
		}
	}
	return total, nil
}
