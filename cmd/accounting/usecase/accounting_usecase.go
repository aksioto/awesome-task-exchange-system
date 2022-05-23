package usecase

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/accounting/internal/model"
	"github.com/aksioto/awesome-task-exchange-system/cmd/accounting/repo"
	"github.com/aksioto/awesome-task-exchange-system/internal/helper"
	v1 "github.com/aksioto/awesome-task-exchange-system/internal/model/streaming/v1"
	v2 "github.com/aksioto/awesome-task-exchange-system/internal/model/streaming/v2"
	"github.com/google/uuid"
	"log"
)

type AccountingUsecase struct {
	accountingRepo *repo.AccountingRepo
}

func NewAccountingUsecase(accountingRepo *repo.AccountingRepo) *AccountingUsecase {
	return &AccountingUsecase{
		accountingRepo: accountingRepo,
	}
}

func (u *AccountingUsecase) CreateUserV1(userData *v1.UserData) error {
	err := u.accountingRepo.CreateUserV1(userData)
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}

func (u *AccountingUsecase) CreateTaskV1(data *v1.TaskData) error {
	err := u.accountingRepo.CreateTaskV1(data)
	if err != nil {
		log.Println(err.Error())
	}
	return nil
}
func (u *AccountingUsecase) CreateTaskV2(data *v2.TaskData) error {
	err := u.accountingRepo.CreateTaskV2(data)
	if err != nil {
		log.Println(err.Error())
	}
	return nil
}

func (u *AccountingUsecase) AssignTaskV1(data *v1.TaskAssigneeData) error {
	err := u.accountingRepo.UpdateTaskUser(data.PublicTaskID, data.PublicUserID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
func (u *AccountingUsecase) CompleteTask(data *v1.TaskAssigneeData) error {
	err := u.accountingRepo.UpdateTaskStatus(data.PublicTaskID, 1)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (u *AccountingUsecase) ApplyTransaction(userID, taskID uuid.UUID, transactionType, description string) error {
	if description == "" {
		taskDescription, err := u.accountingRepo.GetTaskDescription(taskID)
		if err != nil {
			log.Println(err.Error())
		}
		description = taskDescription
	}

	balance, err := u.accountingRepo.GetUserBalance(userID)
	if err != nil {
		return err
	}

	switch transactionType {
	case "increase":
		val := helper.Random(20, 40)
		err := u.accountingRepo.IncreaseTransaction(userID, taskID, description, val, balance+val)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	case "decrease":
		val := helper.Random(10, 20)
		err := u.accountingRepo.DecreaseTransaction(userID, taskID, description, val, balance-val)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	return nil
}

func (u *AccountingUsecase) CloseBillingCycle() error {
	usersIDs, err := u.accountingRepo.GetUsersIDs()
	if err != nil {
		return err
	}

	for _, userID := range usersIDs {
		success, err := u.accountPayout(userID)
		if err != nil {
			log.Println(err)
			continue
		}
		if success {
			_ = u.sendMail(userID)
		}
	}

	return nil
}

func (u *AccountingUsecase) sendMail(userID uuid.UUID) error {
	//TODO: send email logic

	log.Println("Email sent")
	return nil
}

func (u *AccountingUsecase) accountPayout(userID uuid.UUID) (bool, error) {
	balance, err := u.accountingRepo.GetUserBalance(userID)
	if err != nil {
		return false, err
	}
	if balance <= 0 {
		//TODO: nothing
		log.Println("User balance is below or equals 0. Balance = ", balance)
		return false, nil
	}

	err = u.accountingRepo.PayoutTransaction(userID, "Payout ", balance)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	return true, nil
}

func (u *AccountingUsecase) GetAccountingStatistic() (*model.ProfitStatistics, error) {
	profit, err := u.accountingRepo.GetProfit()
	if err != nil {
		return nil, err
	}

	return &model.ProfitStatistics{
		Profit: profit,
	}, nil
}

func (u *AccountingUsecase) GetUserBalanceWithLogForLast24Hours(userID uuid.UUID) (*model.UserDailyStatistics, error) {
	balance, err := u.accountingRepo.GetUserBalance(userID)
	if err != nil {
		return nil, err
	}

	transactions, err := u.accountingRepo.GetUserTransactions(userID)
	if err != nil {
		return nil, err
	}

	return &model.UserDailyStatistics{
		Balance:      balance,
		Transactions: transactions,
	}, nil
}
