package model

import ()

type Roles int64

var (
	ROLE_LOGIN Roles = 1
	// WITH_NONE  Roles = 0

	ROLE_Secret Roles = 1 << 1 // Preserve, used for test.
	ROLE_Sales  Roles = 1 << 2
	// WITH_PERSON  Roles = 1 << 1 // customer or factory
)

// now use strings implementation.
// var (
// 	ROLE_RequireLogin = "require_login"
// )
