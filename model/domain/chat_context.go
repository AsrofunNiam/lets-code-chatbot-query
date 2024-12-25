package domain

import "time"

type ChatContext struct {
	UserID  string
	History []ChatHistoryEntry
}

type ChatHistoryEntry struct {
	Type      string
	Content   string
	Timestamp time.Time
}
