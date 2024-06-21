package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/***************************************************************************/
/* Authenticate_API - The pulled API_KEY data from the proxy "client" side */
/***************************************************************************/

func Authenticate_API() gin.HandlerFunc {
	return func(c *gin.Context) {

		api_key := c.GetHeader("API_KEY")

		if api_key == "" {
			c.JSON(http.StatusOK, gin.H{"error": "api authentication failed"})
			c.Abort()
			return
		}

		temp_value := strings.Split(api_key, ":")

		/* Validate the string properly split */

		if len(temp_value) != 2 {
			c.JSON(http.StatusOK, gin.H{"error": "api authentication failed"})
			c.Abort()
			return
		}

		c.Set("company_uuid", temp_value[0])
		c.Set("api_key", temp_value[1])

	}
}

/******************************************************************************/
/* Cache_Authenticate_API - When we can't connect to the actual Key9 API,  we */
/* use the "cached" API creds                                                 */
/******************************************************************************/

func Cache_Authenticate_API(api_key string) bool {

	temp_api_file := fmt.Sprintf("%s/api.cache", Config.Proxy.Cache_Dir)

	api_cache_data, ret := Read_Cache(temp_api_file)

	if ret == false {
		return false
	}

	if bytes.Equal([]byte(api_key[:]), []byte(api_cache_data[:])) {
		return true
	}

	return false
}
