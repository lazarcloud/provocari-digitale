package globals

import "time"

var AuthIssuer = "lazar"

var AuthSecretJWTKey = []byte("secret")

var AuthRoles = []string{"public", "service", "loggedIn"}

func AuthIsValidRole(role string) bool {
	for _, r := range AuthRoles {
		if r == role {
			return true
		}
	}
	return false
}

var AuthRolePublic = "public"
var AuthRoleService = "service"
var AuthRoleLoggedIn = "loggedIn"

var AuthJWTTypes = []string{"access", "refresh"}

func AuthIsValidType(t string) bool {
	for _, r := range AuthJWTTypes {
		if r == t {
			return true
		}
	}
	return false
}

var AuthAccessType = "access"
var AuthRefreshType = "refresh"
var AuthAccessTypeDuration = time.Hour * 24
var AuthRefreshTypeDuration = time.Hour * 24 * 30

type ContextCustomKey string

const ContextAccessRoleKey ContextCustomKey = "role"

const ContextUserIdKey ContextCustomKey = "user_id"
