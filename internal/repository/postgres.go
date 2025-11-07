package repository

import (
	"context"
	"denet/internal/store"
	"denet/internal/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type PostgresUnitOfWork struct {
	db store.Database
}

func NewPostgresUnitOfWork(db store.Database) UnitOfWork {
	return &PostgresUnitOfWork{db: db}
}

func (uow *PostgresUnitOfWork) Users() UserRepository {
	return &PostgresUserRepository{db: uow.db}
}

func (uow *PostgresUnitOfWork) Tasks() TaskRepository {
	return &PostgresTaskRepository{db: uow.db}
}

func (uow *PostgresUnitOfWork) UserTasks() UserTaskRepository {
	return &PostgresUserTaskRepository{db: uow.db}
}

func (uow *PostgresUnitOfWork) Transactions() TransactionRepository {
	return &PostgresTransactionRepository{db: uow.db}
}

func (uow *PostgresUnitOfWork) Close() error {
	return uow.db.Close()
}

type PostgresUserRepository struct {
	db store.Database
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *model.User) error {
	user.ID = uuid.New().String()
	query := `INSERT INTO users (id, username, email, balance) VALUES ($1, $2, $3, $4)`
	return r.db.Exec(ctx, query, user.ID, user.Username, user.Email, user.Balance)
}

func (r *PostgresUserRepository) CreateWithPassword(ctx context.Context, user *model.User, password string) error {
	user.ID = uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (id, username, email, password_hash, balance) VALUES ($1, $2, $3, $4, $5)`
	return r.db.Exec(ctx, query, user.ID, user.Username, user.Email, string(hashedPassword), user.Balance)
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, username, email, password_hash, balance, referrer_id, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Balance, &user.ReferrerID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id, username, email, password_hash, balance, referrer_id, created_at, updated_at FROM users WHERE username = $1`
	row := r.db.QueryRow(ctx, query, username)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Balance, &user.ReferrerID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, username, email, password_hash, balance, referrer_id, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Balance, &user.ReferrerID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *PostgresUserRepository) UpdateBalance(ctx context.Context, id string, newBalance int) error {
	query := `UPDATE users SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	return r.db.Exec(ctx, query, newBalance, id)
}

func (r *PostgresUserRepository) GetLeaderboard(ctx context.Context, limit int) ([]model.LeaderboardUser, error) {
	query := `SELECT id, username, balance, RANK() OVER (ORDER BY balance DESC) as rank FROM users ORDER BY balance DESC LIMIT $1`
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.LeaderboardUser
	for rows.Next() {
		var user model.LeaderboardUser
		if err := rows.Scan(&user.ID, &user.Username, &user.Balance, &user.Rank); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) SetReferrer(ctx context.Context, userID, referrerID string) error {
	_, err := r.GetByID(ctx, referrerID)
	if err != nil {
		return ErrUserNotFound
	}
	query := `UPDATE users SET referrer_id = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	return r.db.Exec(ctx, query, referrerID, userID)
}

func (r *PostgresUserRepository) VerifyPassword(ctx context.Context, username, password string) (*model.User, error) {
	user, err := r.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidPassword
	}
	return user, nil
}

type PostgresTaskRepository struct {
	db store.Database
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id string) (*model.Task, error) {
	query := `SELECT id, name, description, points, created_at FROM tasks WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var task model.Task
	err := row.Scan(&task.ID, &task.Name, &task.Description, &task.Points, &task.CreatedAt)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return &task, nil
}

func (r *PostgresTaskRepository) GetAll(ctx context.Context) ([]model.Task, error) {
	query := `SELECT id, name, description, points, created_at FROM tasks ORDER BY created_at`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.Points, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

type PostgresUserTaskRepository struct {
	db store.Database
}

func (r *PostgresUserTaskRepository) CompleteTask(ctx context.Context, userID, taskID string) error {
	userTaskID := uuid.New().String()
	query := `INSERT INTO user_tasks (id, user_id, task_id, completed) VALUES ($1, $2, $3, true)`
	return r.db.Exec(ctx, query, userTaskID, userID, taskID)
}

func (r *PostgresUserTaskRepository) GetCompletedTasks(ctx context.Context, userID string) ([]model.UserTask, error) {
	query := `SELECT id, user_id, task_id, completed, created_at FROM user_tasks WHERE user_id = $1 AND completed = true`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userTasks []model.UserTask
	for rows.Next() {
		var userTask model.UserTask
		if err := rows.Scan(&userTask.ID, &userTask.UserID, &userTask.TaskID, &userTask.Completed, &userTask.CreatedAt); err != nil {
			return nil, err
		}
		userTasks = append(userTasks, userTask)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userTasks, nil
}

func (r *PostgresUserTaskRepository) IsTaskCompleted(ctx context.Context, userID, taskID string) (bool, error) {
	var id string
	query := `SELECT id FROM user_tasks WHERE user_id = $1 AND task_id = $2 AND completed = true`
	row := r.db.QueryRow(ctx, query, userID, taskID)
	err := row.Scan(&id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

type PostgresTransactionRepository struct {
	db store.Database
}

func (r *PostgresTransactionRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	if err := fn(ctx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}