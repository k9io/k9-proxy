package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

/*******************************************************/
/* HTTP_Logger - Used to log HTTP requests to k9-proxy */
/*******************************************************/

func HTTP_Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientIP := c.ClientIP()
		now := time.Now()
		log.Printf("[%s] %s %s %s", now.Format(time.RFC3339), c.Request.Method, c.Request.URL.Path, clientIP)
		c.Next()
	}
}
