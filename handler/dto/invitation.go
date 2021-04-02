package dto

import "time"

type CreateInvitationRequest struct {
	Email string `json:"email" valid:"email"`
}

type InvitationResponse struct {
	ID        uint      `json:id`
	Email     string    `json:"email"`
	InviterID uint      `json:"inviter_id"`
	CreatedAt time.Time `json:"created_at"`
}
