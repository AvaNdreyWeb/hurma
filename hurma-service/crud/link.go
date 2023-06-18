package crud

import (
	"context"
	"errors"
	"fmt"
	"hurma-service/hurma-service/models"
	"hurma-service/hurma-service/utils"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkManager struct {
}

var ErrLinkConflict = errors.New("this link already exists")

func (lm *LinkManager) Create(l *models.CreateLinkDTO, cl *mongo.Client) (primitive.ObjectID, error) {
	if lm.FullExists(l.FullUrl, cl) {
		return primitive.ObjectID{}, ErrLinkConflict
	}

	var shortUrl string
	for {
		shortUrl = utils.ShortenURL()
		if !lm.ShortExists(shortUrl, cl) {
			break
		}
	}

	coll := cl.Database("hurma").Collection("links")
	doc := models.Link{
		Title:    l.Title,
		ShortUrl: shortUrl,
		FullUrl:  l.FullUrl,
		Expires: models.ExpireDate{
			CreatedAt: l.CreatedAt,
			ExpiresAt: l.ExpiresAt,
		},
		Clicks: models.ClickStat{
			Daily: []uint64{},
		},
	}
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	log.Printf("Inserted link with id: %v\n", result.InsertedID)
	linkId := result.InsertedID.(primitive.ObjectID)
	return linkId, nil
}

func (lm *LinkManager) Edit(u *models.AuthUserDTO, cl *mongo.Client) error {

	return nil
}

func (lm *LinkManager) FullExists(fullUrl string, cl *mongo.Client) bool {
	coll := cl.Database("hurma").Collection("links")

	link := new(models.Link)
	filter := bson.D{{Key: "fullUrl", Value: fullUrl}}
	err := coll.FindOne(context.TODO(), filter).Decode(link)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}

	return true
}

func (lm *LinkManager) ShortExists(shortUrl string, cl *mongo.Client) bool {
	coll := cl.Database("hurma").Collection("links")
	link := new(models.Link)
	filter := bson.D{{Key: "shortUrl", Value: shortUrl}}
	err := coll.FindOne(context.TODO(), filter).Decode(link)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}

	return true
}

func (lm *LinkManager) GetLinksByIdList(linksId []primitive.ObjectID, cl *mongo.Client) []models.TableLinkDTO {
	links := make([]models.TableLinkDTO, 0)
	for _, id := range linksId {
		links = append(links, lm.GetByID(id, cl))
	}
	return links
}

func (lm *LinkManager) GetByID(linkId primitive.ObjectID, cl *mongo.Client) models.TableLinkDTO {
	coll := cl.Database("hurma").Collection("links")
	filter := bson.D{{Key: "_id", Value: linkId}}
	link := new(models.Link)
	err := coll.FindOne(context.TODO(), filter).Decode(link)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.TableLinkDTO{}
		}
		log.Fatal(err)
	}
	stringId := fmt.Sprintf("%q", linkId.Hex())
	l := len(stringId)
	return models.TableLinkDTO{
		Id:          stringId[1 : l-1],
		Title:       link.Title,
		ShortUrl:    link.ShortUrl,
		ExpiresAt:   link.Expires.ExpiresAt,
		ClicksTotal: link.Clicks.Total,
	}
}
