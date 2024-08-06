package models

import "time"

type User struct {
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Feedback represents the Feedback table.
type Feedback struct {
	FeedbackID     int       `json:"feedback_id"`
	UserID         int       `json:"user_id"`
	Name           string    `json:"name"`
	Age            int       `json:"age"`
	Occupation     string    `json:"occupation"`
	Email          string    `json:"email"`
	WouldRecommend *bool     `json:"would_recommend"`
	Suggestion     string    `json:"suggestion"`
	Likes          string    `json:"likes"`
	CreatedAt      time.Time `json:"created_at"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
