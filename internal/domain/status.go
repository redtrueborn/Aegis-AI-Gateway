package domain
type Status string


const (
	StatusUnknown   Status = "unknown"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)