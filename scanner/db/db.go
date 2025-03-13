package db

import (
	"context"
	"time"

	"github.com/back2basic/siaalert/scanner/logger"
	stypes "go.sia.tech/core/types"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB struct holds the connection and collection details
type MongoDB struct {
	Client   *mongo.Client
	ColHosts *mongo.Collection
	ColScan  *mongo.Collection
	ColApi   *mongo.Collection
	ColAlert *mongo.Collection
	ColRhp   *mongo.Collection
}

// NewMongoDB initializes a new MongoDB instance
func NewMongoDB(uri, dbName, collectionHost, collectionScan, collectionApi, collectionAlert, collectionRhp string) (*MongoDB, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		Client:   client,
		ColHosts: client.Database(dbName).Collection(collectionHost),
		ColScan:  client.Database(dbName).Collection(collectionScan),
		ColApi:   client.Database(dbName).Collection(collectionApi),
		ColAlert: client.Database(dbName).Collection(collectionAlert),
		ColRhp:   client.Database(dbName).Collection(collectionRhp),
	}, nil
}

// Close disconnects the MongoDB client
func (db *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.Client.Disconnect(ctx)
}

// UpdateHosts inserts or updates multiple hosts document
func (db *MongoDB) UpdateHosts(hosts []bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var interfaceHosts []interface{}
	for _, host := range hosts {
		interfaceHosts = append(interfaceHosts, host)
	}

	_, err := db.ColHosts.InsertMany(ctx, interfaceHosts)
	return err
}

// UpdateHost inserts or updates a host document
func (db *MongoDB) UpdateHost(hostID stypes.PublicKey, hostData bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	log := logger.GetLogger()
	defer log.Sync()

	filter := bson.M{"publicKey": hostID.String()}
	update := bson.M{"$set": hostData}
	result, err := db.ColHosts.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if result.MatchedCount == 0 {
		log.Info("UpdateHost: no match", zap.String("hostID", hostID.String()))
		return err
	}
	if result.ModifiedCount == 0 {
		// log.Info("UpdateHost: no change", zap.String("hostID", hostID.String()))
		return err
	}
	return err
}

// FindHosts retrieves hosts based on a filter
func (db *MongoDB) FindHosts(filter bson.M) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.ColHosts.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	for cursor.Next(ctx) {
		var host bson.M
		if err := cursor.Decode(&host); err != nil {
			return nil, err
		}
		results = append(results, host)
	}
	return results, nil
}

func (db *MongoDB) FindScan(filter bson.M) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.ColScan.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	for cursor.Next(ctx) {
		var scan bson.M
		if err := cursor.Decode(&scan); err != nil {
			return nil, err
		}
		results = append(results, scan)
	}
	return results, nil
}

func (db *MongoDB) InsertScan(data bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// filter := bson.M{"publicKey": hostID.String()}
	// insert := bson.M{"$set": data}
	_, err := db.ColScan.InsertOne(ctx, data, options.InsertOne())
	return err
}

func (db *MongoDB) FindRhp(filter bson.M) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.ColRhp.FindOne(ctx, filter)
}

func (db *MongoDB) FindRhps(filter bson.M) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.ColRhp.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	for cursor.Next(ctx) {
		var rhp bson.M
		if err := cursor.Decode(&rhp); err != nil {
			return nil, err
		}
		results = append(results, rhp)
	}
	return results, nil
}

// UpdateRhp inserts or updates a rhp document
func (db *MongoDB) UpdateRhp(publicKey string, online bool, data bson.M, log *zap.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"publicKey": publicKey}
	update := bson.M{"$set": data}
	// log.Info("Attempting update", zap.Any("filter", filter), zap.Any("update", update))

	_, err := db.ColRhp.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	// log.Info("update rhp", zap.Bool("success", docnew.ModifiedCount > 0), zap.Error(err))
	return err
}

// UpdateOtp inserts or updates a otp document
func (db *MongoDB) UpdateOtp(publicKey, email, exp, secret string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"publicKey": publicKey}
	update := bson.M{"$set": bson.M{"secret": secret, "expire": exp, "email": email}}

	_, err := db.ColApi.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (db *MongoDB) FindOtp(filter bson.M) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.ColApi.FindOne(ctx, filter)
}

func (db *MongoDB) DeleteOtp(publicKey, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"publicKey": publicKey, "email": email}
	_, err := db.ColApi.DeleteOne(ctx, filter)
	return err
}

func (db *MongoDB) UpdateAlert(publicKey string, data bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"publicKey": publicKey}
	update := bson.M{"$set": data}
	_, err := db.ColAlert.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (db *MongoDB) FindAlerts(filter bson.M) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.ColAlert.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	for cursor.Next(ctx) {
		var alert bson.M
		if err := cursor.Decode(&alert); err != nil {
			return nil, err
		}
		results = append(results, alert)
	}
	return results, nil
}

func (db *MongoDB) DeleteAlert(publicKey, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"publicKey": publicKey, "email": email}
	_, err := db.ColAlert.DeleteOne(ctx, filter)
	return err
}

// GetSubscribers returns the subscribers of a host
func (db *MongoDB) GetSubscribers(publicKey string, log *zap.Logger) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.ColAlert.Find(ctx, bson.M{"publicKey": publicKey})
	if err != nil {
		log.Warn("Failed to find document", zap.Error(err))
		return nil
	}
	defer cursor.Close(ctx)

	var subscribers []string
	for cursor.Next(ctx) {
		var alert bson.M
		if err := cursor.Decode(&alert); err != nil {
			log.Error("Failed to decode document", zap.Error(err))
			return nil
		}
		if alert["type"] != "email" {
			continue
		}
		subscribers = append(subscribers, alert["sender"].(string))
	}

	// db := GetAppwriteDatabaseService()
	// cfg := config.GetConfig()
	// found, err := db.ListDocuments(
	// 	cfg.Appwrite.Database.Id,
	// 	cfg.Appwrite.ColAlert.Id,
	// 	[]string{
	// 		query.Equal("hostId", hostId),
	// 	},
	// )

	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil
	// }
	// var alertList strict.AlertList
	// found.Decode(&alertList)

	// // loop t and find email subscribers
	// var subscribers []string
	// for _, doc := range alertList.Documents {
	// 	if doc.Type != "email" {
	// 		continue
	// 	}
	// 	subscribers = append(subscribers, doc.Sender)
	// }

	return subscribers
}
