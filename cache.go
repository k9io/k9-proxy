package main

import (
	"bytes"
	"log"
	"os"
	"strings"
)

/************************************************************************/
/* Write_Cache - Write's a "cache" file.  "cache" files are use when we */
/* can't connect to the Key9 API.                                       */
/************************************************************************/

func Write_Cache(sha1_file string, body string) {

	/* Cheap (CPU wise) check for JSON at error.   We don't "cache"
	   errors */

	if strings.Contains(body, "\"error\":") {

		log.Printf("Key9 API returned an error, not caching: %s\n", string(body) )
		return

	}

	/* Have we ever "cached" this query before? */

	_, err_stat := os.Stat(sha1_file)

	/* No cache, create one */

	if err_stat != nil {

		log.Printf("No cache file, creating %s\n", sha1_file)

		err_write := os.WriteFile(sha1_file, []byte(body), 0600)

		if err_write != nil {

			log.Printf("Can't create %s\n", sha1_file)
			return

		}
	}

	/* We stat above, so we ignore the bool return from Read_Cache */

	current_cache, _ := Read_Cache(sha1_file)

	/* Check to see if the API/cache values have changed */

	if bytes.Equal([]byte(current_cache[:]), []byte(body[:])) == false {

		log.Printf("Cache %s is stale, updating.\n", sha1_file)

		err := os.WriteFile(sha1_file, []byte(body), 0600)

		if err != nil {

			log.Printf("Cannot update %s.\n", sha1_file)
			return

		}
	}
}

/****************************************************************/
/* Read_Cache - If Key9 APIs call fails, we pull from "cache".  */
/****************************************************************/

func Read_Cache(filename string) (string, bool) {

	_, err := os.Stat(filename)

	if err != nil {

		log.Printf("Cannot find cache file %s.\n", filename)
		return "", false

	}

	body, err := os.ReadFile(filename)

	if err != nil {

		log.Printf("Cannot read %s.\n", filename)
		return "", false

	}

	return (string(body)), true

}
