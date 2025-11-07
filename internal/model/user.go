package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Balance      int       `json:"balance" db:"balance"`
	ReferrerID   *string   `json:"referrer_id,omitempty" db:"referrer_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Task struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Points      int       `json:"points" db:"points"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UserTask struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	TaskID    string    `json:"task_id" db:"task_id"`
	Completed bool      `json:"completed" db:"completed"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type LeaderboardUser struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Balance  int    `json:"balance" db:"balance"`
	Rank     int    `json:"rank" db:"rank"`
}

type UserStatus struct {
	User           *User      `json:"user"`
	CompletedTasks []UserTask `json:"completed_tasks"`
	TotalPoints    int        `json:"total_points"`
}

type CompleteTaskRequest struct {
	TaskID string `json:"task_id" binding:"required"`
}

type SetReferrerRequest struct {
	ReferrerID string `json:"referrer_id" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}