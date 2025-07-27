package main

import(
   "time"
)

type TaskStatus int

// Task represents a transcription task
type Task struct {
	OBJID         int       `json:"objid" db:"objid"`
	Filename      string    `json:"filename" db:"filename"`
	Label         string    `json:"label" db:"label"`
	SSOAccount    string    `json:"sso_account" db:"sso_account"`
	Status        TaskStatus `json:"status" db:"status"`
	PID           int       `json:"pid" db:"pid"`
	Transcribe    string    `json:"transcribe" db:"transcribe"`
	ContentLength int       `json:"content_length" db:"content_length"`
	Diarize       int       `json:"diarize" db:"diarize"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}