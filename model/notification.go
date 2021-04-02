package model

type Notification struct {
	ID             uint
	RecipientEmail string
	Type           string
}

func NewInvitationNotification(email string) *Notification {
	return &Notification{
		RecipientEmail: email,
		Type:           "Invitation",
	}
}
