package crud

import (
	"context"
	"errors"
	"hurma-service/hurma-service/models"
	"hurma-service/hurma-service/utils"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserManager struct {
}

var ErrEmailConflict = errors.New("user with this email already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrValidationFailed = errors.New("invalid email or password")

func (um *UserManager) Create(u *models.AuthUserDTO, cl *mongo.Client) error {
	user, _ := um.Get(u.Email, cl)
	if user != nil {
		return ErrEmailConflict
	}

	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	coll := cl.Database("hurma").Collection("users")
	doc := models.User{
		Email:    u.Email,
		Password: hash,
		Links:    []primitive.ObjectID{},
	}
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}
	log.Printf("Inserted user with id: %v\n", result.InsertedID)

	return nil
}

func (um *UserManager) Validate(u *models.AuthUserDTO, cl *mongo.Client) error {
	user, err := um.Get(u.Email, cl)
	if err != nil {
		return ErrValidationFailed
	}

	if !utils.CheckPasswordHash(u.Password, user.Password) {
		return ErrValidationFailed
	}

	return nil
}

func (um *UserManager) Get(email string, cl *mongo.Client) (*models.User, error) {
	coll := cl.Database("hurma").Collection("users")
	user := new(models.User)
	filter := bson.D{{Key: "email", Value: email}}
	err := coll.FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		log.Fatal(err)
	}

	return user, nil
}

func (um *UserManager) AddLink(email string, linkId primitive.ObjectID, cl *mongo.Client) error {
	coll := cl.Database("hurma").Collection("users")
	_, err := um.Get(email, cl)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.M{"$push": bson.M{"links": linkId}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Append user links with linkId: %v\n", linkId)

	return nil
}
