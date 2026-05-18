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
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// serveFromCache validates the API key against the local cache and, if valid,
// returns the cached response body. Used when the upstream Key9 API is unreachable.
func serveFromCache(c *gin.Context, apiKey, cacheFile string) {
	if !Cache_Authenticate_API(apiKey) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "api authentication failure [Key9 Proxy]"})
		c.Abort()
		return
	}
	log.Printf("Pulling authentication from cache.\n")
	cacheBody, ok := Read_Cache(cacheFile)
	if !ok {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Can't pull data from proxy cache"})
		c.Abort()
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(cacheBody))
	log.Printf("Key9 API unavailable, serving cache %s\n", cacheFile)
}

func Process_Key9(c *gin.Context) {

	cacheFile := fmt.Sprintf("%s/%x", Config.Proxy.Cache_Dir, sha256.Sum256([]byte(c.Request.URL.RequestURI())))

	client := http.Client{Timeout: time.Duration(Config.Core.Connection_Timeout) * time.Second}

	apiKey := fmt.Sprintf("%s:%s", c.GetString("company_uuid"), c.GetString("api_key"))
	urlTmp := fmt.Sprintf("%s%s", Config.Core.Address, c.Request.URL.RequestURI())

	log.Printf("Proxied request: %s\n", urlTmp)

	var req *http.Request
	var err error

	if c.Request.Method == "GET" {
		req, err = http.NewRequest("GET", urlTmp, nil)
	} else {
		var postBody []byte
		postBody, err = c.GetRawData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
			c.Abort()
			return
		}
		req, err = http.NewRequest("POST", urlTmp, bytes.NewBuffer(postBody))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build upstream request"})
		c.Abort()
		return
	}

	req.Header["API_KEY"] = []string{apiKey}

	res, err := client.Do(req)
	if err != nil {
		serveFromCache(c, apiKey, cacheFile)
		return
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		serveFromCache(c, apiKey, cacheFile)
		return
	}

	if !strings.Contains(string(responseBody), "\"error\":") {
		Write_Cache(fmt.Sprintf("%s/api.cache", Config.Proxy.Cache_Dir), apiKey)
	}

	c.Data(http.StatusOK, "application/json", responseBody)
	Write_Cache(cacheFile, string(responseBody))
}
