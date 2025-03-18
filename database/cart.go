package database

import (
	"context"
	"errors"
	"log"

	"github.com/mgsquare/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct      = errors.New("can't find the product")
	ErrCantDecodeProducts   = errors.New("can't find the product")
	ErrUserIdIsNotValidated = errors.New("this user is not valid")
	ErrCantUpdateUser       = errors.New("cannot add this item from the cart")
	ErrCantRemoveItem       = errors.New("cannot remove this item from the cart")
	ErrCantGetItem          = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem      = errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromDB, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser
	err = searchFromDB.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValidated
	}

	filter := bson.D{primitive.E{Key: "$push", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValidated
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItem
	}
	return nil

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
