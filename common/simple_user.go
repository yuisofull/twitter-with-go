package common

type SimpleUser struct {
	SQLModel  `json:",inline"`
	LastName  string `json:"last_name,omitempty" gorm:"column:last_name;"`
	FirstName string `json:"first_name,omitempty" gorm:"column:first_name;"`
	Role      string `json:"role,omitempty" gorm:"column:role;"`
	Avatar    *Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (SimpleUser) TableName() string {
	return "users"
}

func (u *SimpleUser) Mask(isAdmin bool) {
	u.GenUID(DbTypeUser)
}
