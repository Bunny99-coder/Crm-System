package postgres

import (
	"crm-project/internal/models"
)

// EventRepository defines the interface for event data access
type EventRepository interface {
	CreateEvent(event *models.Event) error
	GetEventByID(id int) (*models.Event, error)
	GetEventsByDealID(dealID int) ([]models.Event, error)
	GetAllEvents() ([]models.Event, error)
	UpdateEvent(event *models.Event) error
	DeleteEvent(id int) error
	GetEventsForUser(userID int) ([]models.Event, error)
}

type TaskRepository interface {
    CreateTask(task *models.Task) error
    GetTaskByID(id int) (*models.Task, error)
    GetTasksByDealID(dealID int) ([]models.Task, error)
    GetAllTasks() ([]models.Task, error)
    UpdateTask(task *models.Task) error
    DeleteTask(id int) error
    GetTasksForUser(userID int) ([]models.Task, error)
    GetTasksByDealIDForUser(dealID int, userID int) ([]models.Task, error)
}

// CommLogRepository defines the interface for communication log data access


// CommLogRepository defines the interface for communication log data access
type CommLogRepository interface {
    CreateCommLog(log *models.CommLog) error
    GetCommLogByID(id int) (*models.CommLog, error)
    GetCommLogsByDealID(dealID int) ([]models.CommLog, error)
    GetCommLogsByContactID(contactID int) ([]models.CommLog, error) // Added
    GetAllCommLogs() ([]models.CommLog, error)
    UpdateCommLog(log *models.CommLog) error
    DeleteCommLog(id int) error
    GetCommLogsForUser(userID int) ([]models.CommLog, error)
}