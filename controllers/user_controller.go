package controllers

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.Db, "Users")
var validate = validator.New()

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cencel()
		var user models.User

		//validate request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Error",
				Data: gin.H{
					"data": err.Error(),
				},
			})
			return
		}
		//user the validator library to validete required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: validationErr.Error(),
				Data: gin.H{
					"data": validationErr.Error(),
				},
			})
			return
		}
		//we create new user and enter data this model user
		newUser := models.User{
			Id:       primitive.NewObjectID(),
			Name:     user.Name,
			Location: user.Location,
			Title:    user.Title,
		}

		resault, err := userCollection.InsertOne(ctx, newUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: gin.H{
					"data": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data: gin.H{
				"data": resault,
			},
		})
		return
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cencel()

		var users []models.User

		resault, err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data: gin.H{
					"data": err.Error(),
				},
			})
		}
		defer resault.Close(ctx)
		for resault.Next(ctx) {
			var signleUser models.User
			if err = resault.Decode(&signleUser); err != nil {
				c.JSON(http.StatusInternalServerError, responses.UserResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data: gin.H{
						"data": err.Error(),
					},
				})
			}
			users = append(users, signleUser)
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: gin.H{
				"data": users,
			},
		})
	}
}

func GetaUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var user models.User
		defer cencel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.UserResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": "user tidak di temukan"},
				},
			)
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: gin.H{
				"data": user,
			},
		})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cencel()
		var user models.User

		objId, _ := primitive.ObjectIDFromHex(userId)

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusOK,
				Message: "error",
				Data: gin.H{
					"data": err.Error(),
				},
			})
		}

		if validateErr := validate.Struct(&user); validateErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data: gin.H{
					"data": validateErr.Error(),
				},
			})
		}

		update := gin.H{
			"name":     user.Name,
			"location": user.Location,
			"title":    user.Title,
		}

		resault, err := userCollection.UpdateOne(ctx, gin.H{"_id": objId}, gin.H{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: gin.H{
					"data": err.Error(),
				},
			})
			return
		}

		var updateUser models.User

		if resault.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, gin.H{"_id": objId}).Decode(&updateUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.UserResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data: gin.H{
						"data": err.Error(),
					},
				})
				return
			}
		}
		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: gin.H{
				"data": updateUser,
			},
		})

	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cencel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		resault, err := userCollection.DeleteOne(ctx, gin.H{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: gin.H{
					"data": err.Error(),
				},
			})
		}
		if resault.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.UserResponse{
				Status:  http.StatusNotFound,
				Message: "error",
				Data: gin.H{
					"data": "User with specified ID not found!",
				},
			})
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: gin.H{
				"data": "user success delete",
			},
		})
	}
}
