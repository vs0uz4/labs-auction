package auction

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vs0uz4/labs-auction/config/logger"
	"github.com/vs0uz4/labs-auction/internal/infra/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const URI = "mongodb://admin:admin@localhost:27016/auctions?authSource=admin"

type MockAuctionRepository struct {
	mock.Mock
}

func (m *MockAuctionRepository) Find(ctx context.Context, filter interface{}) ([]AuctionEntityMongo, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]AuctionEntityMongo), args.Error(1)
}

func (m *MockAuctionRepository) UpdateOne(ctx context.Context, filter, update interface{}) error {
	args := m.Called(ctx, filter, update)
	return args.Error(0)
}

func setupTestLoggerForClass() (*bytes.Buffer, func()) {
	var buffer bytes.Buffer
	writeSyncer := zapcore.AddSync(&buffer)

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	mockedLogger := zap.New(core)

	originalLogger := logger.GetLogger()
	logger.SetLogger(mockedLogger)

	restore := func() {
		logger.SetLogger(originalLogger)
		mockedLogger.Sync()
	}

	return &buffer, restore
}

func TestClosingExpiredAuctionsWithMock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockRepo := new(MockAuctionRepository)

	mockRepo.On("Find", mock.Anything, mock.Anything).Return([]AuctionEntityMongo{
		{Id: "a3739af2-9c8d-4ed5-b396-3e135c8153fc", Timestamp: time.Now().Add(-10 * time.Minute).Unix()},
		{Id: "5a0a1e96-6cb1-44ca-a88e-1966e94a7591", Timestamp: time.Now().Add(5 * time.Minute).Unix()},
	}, nil)

	filter := bson.M{"status": 0}
	update := bson.M{"$set": bson.M{"status": 1}}

	mockRepo.On("UpdateOne", mock.Anything, bson.M{"_id": "a3739af2-9c8d-4ed5-b396-3e135c8153fc"}, update).Return(nil)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				auctions, _ := mockRepo.Find(ctx, filter)
				for _, auction := range auctions {
					if utils.IsAuctionExpired(auction.Timestamp) {
						_ = mockRepo.UpdateOne(ctx, bson.M{"_id": auction.Id}, update)
					}
				}
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	time.Sleep(2 * time.Second)
	mockRepo.AssertCalled(t, "UpdateOne", mock.Anything, bson.M{"_id": "a3739af2-9c8d-4ed5-b396-3e135c8153fc"}, update)
	mockRepo.AssertNotCalled(t, "UpdateOne", mock.Anything, bson.M{"_id": "5a0a1e96-6cb1-44ca-a88e-1966e94a7591"}, update)
}

func TestClosingExpiredAuctionsMockingFetchingErrorWithMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	buffer, restoreLogger := setupTestLoggerForClass()
	defer restoreLogger()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("auctions")
	defer db.Drop(ctx)

	collection := db.Collection("auctions")

	auctionRepo := &AuctionRepository{
		Collection: collection,
	}

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, err := auctionRepo.Collection.Find(ctx, bson.M{"$invalid": "field"})
				if err != nil {
					logger.Error("Expected error in fetching auctions", err)
				}
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	time.Sleep(time.Second * 1)

	logs := buffer.String()
	assert.Contains(t, logs, "Expected error in fetching auctions")
	assert.Contains(t, logs, "unknown top level operator")
}

func TestClosingExpiredAuctionsMockingUpdateErrorWithMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	buffer, restoreLogger := setupTestLoggerForClass()
	defer restoreLogger()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("auctions")
	defer db.Drop(ctx)

	collection := db.Collection("auctions")

	_, err = collection.InsertOne(ctx, bson.M{
		"_id":          "expired-auction-id",
		"product_name": "Expired Product",
		"category":     "expired-auction-category",
		"description":  "expired-auction-description",
		"status":       0,
		"condition":    1,
		"timestamp":    time.Now().Add(-10 * time.Minute).Unix(),
	})
	if err != nil {
		t.Fatalf("Failed to insert expired auction: %v", err)
	}

	auctionRepo := &AuctionRepository{
		Collection: collection,
	}

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, err := auctionRepo.Collection.UpdateOne(ctx, bson.M{"_id": "expired-auction-id"}, bson.M{"$set": bson.M{"$status": 1}})
				if err != nil {
					logger.Error("Expected error in update auctions", err)
				}
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	time.Sleep(time.Second * 1)

	logs := buffer.String()
	assert.Contains(t, logs, "Expected error in update auctions")
	assert.Contains(t, logs, "is not allowed in the context of an update")
}

func TestClosingExpiredAuctionsWithMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	os.Setenv("BATCH_INSERT_INTERVAL", "100ms")
	defer os.Unsetenv("BATCH_INSERT_INTERVAL")

	os.Setenv("MAX_BATCH_SIZE", "4")
	defer os.Unsetenv("MAX_BATCH_SIZE")

	os.Setenv("AUCTION_INTERVAL", "300s")
	defer os.Unsetenv("AUCTION_INTERVAL")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@localhost:27016/auctions?authSource=admin"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("auctions")
	defer db.Drop(ctx)

	collection := db.Collection("auctions")
	_, err = collection.InsertMany(ctx, []interface{}{
		bson.M{
			"_id":          "expired-auction-id",
			"product_name": "Expired Product",
			"category":     "expired-auction-category",
			"description":  "expired-auction-description",
			"condition":    1,
			"status":       0,
			"timestamp":    time.Now().Add(-10 * time.Minute).Unix(),
		},
		bson.M{
			"_id":          "valid-auction-id",
			"product_name": "Valid Product",
			"category":     "valid-auction-category",
			"description":  "expired-auction-description",
			"condition":    1,
			"status":       0,
			"timestamp":    time.Now().Add(10 * time.Minute).Unix(),
		},
	})
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	auctionRepo := &AuctionRepository{
		Collection: collection,
	}

	go auctionRepo.ClosingExpiredAuctions(ctx)

	time.Sleep(3 * time.Second)

	var expiredAuction bson.M
	err = collection.FindOne(ctx, bson.M{"_id": "expired-auction-id"}).Decode(&expiredAuction)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), expiredAuction["status"])

	var validAuction bson.M
	err = collection.FindOne(ctx, bson.M{"_id": "valid-auction-id"}).Decode(&validAuction)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), validAuction["status"])
}

func TestClosingExpiredAuctionsDecodingErrorWithMongo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	os.Setenv("BATCH_INSERT_INTERVAL", "100ms")
	defer os.Unsetenv("BATCH_INSERT_INTERVAL")

	os.Setenv("MAX_BATCH_SIZE", "4")
	defer os.Unsetenv("MAX_BATCH_SIZE")

	os.Setenv("AUCTION_INTERVAL", "300s")
	defer os.Unsetenv("AUCTION_INTERVAL")

	buffer, restoreLogger := setupTestLoggerForClass()
	defer restoreLogger()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("auctions")
	defer db.Drop(ctx)

	collection := db.Collection("auctions")

	_, err = collection.InsertOne(ctx, bson.M{
		"_id":         "invalid-auction-id",
		"timestamp":   "not-a-timestamp",
		"status":      0,
		"condition":   "not-a-integer",
		"extra_field": "unexpected_data",
	})
	if err != nil {
		t.Fatalf("Failed to insert invalid data: %v", err)
	}

	auctionRepo := &AuctionRepository{
		Collection: collection,
	}

	go auctionRepo.ClosingExpiredAuctions(ctx)

	time.Sleep(time.Second * 1)

	logs := buffer.String()
	assert.Contains(t, logs, "Error decoding auctions")
}
