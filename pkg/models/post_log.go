package model

import (
	"time"

	"github.com/arsura/gourney/pkg/constant"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	POST_ID   = "post_id"
	EVENT     = "event"
	ACTION_AT = "action_at"
)

type PostLog struct {
	Id        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	PostId    primitive.ObjectID   `json:"post_id" bson:"post_id"`
	Event     constant.EventAction `json:"event" bson:"event"`
	ActionAt  time.Time            `json:"action_at" bson:"action_at"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
}
