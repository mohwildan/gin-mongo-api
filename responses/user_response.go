package responses

import "github.com/gin-gonic/gin"

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}
