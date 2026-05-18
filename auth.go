/*
 * Copyright (C) 2024-2025 Key9 Identity, Inc <k9.io>
 * Copyright (C) 2024-2025 Champ Clark III <cclark@k9.io>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License Version 2 as
 * published by the Free Software Foundation.  You may not use, modify or
 * distribute this program under any other version of the GNU General
 * Public License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA 02111-1307, USA.
 */

package main

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/***************************************************************************/
/* Authenticate_API - The pulled API_KEY data from the proxy "client" side */
/***************************************************************************/

func Authenticate_API() gin.HandlerFunc {
	return func(c *gin.Context) {

		api_key := c.GetHeader("API_KEY")

		if api_key == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "api authentication failed"})
			c.Abort()
			return
		}

		temp_value := strings.SplitN(api_key, ":", 2)

		if len(temp_value) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "api authentication failed"})
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

	return subtle.ConstantTimeCompare([]byte(api_key), []byte(api_cache_data)) == 1
}
