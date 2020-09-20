// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Status string

const (
	StatusWaiting  Status = "waiting"
	StatusWorking  Status = "working"
	StatusFinished Status = "finished"
)

func (e *Status) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Status(s)
	case string:
		*e = Status(s)
	default:
		return fmt.Errorf("unsupported scan type for Status: %T", src)
	}
	return nil
}

type Problem struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	TestName     string    `json:"test_name"`
	AuthorID     int64     `json:"author_id"`
	TimeLimit    float64   `json:"time_limit"`
	MemoryLimit  int32     `json:"memory_limit"`
	StackLimit   int32     `json:"stack_limit"`
	SourceSize   int32     `json:"source_size"`
	ConsoleInput bool      `json:"console_input"`
	Visible      bool      `json:"visible"`
}

type Submission struct {
	ID             int64          `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UserID         int64          `json:"user_id"`
	ProblemID      int64          `json:"problem_id"`
	Language       string         `json:"language"`
	Code           string         `json:"code"`
	Status         Status         `json:"status"`
	CompileError   sql.NullBool   `json:"compile_error"`
	CompileMessage sql.NullString `json:"compile_message"`
	Score          int32          `json:"score"`
	Visible        bool           `json:"visible"`
}

type SubmissionTest struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Done         bool      `json:"done"`
	Verdict      string    `json:"verdict"`
	Time         float64   `json:"time"`
	Memory       int32     `json:"memory"`
	Score        int32     `json:"score"`
	TestID       int64     `json:"test_id"`
	UserID       int64     `json:"user_id"`
	SubmissionID int64     `json:"submission_id"`
}

type Test struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Score     int32     `json:"score"`
	ProblemID int64     `json:"problem_id"`
	VisibleID int32     `json:"visible_id"`
}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Admin     bool      `json:"admin"`
	Proposer  bool      `json:"proposer"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Bio       string    `json:"bio"`
}
