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

	log.Printf("* Using configuration file - %s\n", Config_File)

	Config := LoadConfig(Config_File)

	DropPrivileges(Config.Authentication.Runas)

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

	log.Printf("Listening traffic on %s.", Config.Proxy.HTTP_Listen)

	/* Non-TLS */

	if Config.Proxy.HTTP_TLS == false { 

		err = router.Run(Config.Proxy.HTTP_Listen)
	
	} else {

	/* TLS */

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
