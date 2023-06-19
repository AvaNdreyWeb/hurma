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

func (lm *LinkManager) EditTitle(title string, id primitive.ObjectID, cl *mongo.Client) error {
	coll := cl.Database("hurma").Collection("links")
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: title}}}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("link title updated: %v\n", id)

	return nil
}

func (lm *LinkManager) EditExpires(expiresAt string, id primitive.ObjectID, cl *mongo.Client) error {
	link := lm.GetByID(id, cl)
	coll := cl.Database("hurma").Collection("links")
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "expires", Value: bson.D{{Key: "expiresAt", Value: expiresAt}, {Key: "createdAt", Value: link.Expires.CreatedAt}}}}}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("link expire date is updated: %v\n", id)

	return nil
}

func (lm *LinkManager) Delete(email string, id primitive.ObjectID, cl *mongo.Client) error {
	coll := cl.Database("hurma").Collection("links")
	filter := bson.D{{Key: "_id", Value: id}}

	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	log.Printf("link deleted: %v\n", id)
	um := new(UserManager)
	err = um.DeleteFromLinks(email, id, cl)
	if err != nil {
		return err
	}
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
		link := lm.GetByID(id, cl)
		stringId := fmt.Sprintf("%q", id.Hex())
		length := len(stringId)
		l := models.TableLinkDTO{
			Id:          stringId[1 : length-1],
			Title:       link.Title,
			ShortUrl:    link.ShortUrl,
			ExpiresAt:   link.Expires.ExpiresAt,
			ClicksTotal: link.Clicks.Total,
		}
		links = append(links, l)
	}
	return links
}

func (lm *LinkManager) GetByID(linkId primitive.ObjectID, cl *mongo.Client) *models.Link {
	coll := cl.Database("hurma").Collection("links")
	filter := bson.D{{Key: "_id", Value: linkId}}
	link := new(models.Link)
	err := coll.FindOne(context.TODO(), filter).Decode(link)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.Link{}
		}
		log.Fatal(err)
	}

	return link
}

func (lm *LinkManager) GetFullUrl(shortUrl string, cl *mongo.Client) (*models.Link, error) {
	coll := cl.Database("hurma").Collection("links")
	filter := bson.D{{Key: "shortUrl", Value: shortUrl}}
	link := new(models.Link)
	err := coll.FindOne(context.TODO(), filter).Decode(link)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.Link{}, err
		}
		log.Fatal(err)
	}

	return link, nil
}

func (lm *LinkManager) IncTotal(id primitive.ObjectID, cl *mongo.Client) error {
	coll := cl.Database("hurma").Collection("links")
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "clicks", Value: bson.D{{Key: "total", Value: 1}}}}}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("link inc updated: %v\n", id)

	return nil
}
