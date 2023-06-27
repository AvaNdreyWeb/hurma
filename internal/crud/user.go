package crud

import (
	"context"
	"errors"
	"hurma/internal/models"
	"hurma/internal/utils"

	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserManager struct {
}

var ErrEmailConflict = errors.New("user with this email already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrPageNotFound = errors.New("page not found")
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
	_, err := um.Get(email, cl)
	if err != nil {
		return err
	}
	coll := cl.Database("hurma").Collection("users")
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.M{"$push": bson.M{"links": linkId}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Append user links with linkId: %v\n", linkId)

	return nil
}

func (um *UserManager) GetLinks(email string, page int, cl *mongo.Client) ([]models.TableLinkDTO, error) {
	user, err := um.Get(email, cl)
	if err != nil {
		return []models.TableLinkDTO{}, err
	}
	count := len(user.Links)
	start := (page - 1) * 10
	if start >= count {
		return []models.TableLinkDTO{}, ErrPageNotFound
	}
	var end int
	if start+10 > count {
		end = count
	} else {
		end = start + 10
	}

	linksId := user.Links[start:end]
	lm := new(LinkManager)
	links := lm.GetLinksByIdList(linksId, cl)
	return links, nil
}

func (um *UserManager) Subscribe(email string, cl *mongo.Client) error {
	_, err := um.Get(email, cl)
	if err != nil {
		return err
	}
	coll := cl.Database("hurma").Collection("users")
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.M{"$set": bson.M{"subscription": true}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	log.Printf("user subscribed to statistics: %s\n", email)

	return nil
}

func (um *UserManager) Unsubscribe(email string, cl *mongo.Client) error {
	_, err := um.Get(email, cl)
	if err != nil {
		return err
	}
	coll := cl.Database("hurma").Collection("users")
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.M{"$set": bson.M{"subscription": false}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	log.Printf("user unsubscribed from statistics: %s\n", email)

	return nil
}

func (um *UserManager) DeleteFromLinks(email string, id primitive.ObjectID, cl *mongo.Client) error {
	user, err := um.Get(email, cl)
	if err != nil {
		return err
	}
	linksId := make([]primitive.ObjectID, 0)
	for _, linkId := range user.Links {
		if linkId != id {
			linksId = append(linksId, linkId)
		}
	}

	coll := cl.Database("hurma").Collection("users")
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.M{"$set": bson.M{"links": linksId}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("link deleted from user: %v\n", id)

	return nil
}

func (um *UserManager) StatisticsAccess(email string, id primitive.ObjectID, cl *mongo.Client) bool {
	user, err := um.Get(email, cl)
	if err != nil {
		return false
	}
	for _, linkId := range user.Links {
		if linkId != id {
			return true
		}
	}
	return false
}
