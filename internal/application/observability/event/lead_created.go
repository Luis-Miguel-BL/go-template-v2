package event

import "time"

type LeadCreated struct {
	LeadID    string
	CreatedAt time.Time
}

func (e LeadCreated) Name() string {
	return "LeadCreated"
}

func (e LeadCreated) Attributes() map[string]any {
	return map[string]any{
		"lead_id":    e.LeadID,
		"created_at": e.CreatedAt.Format(time.RFC3339),
	}
}
