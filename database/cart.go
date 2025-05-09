package database

import (
	"context"
	"errors"
	"log"
	"time"

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
	ErrEmptyCart            = errors.New("cannot checkout empty cart")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromDB, err := prodCollection.Find(ctx, bson.M{"_id": productID})

	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productcart []models.ProductUser
	err = searchFromDB.All(ctx, &productcart)

	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	} else if len(productcart) == 0 {
		log.Println("No products found in DB for given ID")
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValidated
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productcart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return ErrCantRemoveItem
	}
	return nil

}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValidated
	}
	var getCartItems models.User
	var orderCart models.Order

	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Ordered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "total", Value: bson.D{
				{Key: "$sum", Value: "$usercart.price"},
			}},
		}},
	}

	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}
	var getUserCart []bson.M

	if err := currentResults.All(ctx, &getUserCart); err != nil {
		panic(err)
	}

	var totalPrice int32

	for _, userItem := range getUserCart {
		price := userItem["total"]
		totalPrice = price.(int32)
	}
	if totalPrice == 0 {
		return ErrEmptyCart
	}
	orderCart.Price = int(totalPrice)

	filter := bson.D{{Key: "_id", Value: id}}

	update := bson.D{{Key: "$push", Value: bson.D{{Key: "orders", Value: orderCart}}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	err = userCollection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)

	if err != nil {
		log.Println(err)
	}

	usercart_empty := make([]models.ProductUser, 0)

	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}

	_, err = userCollection.UpdateOne(ctx, filter3, update3)

	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValidated
	}

	var productDetails models.ProductUser

	var orderDetails models.Order

	orderDetails.Order_ID = primitive.NewObjectID()

	orderDetails.Ordered_At = time.Now()

	orderDetails.Order_Cart = make([]models.ProductUser, 0)

	orderDetails.Payment_Method.COD = true

	err = prodCollection.FindOne(ctx, bson.D{{Key: "_id", Value: productID}}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
	}

	orderDetails.Price = productDetails.Price

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "orders", Value: orderDetails}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": productDetails}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)

	if err != nil {
		log.Println(err)
	}

	return nil

}
