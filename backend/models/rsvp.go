package models

import (
	"time"
)

type RSVP struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     *uint     `json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	GuestName  string    `json:"guest_name"`
	GuestEmail string    `json:"guest_email"`
	GuestPhone string    `json:"guest_phone"`
	Status     string    `json:"status"`
	GuestCount int       `json:"guest_count"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

//
//// Get display name
//func (r *RSVP) GetDisplayName() string {
//	if r.User != nil {
//		return r.User.FullName
//	}
//	return r.GuestName
//}
//
//// Get display email
//func (r *RSVP) GetDisplayEmail() string {
//	if r.User != nil {
//		return r.User.Email
//	}
//	return r.GuestEmail
//}
//
//// Get display phone
//func (r *RSVP) GetDisplayPhone() string {
//	if r.User != nil {
//		return r.User.Phone
//	}
//	return r.GuestPhone
//}
