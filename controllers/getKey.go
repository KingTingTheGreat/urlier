package controllers

import (
	"fmt"
	"net/http"
	"time"
	"urlier/configs"
	"urlier/models"
	"urlier/templates"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var entryCollection2 *mongo.Collection = configs.GetCollection(configs.DB, "Entries")

func GetKey(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := c.Param("key")

	var entry models.Entry
	err := entryCollection2.FindOne(ctx, bson.M{"key": key}).Decode(&entry)
	if err != nil {
		fmt.Println("Key not found:", key)
		fmt.Println(entry)
		return renderPage(c, http.StatusNotFound, templates.KeyNotFound(key), "Not Found")
	}

	// increment the number of visits
	entry.Visits++
	_, err = entryCollection2.UpdateOne(ctx, bson.M{"key": key}, bson.M{"$set": bson.M{"visits": entry.Visits}})
	if err != nil {
		fmt.Println("Failed to update visits:", err)
	}

	fmt.Println("Redirecting to", entry.Url)
	return c.Redirect(http.StatusTemporaryRedirect, entry.Url)
}
