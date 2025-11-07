package repository

import (
	"context"
	"errors"
	"denet/internal/model"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrTaskNotFound    = errors.New("task not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	CreateWithPassword(ctx context.Context, user *model.User, password string) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateBalance(ctx context.Context, id string, newBalance int) error
	SetReferrer(ctx context.Context, userID, referrerID string) error
	GetLeaderboard(ctx context.Context, limit int) ([]model.LeaderboardUser, error)
	VerifyPassword(ctx context.Context, username, password string) (*model.User, error)
}

type TaskRepository interface {
	GetByID(ctx context.Context, id string) (*model.Task, error)
	GetAll(ctx context.Context) ([]model.Task, error)
}

type UserTaskRepository interface {
	CompleteTask(ctx context.Context, userID, taskID string) error
	GetCompletedTasks(ctx context.Context, userID string) ([]model.UserTask, error)
	IsTaskCompleted(ctx context.Context, userID, taskID string) (bool, error)
}

type TransactionRepository interface {
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}

type UnitOfWork interface {
	Users() UserRepository
	Tasks() TaskRepository
	UserTasks() UserTaskRepository
	Transactions() TransactionRepository
	Close() error
}