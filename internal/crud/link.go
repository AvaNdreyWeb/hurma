package crud

import (
	"context"
	"errors"
	"fmt"
	"hurma/internal/models"
	"hurma/internal/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// cachedData, err := redisClient.Get(linkId).Result()
	// if err == nil {
	// 	var link models.Link
	// 	err := json.Unmarshal([]byte(cachedData), &link)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &link, nil
	// }
	// if not in cache
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

func (lm *LinkManager) GetByShortUrl(shortUrl string, cl *mongo.Client) (*models.Link, error) {
	coll := cl.Database("hurma").Collection("links")
	filter := bson.D{{Key: "shortUrl", Value: shortUrl}}
	link := new(models.Link)
	log.Println("short url", shortUrl)
	err := coll.FindOne(context.TODO(), filter).Decode(link)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.Link{}, err
		}
		log.Fatal(err)
	}
	log.Println("SUCCESS short url")

	return link, nil
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

func (lm *LinkManager) GetLinkStatistics(link *models.Link, days int, cl *mongo.Client) ([]models.DailyDTO, error) {
	startString := link.Expires.CreatedAt
	layout := "2006-01-02T15:04:05Z"

	startDate, err := time.Parse(layout, startString)
	if err != nil {
		return []models.DailyDTO{}, err
	}

	daily := link.Clicks.Daily
	diff := 0
	if days < len(daily) {
		diff = len(daily) - days
		startDate = startDate.AddDate(0, 0, diff)
		daily = daily[diff:]
	} else if days > len(daily) {
		days = len(daily)
	}

	DataList := make([]models.DailyDTO, days)

	for i := 0; i < days; i++ {
		curDate := startDate.AddDate(0, 0, i)
		DataList[i] = models.DailyDTO{Clicks: daily[i], Date: curDate.Format(layout)}
	}

	return DataList, nil
}

func UpdateAll(cl *mongo.Client) {
	coll := cl.Database("hurma").Collection("links")
	filter := bson.M{}
	findOptions := options.Find()
	cursor, err := coll.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var link models.Link
		if err := cursor.Decode(&link); err != nil {
			log.Fatal(err)
		}

		var yesterdayClicks uint64
		n := len(link.Clicks.Daily)
		if n == 0 {
			yesterdayClicks = link.Clicks.Total
		} else {
			yesterdayClicks = link.Clicks.Total - link.Clicks.Daily[n-1]
		}

		update := bson.M{"$push": bson.M{"Clicks.Daily": yesterdayClicks}}

		_, err = coll.UpdateOne(context.TODO(), bson.M{"_id": link.Id}, update)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

}
