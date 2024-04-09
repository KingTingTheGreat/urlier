package controllers

import (
	"net/http"
	"sort"
	"time"
	"urlier/configs"
	"urlier/models"
	"urlier/templates"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var entryCollection3 *mongo.Collection = configs.GetCollection(configs.DB, "Entries")

func GetTrending(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get the top 10 entries by number of visits
	var entries []models.Entry
	cursor, err := entryCollection3.Find(context.Background(), bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &entries); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// sort entries by number of visits
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Visits > entries[j].Visits
	})

	// get the top 10 entries
	if len(entries) > 10 {
		entries = entries[:10]
	}

	return renderPage(c, http.StatusOK, templates.Trending(entries), "Trending")
}
