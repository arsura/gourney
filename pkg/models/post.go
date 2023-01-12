package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	POST_TITLE               = "title"
	POST_CONTENT             = "content"
	POST_SOCIAL_NETWORK_TYPE = "social_network_type"
)

type PostSocialNetworkType string

const (
	PostSocialNetworkTypeFacebook PostSocialNetworkType = "facebook"
	PostSocialNetworkTypeTwitter  PostSocialNetworkType = "twitter"
)

type Post struct {
	Id                primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Title             string                `json:"title" bson:"title"`
	Content           string                `json:"content" bson:"content"`
	SocialNetworkType PostSocialNetworkType `json:"social_network_type" bson:"social_network_type"`
	CreatedAt         time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at" bson:"updated_at"`
}
