package notify

import (
 "context"
 "time"
)

type Status string
const (
 StatusQueued Status = "queued"
 StatusAccepted Status = "accepted"
 StatusRetryWait Status = "retry_wait"
 StatusFailed Status = "failed"
)

type Notification struct {
 ID string `json:"id"`
 IdempotencyKey string `json:"-"`
 Recipient string `json:"recipient_ref"`
 Channel string `json:"channel"`
 Template string `json:"template_id"`
 Payload map[string]string `json:"payload"`
 Status Status `json:"status"`
 Attempts int `json:"attempts"`
 NextAttemptAt time.Time `json:"next_attempt_at"`
 LastError string `json:"last_error,omitempty"`
 CreatedAt time.Time `json:"created_at"`
}

type ProviderRequest struct {
 NotificationID string `json:"notification_id"`
 Recipient string `json:"recipient_ref"`
 Channel string `json:"channel"`
 Template string `json:"template_id"`
 Payload map[string]string `json:"payload"`
}
type ProviderResult struct { ProviderMessageID string }
type Provider interface { Send(context.Context, ProviderRequest) (ProviderResult,error) }
