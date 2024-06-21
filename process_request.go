/*
 * Copyright (C) 2024 Key9 Identity, Inc <k9.io>
 * Copyright (C) 2024 Champ Clark III <cclark@k9.io>
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
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Process_Key9(c *gin.Context) {

	var cache_body string
	var cache_ret bool

	var req *http.Request
	var err error
	var jsondata []uint8

	sha1_file := fmt.Sprintf("%s/%x", Config.Proxy.Cache_Dir, sha1.Sum([]byte(c.Request.URL.Path)))

	client := http.Client{}

	api_key_temp := fmt.Sprintf("%s:%s", c.GetString("company_uuid"), c.GetString("api_key"))
	url_tmp := fmt.Sprintf("%s%s", Config.Core.Address, c.Request.URL.Path)

	log.Printf("Proxied request: %s\n", url_tmp) /* Display/log the URL proxied */

	/* Determine if this is a GET or POST request.  POST is used to pass the "remote"
	   IP addresses.   This is used with Geolocks */

	if c.Request.Method == "GET" {

		req, err = http.NewRequest("GET", url_tmp, nil)

	} else {

		jsondata, _ = c.GetRawData()
		req, err = http.NewRequest("POST", url_tmp, bytes.NewBuffer(jsondata))

	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Can't pull data from proxy cache [http.NewRequest]"})
		c.Abort()
		return
	}

	req.Header["API_KEY"] = []string{api_key_temp}

	res, err := client.Do(req)

	/* If the requests fails,  we pull data from cache.  'Fail' means a connection cannot
	   be established to the webserver,  not an authentication issue */

	if err != nil {

		/* Request has failed,  pull API key from cache and validate and authenticate */

		if Cache_Authenticate_API(api_key_temp) == true {

			log.Printf("Pulling authentication from cache.\n")

			cache_body, cache_ret = Read_Cache(sha1_file)

			/* We pull the API keys from cache.  This means that in the past, the proxy has
			   to have made a valid/authenticated connection to the Key9 API!  If it hasn't
			   ( cache_ret == false ),  there is nothing we can do.  */

			if cache_ret == true {

				c.Data(http.StatusOK, "application/json", []byte(cache_body))
				log.Printf("Error pulling from Key9 API [client.Do],  using cache %s\n", sha1_file)
				return

			} else {

				c.JSON(http.StatusOK, gin.H{"error": "Can't pull data from proxy cache [client.Do]"})
				c.Abort()
				return

			}

		} else {

			/* Authentication doesn't match up from cache and the client */

			c.JSON(http.StatusOK, gin.H{"error": "api authentication failure [Key9 Proxy]"})
			c.Abort()
			return
		}

	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {

		/* If the connect was establish but it isn't possible to read the body,  we pull from cache
		   data */

		if Cache_Authenticate_API(api_key_temp) == true {

			log.Printf("Pulling authentication from cache.\n")

			cache_body, cache_ret = Read_Cache(sha1_file)

			if cache_ret == true {

				c.Data(http.StatusOK, "application/json", []byte(cache_body))
				log.Printf("Error pulling from Key9 API [client.Do],  using cache %s\n", sha1_file)
				return

			} else {

				c.JSON(http.StatusOK, gin.H{"error": "Can't pull data from proxy cache [client.Do]"})
				c.Abort()
				return

			}

		} else {

			c.JSON(http.StatusOK, gin.H{"error": "api authentication failure [Key9 Proxy]"})
			c.Abort()
			return
		}

	}

	/* Cache API creds if there isn't an error */

	if strings.Contains(string(body), "\"error\":") == false {

		api_file_temp := fmt.Sprintf("%s/api.cache", Config.Proxy.Cache_Dir)
		Write_Cache(api_file_temp, api_key_temp)
	}

	/* Send data to the client */

	c.Data(http.StatusOK, "application/json", body)

	/* Write cache of whatever we have looked up */

	Write_Cache(sha1_file, string(body))

}
