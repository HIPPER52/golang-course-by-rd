package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	db *mongo.Database
}

func NewStore(db *mongo.Database) *Store {
	return &Store{db: db}
}

func (s *Store) PutDocument(ctx context.Context, collectionName string, doc any) error {
	collection := s.db.Collection(collectionName)

	bdoc, ok := doc.(map[string]interface{})
	if !ok {
		return errors.New("document must be a map")
	}

	id, exists := bdoc["_id"]
	if !exists {
		return errors.New("document must contain _id field")
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bdoc}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (s *Store) GetDocument(ctx context.Context, collectionName, keyField, keyValue string) (map[string]any, error) {
	collection := s.db.Collection(collectionName)
	filter := bson.M{keyField: keyValue}

	var result map[string]any
	err := collection.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return result, err
}

func (s *Store) DeleteDocument(ctx context.Context, collectionName, keyField, keyValue string) (bool, error) {
	collection := s.db.Collection(collectionName)
	filter := bson.M{keyField: keyValue}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	var result bool = res.DeletedCount > 0
	return result, nil
}

func (s *Store) ListDocuments(ctx context.Context, collectionName string) ([]map[string]any, error) {
	collection := s.db.Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]any
	for cursor.Next(ctx) {
		var doc map[string]any
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		results = append(results, doc)
	}
	return results, nil
}

func (s *Store) CreateIndex(ctx context.Context, collectionName, field string, unique bool) error {
	collection := s.db.Collection(collectionName)

	indexName := "ix_" + field

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}},
		Options: options.Index().SetUnique(unique).SetName(indexName),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (s *Store) DeleteIndex(ctx context.Context, collectionName, indexName string) error {
	collection := s.db.Collection(collectionName)
	_, err := collection.Indexes().DropOne(ctx, indexName)
	return err
}

func (s *Store) CreateCollection(ctx context.Context, name string) error {
	err := s.db.CreateCollection(ctx, name)
	if err != nil {
		return err
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("ix_id"),
	}

	_, err = s.db.Collection(name).Indexes().CreateOne(ctx, indexModel)
	return err
}

func (s *Store) DeleteCollection(ctx context.Context, name string) error {
	return s.db.Collection(name).Drop(ctx)
}

func (s *Store) ListCollections(ctx context.Context) ([]string, error) {
	names, err := s.db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	return names, nil
}
