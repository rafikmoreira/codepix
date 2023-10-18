package model

import (
	"errors"
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	TransactionPeding    string = "pending"
	TransactionCompleted string = "completed"
	TransactionError     string = "error"
	TransactionCancelled string = "cancelled"
	TransactionConfirmed string = "confirmed"
)

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transactions []*Transaction
}

type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"-"`
	Amount            float64  `json:"amount" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	Status            string   `json:"status" valid:"notnull"`
	Description       string   `json:"description" valid:"notnull"`
	CancelDescription string   `json:"cancel_description" valid:"notnull"`
}

func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Confirm() error {
	t.Status = TransactionConfirmed
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Cancel(description string) error {
	t.Status = TransactionCancelled
	t.UpdatedAt = time.Now()
	t.CancelDescription = description
	err := t.isValid()
	return err
}

func (t *Transaction) Error(description string) error {
	t.Status = TransactionError
	t.UpdatedAt = time.Now()
	t.CancelDescription = description
	err := t.isValid()
	return err
}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)

	if t.Amount <= 0 {
		return errors.New("the amount must be greater than 0")
	}

	if t.Status != TransactionPeding &&
		t.Status != TransactionCompleted &&
		t.Status != TransactionError &&
		t.Status != TransactionCancelled &&
		t.Status != TransactionConfirmed {
		return errors.New("invalid status for the transaction")
	}

	if t.PixKeyTo.AccountID == t.AccountFrom.ID {
		return errors.New("the source and destination account cannot be the same")
	}

	if err != nil {
		return err
	}

	return nil
}

func NewTransaction(
	accountFrom *Account,
	amount float64,
	pixKeyTo *PixKey,
	description string,
	id string,
) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount:      amount,
		PixKeyTo:    pixKeyTo,
		Status:      TransactionPeding,
		Description: description,
	}

	if id == "" {
		transaction.ID = uuid.NewV4().String()
	} else {
		transaction.ID = id
	}

	transaction.CreatedAt = time.Now()

	err := transaction.isValid()

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
