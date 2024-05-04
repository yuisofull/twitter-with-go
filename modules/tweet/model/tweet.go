package tweetmodel

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"twitter/common"
)

const (
	EntityName = "Tweet"
	TableName  = "tweets"
	Index      = "tweet_index"
)

type Tweet struct {
	common.SQLModel `json:",inline"`
	UserID          int                `json:"-" gorm:"column:user_id;"`
	User            *common.SimpleUser `json:"user" gorm:"preload:false;foreignKey:Id;"`
	Text            string             `json:"text_content" gorm:"column:text_content;"`
	Images          *common.Images     `json:"image" gorm:"column:image;" form:"-"`
}

func (Tweet) TableName() string { return TableName }

func (r *Tweet) Mask(isAdminOrOwner bool) {
	r.GenUID(common.DbTypeTweet)

	if u := r.User; u != nil {
		u.Mask(false)
	}
}

type TweetES struct {
	Id     int    `json:"id"`
	UserID int    `json:"user_id"`
	Text   string `json:"text_content"`
	Images string `json:"images"`
}

func (tw *TweetES) ToTweet() Tweet {
	uid := common.NewUID(uint32(tw.Id), common.DbTypeTweet, 1)
	fakeUserID := common.NewUID(uint32(tw.UserID), common.DbTypeUser, 1)
	image := &common.Images{}
	err := json.Unmarshal([]byte(tw.Images), image)
	if err != nil {
		log.Println("error unmarshal image", err)
	}

	return Tweet{
		SQLModel: common.SQLModel{
			Id:     int(uid.GetLocalID()),
			FakeID: &uid,
			Status: 1,
		},
		UserID: tw.UserID,
		User: &common.SimpleUser{
			SQLModel: common.SQLModel{
				FakeID: &fakeUserID,
				Status: 1,
			},
		},
		Text:   tw.Text,
		Images: image,
	}
}

type TweetCreate struct {
	common.SQLModel `json:",inline" swaggerignore:"true"`
	UserID          int            `json:"-" gorm:"column:user_id;"`
	Text            string         `json:"text_content" gorm:"column:text_content;"`
	ImageUIDs       []string       `json:"imageIDs" gorm:"-" form:"imageIDs"`
	Images          *common.Images `json:"-" gorm:"column:images;" form:"-" swaggerignore:"true"`
}

func (TweetCreate) TableName() string { return TableName }

func (data *TweetCreate) Mask(isAdminOrOwner bool) {
	data.GenUID(common.DbTypeTweet)
}

func (data *TweetCreate) Validate() error {
	data.Text = strings.TrimSpace(data.Text)
	if data.Text == "" && len(data.ImageUIDs) == 0 {
		return ErrTextOrImageEmpty
	}
	return nil
}

type UpdateTweet struct {
	Text      string        `json:"text_content" gorm:"column:text_content;"`
	ImageUIDs []common.UID  `json:"-" gorm:"-" form:"images"`
	Images    common.Images `json:"image" gorm:"column:image;" form:"-"`
}

func (UpdateTweet) TableName() string { return Tweet{}.TableName() }

var (
	ErrNameIsEmpty      = errors.New("name can not be empty")
	ErrTextOrImageEmpty = errors.New("text or image can not be empty")
)
