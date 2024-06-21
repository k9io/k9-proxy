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
	"log"
	"os/user"
	"strconv"
	"syscall"
)

/*********************************************************************************************/
/* DropPrivileges - When started by "root",  this "drops" the privs to a lesser user         */
/*                                                                                           */
/* This is from https://stackoverflow.com/questions/41248866/golang-dropping-privileges-v1-7 */
/*********************************************************************************************/

func DropPrivileges(userToSwitchTo string) {

	// Lookup user and group IDs for the user we want to switch to.

	userInfo, err := user.Lookup(userToSwitchTo)
	if err != nil {
		log.Fatal(err)
	}

	// Convert group ID and user ID from string to int.

	gid, err := strconv.Atoi(userInfo.Gid)
	if err != nil {
		log.Fatal(err)
	}

	uid, err := strconv.Atoi(userInfo.Uid)
	if err != nil {
		log.Fatal(err)
	}

	// Unset supplementary group IDs.

	err = syscall.Setgroups([]int{})
	if err != nil {
		log.Fatal("Failed to unset supplementary group IDs: " + err.Error())
	}

	// Set group ID (real and effective).

	err = syscall.Setgid(gid)
	if err != nil {
		log.Fatal("Failed to set group ID: " + err.Error())
	}

	// Set user ID (real and effective).

	err = syscall.Setuid(uid)
	if err != nil {
		log.Fatal("Failed to set user ID: " + err.Error())
	}

}
