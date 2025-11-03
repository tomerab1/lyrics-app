package models

type User struct {
	Id   string `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
}
