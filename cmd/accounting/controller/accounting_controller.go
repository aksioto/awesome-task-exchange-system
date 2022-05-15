package controller

import (
	"encoding/json"
	"github.com/aksioto/awesome-task-exchange-system/cmd/accounting/usecase"
	"github.com/aksioto/awesome-task-exchange-system/internal/event"
	"github.com/aksioto/awesome-task-exchange-system/internal/service/rabbitmq"
	"log"
)

type AccountingController struct {
	tasksUsecase    *usecase.AccountingUsecase
	rabbitmqService *rabbitmq.RabbitmqService
}

func NewAccountingController(tasksUsecase *usecase.AccountingUsecase, rabbitmqService *rabbitmq.RabbitmqService) *AccountingController {
	return &AccountingController{
		tasksUsecase:    tasksUsecase,
		rabbitmqService: rabbitmqService,
	}
}

// MQ
func (tc *AccountingController) HandleNewTasks(body []byte) {
	e := &rabbitmq.Event{}
	err := json.Unmarshal(body, e)
	if err != nil {
		log.Println("Failed to unmarshal message")
	}

	switch e.Name {
	case event.NEW_TASK_ADDED:
		tc.ApplyTransaction(e)
	case event.TASKS_SHUFFLED:
		tc.AccountPayout(e)
	case event.TASK_COMPLETED:
		tc.AccountPayout(e)
	}
}

func (tc *AccountingController) ApplyTransaction(e *rabbitmq.Event) {
	switch e.Version {
	case 1:
		log.Println("[1]", e.Name, e.Version)
	case 2:
		log.Println("[2]", e.Name, e.Version)
	}
}

func (tc *AccountingController) AccountPayout(e *rabbitmq.Event) {
	log.Println(e.Name, e.Version)
}

func (tc *AccountingController) CloseBillingCycle(e *rabbitmq.Event) {
	log.Println(e.Name, e.Version)
}

func (tc *AccountingController) HandleShuffledTasks(body []byte) {
	e := &rabbitmq.Event{}
	err := json.Unmarshal(body, e)
	if err != nil {
		log.Println("Failed to unmarshal message")
	}

	switch e.Name {
	case event.NEW_TASK_ADDED:
		tc.ApplyTransaction(e)
	case event.TASKS_SHUFFLED:
		tc.AccountPayout(e)
	case event.TASK_COMPLETED:
		tc.AccountPayout(e)
	}
}
