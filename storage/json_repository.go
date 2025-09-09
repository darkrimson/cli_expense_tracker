package storage

import (
	"cli_expense_tracker/model"
	"encoding/json"
	"errors"
	"os"
)

type JSONRepository struct {
	filePath string
}

func NewJSONRepository(filePath string) *JSONRepository {
	return &JSONRepository{filePath: filePath}
}

func (r *JSONRepository) load() ([]model.Expense, error) {
	var expenses []model.Expense

	file, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Expense{}, nil // если файла нет → пустой список
		}
		return nil, err
	}

	if len(file) > 0 {
		err = json.Unmarshal(file, &expenses)
		if err != nil {
			return nil, err
		}
	}
	return expenses, nil
}

func (r *JSONRepository) saveAll(expenses []model.Expense) error {
	data, err := json.MarshalIndent(expenses, " ", "")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *JSONRepository) Save(exp model.Expense) (model.Expense, error) {
	expenses, err := r.load()
	if err != nil {
		return model.Expense{}, err
	}

	exp.ID = len(expenses) + 1
	expenses = append(expenses, exp)

	if err := r.saveAll(expenses); err != nil {
		return model.Expense{}, nil
	}
	return exp, nil
}

func (r *JSONRepository) Update(id int, updated model.Expense) error {
	expenses, err := r.load()
	if err != nil {
		return err
	}

	found := false
	for i, e := range expenses {
		if e.ID == id {
			updated.ID = id
			expenses[i] = updated
			found = true
			break
		}
	}

	if !found {
		return errors.New("expense not found")
	}

	return r.saveAll(expenses)
}

func (r *JSONRepository) Delete(id int) error {
	expenses, err := r.load()
	if err != nil {
		return err
	}

	newExpenses := []model.Expense{}
	found := false

	for _, e := range expenses {
		if e.ID == id {
			found = true
			continue
		}
		newExpenses = append(newExpenses, e)
	}

	if !found {
		return errors.New("expense not found")
	}

	return r.saveAll(newExpenses)
}

func (r *JSONRepository) GetAll() ([]model.Expense, error) {
	return r.load()
}

func (r *JSONRepository) GetByID(id int) (model.Expense, error) {
	expenses, err := r.load()
	if err != nil {
		return model.Expense{}, err
	}

	for _, e := range expenses {
		if e.ID == id {
			return e, nil
		}
	}

	return model.Expense{}, errors.New("expense not found")
}
