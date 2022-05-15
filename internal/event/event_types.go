package event

// CUD
const (
	TASK_CREATED = "task.created"
	TASK_UPDATED = "task.updated"
	TASK_DELETED = "task.deleted"

	USER_CREATED = "user.created"
	USER_UPDATED = "user.updated"
	USER_DELETED = "user.deleted"
)

// BE
const (
	NEW_TASK_ADDED = "task.new_task_added"
	TASKS_SHUFFLED = "task.shuffled"
	TASK_COMPLETED = "task.completed"

	TRANSACTION_APPLIED = "accounting.transaction_applied"
	PAYMENT_COMPLETED   = "accounting.payment_completed"
)
