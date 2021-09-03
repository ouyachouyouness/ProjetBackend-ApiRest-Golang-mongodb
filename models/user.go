package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
type User struct {
	Created   time.Time          `json:"created",bson:"created"`
	Name      string             `json:"name",bson:"name"`
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Password  string             `json:"password,omitempty" bson:"Password"`
	isActive  bool               `json:"isActive,omitempty" bson:"isActive"`
	balance   string             `json:"balance,omitempty"bson:"balance"`
	age       float32            `json:"age,omitempty"bson:"age"`
	gender    string             `json:"gender,omitempty"bson:"gender"`
	company   string             `json:"company,omitempty"bson:"company"`
	Email     string             `json:"email,omitempty" bson:"email"`
	phone     string             `json:"phone,omitempty"bson:"phone"`
	adress    string             `json:"adress,omitempty"bson:"adress"`
	about     string             `json:"about,omitempty"bson:"about"`
	registred string             `json:"registred,omitempty"bson:"registred"`
	latitude  string             `json:"latitude,omitempty"bson:"latitude"`
	tags      []string           `json:"tags,omitempty"bson:"tags"`
	friends   []string           `json:"friends,omitempty"bson:"friends"`
	data      string             `json:"data,omitempty"bson:"data"`
}

type Session struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email   string             `json:"email",bson:"email"`
	Expiry  time.Time          `json:"expiry",bson:"expiry"`
	Created time.Time          `json:"created",bson:"created"`
}

type UpdateUser struct {
	Name      string             `json:"name",bson:"name"`
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Password  string             `json:"password,omitempty" bson:"Password"`
	isActive  bool               `json:"isActive,omitempty" bson:"isActive"`
	balance   string             `json:"balance,omitempty"bson:"balance"`
	age       float32            `json:"age,omitempty"bson:"age"`
	gender    string             `json:"gender,omitempty"bson:"gender"`
	company   string             `json:"company,omitempty"bson:"company"`
	Email     string             `json:"email,omitempty" bson:"email"`
	phone     string             `json:"phone,omitempty"bson:"phone"`
	adress    string             `json:"adress,omitempty"bson:"adress"`
	about     string             `json:"about,omitempty"bson:"about"`
	registred string             `json:"registred,omitempty"bson:"registred"`
	latitude  string             `json:"latitude,omitempty"bson:"latitude"`
	tags      []string           `json:"tags,omitempty"bson:"tags"`
	friends   []string           `json:"friends,omitempty"bson:"friends"`
	data      string             `json:"data,omitempty"bson:"data"`
}

func (u *UpdateUser) Valid() error {
	return nil
}

func (u *UpdateUser) Update() bson.M {
	update := bson.M{}
	if u.Name != "" {
		update["Name"] = u.Name
	}
	if u.age != 0 {
		update["model"] = u.age
	}
	if u.adress != "" {
		update["adress"] = u.adress
	}
	if u.phone != "" {
		update["phone"] = u.phone
	}
	if len(update) == 0 {
		return nil
	}
	return bson.M{"$set": update}
}
