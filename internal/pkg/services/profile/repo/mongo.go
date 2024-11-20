package repo

import (
	"context"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/db"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const usersCollection = "users"

type MongoParams struct {
	fx.In

	Log     *zap.Logger
	Wrapper *db.MongoWrapper
}

type MongoProfileRepo struct {
	collection *mongo.Collection
	log        *zap.Logger
}

func NewMongoProfileRepo(p MongoParams) *MongoProfileRepo {
	collection := p.Wrapper.Client.Database(p.Wrapper.DB).Collection(usersCollection)
	return &MongoProfileRepo{
		collection: collection,
		log:        p.Log,
	}
}

func (pr *MongoProfileRepo) Create(ctx context.Context, user *model.User) error {
	_, err := pr.collection.InsertOne(ctx, user)
	return err
}

func (pr *MongoProfileRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := pr.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	return &user, err
}

func (pr *MongoProfileRepo) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	err := pr.collection.FindOne(ctx, bson.M{"login": login}).Decode(&user)
	return &user, err
}

func (pr *MongoProfileRepo) Update(ctx context.Context, user *model.User) error {
	result, err := pr.collection.UpdateOne(
		ctx,
		bson.M{"id": user.ID},
		bson.D{{"$set", bson.M{
			"name": user.Name,
		}}},
	)
	if err != nil || result.ModifiedCount == 0 {
		return err
	}
	return nil
}
