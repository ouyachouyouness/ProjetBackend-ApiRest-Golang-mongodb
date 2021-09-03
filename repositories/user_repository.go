package repositories

import (
	"context"
	// "errors"
	// "fmt"
	"time"
	// "log"
	"goproj/models"
	// "database/sql"
	// "github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Database
}

//////////////////////////////////////////////////////////////////////////////////////

func NewInstanceOfUserRepository(db *mongo.Database) UserRepository {
	return UserRepository{db: db}
}

/////////////////////////////////////GetUserByEmail/////////////////////////////////////////////////

func (u *UserRepository) GetUserByEmail(Email string) (bool, models.User, error) {
	var user models.User
	filter := bson.M{"email": Email}
	count, err := u.db.Collection("users").CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, models.User{}, err
	}
	if count != 1 {
		return false, models.User{}, nil
	}
	err = u.db.Collection("users").FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return false, models.User{}, err
	}

	return true, user, nil
}

/////////////////////////////////////GetUserById/////////////////////////////////////////////////

func (u *UserRepository) GetUserByID(ID string) (bool, models.User, error) {
	var user models.User
	filter := bson.M{"ID": ID}
	count, err := u.db.Collection("users").CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, models.User{}, err
	}
	if count != 1 {
		return false, models.User{}, nil
	}
	err = u.db.Collection("users").FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return false, models.User{}, err
	}

	return true, user, nil
}

//////////////////////////////////////DoesUserExist////////////////////////////////////////////////

func (u *UserRepository) DoesUserExist(email string) (bool, error) {
	exists, _, err := u.GetUserByEmail(email)
	return exists, err
}

////////////////////////////////SaveUser//////////////////////////////////////////////////////

func (u *UserRepository) SaveUser(user models.User) error {
	_, err := u.db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

////////////////////////////////////SaveSession//////////////////////////////////////////////////

func (u *UserRepository) SaveSession(session models.Session) (string, error) {
	insertResult, err := u.db.Collection("sessions").InsertOne(context.TODO(), session)
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

////////////////////////////////GetSessionById//////////////////////////////////////////////////////

func (u *UserRepository) GetSessionById(token string) (bool, models.Session, error) {
	docID, err := primitive.ObjectIDFromHex(token)
	if err != nil {
		return false, models.Session{}, err
	}

	var session models.Session
	filter := bson.M{
		"_id": docID,
		"expiry": bson.M{
			"$gte": time.Now(),
		},
	}

	count, err := u.db.Collection("sessions").CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, models.Session{}, err
	}
	if count != 1 {
		return false, models.Session{}, nil
	}
	err = u.db.Collection("sessions").FindOne(context.TODO(), filter).Decode(&session)
	if err != nil {
		return false, models.Session{}, err
	}
	return true, session, nil
}

//////////////////////////////////MarkSessionAsExpired////////////////////////////////////////////////////

func (u *UserRepository) MarkSessionAsExpired(authToken string) error {
	filter := bson.M{"_id": authToken}
	update := bson.M{"$set": bson.M{"expiry": time.Now()}}
	_, err := u.db.Collection("sessions").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

/////////////////////////////////Update/////////////////////////////////////////////////////

func (c *UserRepository) Update(email string, carID string, body models.UpdateUser) error {
	docID, err := primitive.ObjectIDFromHex(carID)
	if err != nil {
		return err
	}

	update := body.Update()
	if update == nil {
		return nil
	}
	filter := bson.M{"_id": docID, "email": email}

	_, err = c.db.Collection("user").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

////////////////////////////Delete//////////////////////////////////////////////////////////

func (c *UserRepository) Delete(email string, carID string) error {
	docID, err := primitive.ObjectIDFromHex(carID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": docID, "email": email}
	_, err = c.db.Collection("user").DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
