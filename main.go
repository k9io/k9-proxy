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
	"flag"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	var err error
	var Config_File string = "/opt/k9/etc/k9-proxy.yaml"

	config_file := flag.String("config", "", "Configuration file to use")

	flag.Parse()

	if *config_file != "" {
		Config_File = *config_file
	}

	log.Printf("Using configuration file - %s\n", Config_File)

	Config := LoadConfig(Config_File)

	DropPrivileges(Config.Core.Runas)

	log.Printf("Setting gin to \"%s\" mode.\n", Config.Proxy.HTTP_Mode)
	gin.SetMode(Config.Proxy.HTTP_Mode) /* 'debug', 'release' or 'test' */

	router := gin.Default()
	router.Use(HTTP_Logger())

	router.Use(Authenticate_API())

	router.POST("/api/v1/ssh/query/:username/:machine_group", Process_Key9)

	router.GET("/api/v1/query/passwd/username/:username", Process_Key9)
	router.GET("/api/v1/query/passwd/uid/:uid", Process_Key9)
	router.GET("/api/v1/query/passwd/id/:id", Process_Key9)
	router.GET("/api/v1/query/group/gid/:gid", Process_Key9)
	router.GET("/api/v1/query/group/name/:group", Process_Key9)
	router.GET("/api/v1/query/group/id/:id", Process_Key9)
	router.GET("/api/v1/query/k9/all_users", Process_Key9)

	/* Non-TLS */

	if Config.Proxy.HTTP_TLS == false { 

		log.Printf("Listening for unencrypted traffic on %s.", Config.Proxy.HTTP_Listen)
		err = router.Run(Config.Proxy.HTTP_Listen)
	
	} else {

	/* TLS */

		log.Printf("Listening for TLS traffic on %s.", Config.Proxy.HTTP_Listen)
		err = router.RunTLS(Config.Proxy.HTTP_Listen, Config.Proxy.HTTP_Cert, Config.Proxy.HTTP_Key)

	}

	if err != nil {
		
		if Config.Proxy.HTTP_TLS == false {

		log.Fatalf("Cannot bind to %s", Config.Proxy.HTTP_Listen)

		} else { 

		log.Fatalf("Cannot bind it %s or cannot open %s or %s.\n", Config.Proxy.HTTP_Listen, Config.Proxy.HTTP_Cert, Config.Proxy.HTTP_Key)

		}

	}

}
