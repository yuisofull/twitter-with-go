package tweetmodel

type Filter struct {
	FakeUserID string `json:"-" form:"user_id"`
	UserID     int    `json:"user_id,omitempty" form:"-"`
	Status     []int  `json:"-"`
	Search     string `json:"search,omitempty" form:"search"`
}
