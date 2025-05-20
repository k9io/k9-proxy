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
