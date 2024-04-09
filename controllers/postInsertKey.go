package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"urlier/configs"
	"urlier/models"
	"urlier/responses"
	"urlier/templates"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var entryCollection *mongo.Collection = configs.GetCollection(configs.DB, "Entries")

// PostKey is a controller that renders the key page
func PostInsertKey(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check request Accept header
	// send html if Accept header is text/html; otherwise, send json
	html := c.Request().Header.Get("Accept") == "text/html"

	// get request values
	key := c.FormValue("key")
	url := c.FormValue("url")
	pw := c.FormValue("password")

	if key == "" || url == "" {
		if html {
			return renderComponent(c, http.StatusBadRequest, templates.ErrorMessage("400 Bad Request", "Key and URL are required"))
		}
		return c.JSON(http.StatusBadRequest, responses.Response{
			Status:  http.StatusBadRequest,
			Message: "Key and URL are required",
		})
	}

	// key and url may not match
	if key == url {
		if html {
			return renderComponent(c, http.StatusBadRequest, templates.ErrorMessage("400 Bad Request", "Key and URL must not match"))
		}
		return c.JSON(http.StatusBadRequest, responses.Response{
			Status:  http.StatusBadRequest,
			Message: "Key and URL must not match",
		})
	}

	var entry models.Entry

	// check if the key already exists
	err := entryCollection.FindOne(ctx, bson.M{"key": key}).Decode(&entry)
	if err == nil {
		// key exists and is immutable
		if entry.Password == "" {
			fmt.Println("Key already exists")
			if html {
				return renderComponent(c, http.StatusConflict, templates.ErrorMessage("409 Status Conflit", "Key already exists"))
			}
			return c.JSON(http.StatusConflict, responses.Response{Status: http.StatusConflict, Message: "Key already exists"})
		}
		if !CheckPasswordHash(pw, entry.Password) {
			fmt.Println("Incorrect password")
			if html {
				return renderComponent(c, http.StatusUnauthorized, templates.ErrorMessage("401 Unauthorized", "Incorrect password"))
			}
			return c.JSON(http.StatusUnauthorized, responses.Response{Status: http.StatusUnauthorized, Message: "Incorrect password"})
		}
		// check if password exists and matches
		if entry.Password != "" && CheckPasswordHash(pw, entry.Password) {
			// update the entry
			entry.Url = url
			_, err = entryCollection.ReplaceOne(ctx, bson.M{"key": key, "visits": 0}, entry)
			if err != nil {
				if html {
					return renderComponent(c, http.StatusInternalServerError, templates.InternalServerError())
				}
				return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "Failed to update data"})
			}
			if html {
				return renderComponent(c, http.StatusOK, templates.SuccessMessage("Success", "Successfully updated URL"))
			}
			return c.JSON(http.StatusOK, responses.Response{Message: "Success"})
		}

	}

	// insert this new entry to the database
	var pwHash string
	if pw == "" {
		pwHash = ""
	} else {
		pwHash, err = HashPassword(pw)
		if err != nil {
			if html {
				return renderComponent(c, http.StatusInternalServerError, templates.InternalServerError())
			}
			return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "Failed to hash password"})
		}
	}
	entry = models.Entry{Key: key, Url: url, Password: pwHash, Visits: 0}
	_, err = entryCollection.InsertOne(ctx, entry)
	if err != nil {
		if html {
			return renderComponent(c, http.StatusInternalServerError, templates.InternalServerError())
		}
		return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "Failed to insert data"})
	}

	fmt.Println("Successfully wrote to database")
	if html {
		return renderComponent(c, http.StatusOK, templates.SuccessMessage("Success", "Successfully inserted data"))
	}
	return c.JSON(http.StatusOK, responses.Response{Message: "Success"})
}
