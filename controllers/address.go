package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mgsquare/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {

}
func EditHomeAddress() gin.HandlerFunc {

}
func EditWorkAddress() gin.HandlerFunc {

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
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully deleted")

	}
}
