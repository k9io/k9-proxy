package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

/**********************/
/* Configuration Data */
/**********************/

type Configuration struct {
	Authentication struct {
		Address string `yaml:"address"`
		Runas   string `yaml:"runas"`
	}

	Proxy struct {
		HTTP_Listen string `yaml:"http_listen"`
		HTTP_Mode   string `yaml:"http_mode"`
		HTTP_TLS    bool   `yaml:"http_tls"`
		HTTP_Cert   string `yaml:"http_cert"`
		HTTP_Key    string `yaml:"http_key"`
		Logs        string `yaml:"logs"`
		Cache_Dir   string `yaml:"cache_dir"`
	}
}

var Config *Configuration

/****************************************************/
/* LoadConfig - Doh.. what do you think it does? :) */
/****************************************************/

func LoadConfig(config_file string) *Configuration {

	/* Load config file */

	file, err := os.Open(config_file)

	if err != nil {
		log.Fatalf("Cannot open '%s' YAML file.", config_file)
	}

	defer file.Close()

	/* Init new YAML decode */

	d := yaml.NewDecoder(file)

	err = d.Decode(&Config)

	if err != nil {
		log.Fatalf("Cannot parse '%s'.", config_file)
	}

	/* ---- Sanity Checks ---- Authentication ---- */

	if Config.Authentication.Address == "" {
		log.Fatalf("'address' key not found in %s.\n", config_file)
	}

	if Config.Authentication.Runas == "" {
		log.Fatalf("'runas' key not found in %s.\n", config_file)
	}

	/* ---- Proxy ---- */

	if Config.Proxy.HTTP_Listen == "" {
		log.Fatalf("'http_listen' key not found in %s.\n", config_file)
	}

	/* Only check if TLS is enabled, otherwise,  we don't care */

	if Config.Proxy.HTTP_TLS == true {

		if Config.Proxy.HTTP_Cert == "" {
			log.Fatalf("'http_cert' key not found in %s.\n", config_file)
		}

		if Config.Proxy.HTTP_Cert == "" {
			log.Fatalf("'http_key' key not found in %s.\n", config_file)
		}

	}

	if Config.Proxy.HTTP_Mode == "" {
		log.Fatalf("'http_mode' key not found in %s.\n", config_file)
	}

        if Config.Proxy.HTTP_Mode != "release" && Config.Proxy.HTTP_Mode != "debug" && Config.Proxy.HTTP_Mode != "test" {
                log.Fatalf("Invalid 'http_mode' :  %s.  Valid 'http_modes' are 'release', 'debug' and 'test'\n", Config.Proxy.HTTP_Mode)
	        }

	if Config.Proxy.Logs == "" {
		log.Fatalf("'logs' key not found in %s.\n", config_file)
	}

	if Config.Proxy.Cache_Dir == "" {
		log.Fatalf("'cache_dir' key not found in %s.\n", config_file)
	}

	return Config
}
