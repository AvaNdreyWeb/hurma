package crud

import (
	"context"
	"errors"
	"hurma-service/hurma-service/models"
	"hurma-service/hurma-service/utils"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserManager struct {
}

var ErrUsernameConflict = errors.New("user with this username already exists")

func (um *UserManager) Create(u *models.AuthUserDTO, cl *mongo.Client) error {

	if um.Exists(u.Username, cl) {
		return ErrUsernameConflict
	}

	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	coll := cl.Database("hurma").Collection("users")
	doc := models.User{
		Username: u.Username,
		Password: hash,
	}
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}
	log.Printf("Inserted user with id: %v\n", result.InsertedID)

	return nil
}

func (um *UserManager) Validate(u *models.AuthUserDTO, cl *mongo.Client) error {
	// log.Println("Getting user with username", u.Username, "...")
	// log.Println("utils.CheckPasswordHash(u.Password, user.Password)...")

	return nil
}

func (um *UserManager) Exists(username string, cl *mongo.Client) bool {
	coll := cl.Database("hurma").Collection("users")

	user := new(models.User)
	filter := bson.D{{Key: "username", Value: username}}
	err := coll.FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}

	return true
}
