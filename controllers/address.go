package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mgsquare/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid code"})
			ctx.Abort()
			return
		}
		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(500, "Internal Server Error")
		}
		var addresses models.Address

		addresses.Address_ID = primitive.NewObjectID()
		if err = ctx.BindJSON(&addresses); err != nil {
			ctx.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointCursor, err := UserCollection.Aggregate(context, mongo.Pipeline{match_filter, unwind, group})

		if err != nil {
			ctx.IndentedJSON(500, "Internal server error")

		}
		var addressInfo []bson.M
		if err := pointCursor.All(context, &addressInfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := UserCollection.UpdateOne(context, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			ctx.IndentedJSON(400, "Not Allowed")
		}
		defer cancel()
		context.Done()

	}
}
func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id_hex := c.Query("id")
		if user_id_hex == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}
		user_id, err := primitive.ObjectIDFromHex(user_id_hex)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "address.0.house_name", Value: editaddress.House},
				{Key: "address.0.street_name", Value: editaddress.Street},
				{Key: "address.0.city_name", Value: editaddress.City},
				{Key: "address.0.pin_code", Value: editaddress.Pincode},
			}},
		}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something went wrong")
			return
		}

		c.IndentedJSON(200, "Successfully updated the home address")

	}
}
func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id_hex := c.Query("id")
		if user_id_hex == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}
		user_id, err := primitive.ObjectIDFromHex(user_id_hex)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "address.1.house_name", Value: editaddress.House},
				{Key: "address.1.street_name", Value: editaddress.Street},
				{Key: "address.1.city_name", Value: editaddress.City},
				{Key: "address.1.pin_code", Value: editaddress.Pincode},
			}},
		}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something went wrong")
			return
		}

		c.IndentedJSON(200, "Successfully updated the work address")

	}
}
func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id_hex := c.Query("id")
		if user_id_hex == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		emptyAddress := make([]models.Address, 0)
		user_id, err := primitive.ObjectIDFromHex(user_id_hex)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresses", Value: emptyAddress}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "Address couldn't be deleted")
			return
		}

		c.IndentedJSON(200, "Successfully deleted")

	}
}
