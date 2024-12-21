package auction

import (
	"context"
	"time"

	"github.com/vs0uz4/labs-auction/config/logger"
	"github.com/vs0uz4/labs-auction/internal/entity/auction_entity"
	"github.com/vs0uz4/labs-auction/internal/infra/utils"
	"github.com/vs0uz4/labs-auction/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func (ar *AuctionRepository) ClosingExpiredAuctions(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(utils.GetMaxBatchSizeInterval())
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				filter := bson.M{"status": 0}
				cursor, err := ar.Collection.Find(ctx, filter)
				if err != nil {
					logger.Error("Error fetching auctions", err)
					continue
				}
				defer cursor.Close(ctx)

				var auctions []AuctionEntityMongo
				if err = cursor.All(ctx, &auctions); err != nil {
					logger.Error("Error decoding auctions", err)
					continue
				}

				for _, auction := range auctions {
					if utils.IsAuctionExpired(auction.Timestamp) {
						update := bson.M{"$set": bson.M{"status": 1}}
						_, err := ar.Collection.UpdateOne(ctx, bson.M{"_id": auction.Id}, update)
						if err != nil {
							logger.Error("Error updating auction", err)
						}
					}
				}

			case <-ctx.Done():
				return
			}
		}
	}()
}
