package models

import "time"

const READING_COLLECTION = "readings"

type ParameterValue struct {
	Label string  `json:"label" bson:"label"`
	Value float64 `json:"value" bson:"value"`
}

type Reading struct {
	Uid      string           `json:"uid" bson:"uid"`
	DateTime int64            `json:"datetime" bson:"datetime"`
	Values   []ParameterValue `json:"values" bson:"values"`
}

func NewReadingWithDatetime(uid string, datetime int64, values ...ParameterValue) *Reading {
	return &Reading{
		Uid:      uid,
		DateTime: datetime,
		Values:   values,
	}
}

func NewReading(uid string, values ...ParameterValue) *Reading {
	return &Reading{
		Uid:      uid,
		DateTime: time.Now().Unix(),
		Values:   values,
	}
}

// db.readings.insert({
// 	"uid": "abcdef",
// 	"datetime": 0,
// 	"values": [
// 	{
// 		"label": "temperature",
// 		"value": 69.9
// 	}
// 	]
// })
