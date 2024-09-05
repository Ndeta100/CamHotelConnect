package db

import (
	"context"
	"fmt"
	"github.com/Ndeta100/CamHotelConnect/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

const userColl = "users"

type Dropper interface {
	Drop(ctx context.Context) error
}

// AuthStore This field are for user auth
type AuthStore interface {
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
}

type UserStore interface {
	Dropper
	AuthStore
	GetUserByID(ctx context.Context, s primitive.ObjectID) (*types.User, error)
	GetUsers(ctx context.Context) ([]*types.User, error)
	InsertUser(ctx context.Context, user *types.User) (*types.User, error)
	DeleteUser(ctx context.Context, s string) error
	UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	UserExists(ctx context.Context, user *types.User) (bool, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	dbname := os.Getenv("MONGO_DB_NAME")
	coll := client.Database(dbname).Collection(userColl)
	return &MongoUserStore{
		client: client,
		coll:   coll,
	}

}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("-----dropping user collection-----")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	exist, err := s.UserExists(ctx, user)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("user already exists")
	}
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	//TODO: maybe a good idea to handle when no user is deleted
	//check the code below
	//res,err:=s.coll.DeleteOne(ctx,bson.M{"_id":oid})
	//if res.DeletedCount==0{}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {
	//update := bson.D{
	//	{
	//		"$set", params.ToBSON(),
	//	},
	//}
	update := bson.M{
		"$set": params,
	}
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id primitive.ObjectID) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) UserExists(ctx context.Context, user *types.User) (bool, error) {
	filter := bson.M{"email": user.Email}
	count, err := s.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
