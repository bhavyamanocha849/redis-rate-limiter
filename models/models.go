package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApiConfig struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ApiConfigDetails ApiConfigDetails   `json:"api_config_details" bson:"api_config_details"`
}

type ApiConfigDetails struct {
	TimeWindow time.Duration `json:"time_window" bson:"time_window"`
	Capacity   int64         `json:"capacity" bson:"capacity"`
}
