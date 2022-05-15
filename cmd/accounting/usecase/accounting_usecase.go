package usecase

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/accounting/repo"
)

type AccountingUsecase struct {
	tasksRepo *repo.AccountingRepo
}

func NewAccountingUsecase(tasksRepo *repo.AccountingRepo) *AccountingUsecase {
	return &AccountingUsecase{
		tasksRepo: tasksRepo,
	}
}

func (tu *AccountingUsecase) AddLog() {
	//tu.tasksRepo.AddLog()
}
