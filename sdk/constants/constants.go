package constants

import "hei-gin/sdk/enums"

// SUPER_ADMIN_CODE is the role code for the built-in super administrator.
const SUPER_ADMIN_CODE = "SUPER_ADMIN"

// Permission
const PERMISSION_CACHE_KEY = "hei:permission:keys"

// Auth token / session Redis keys
var (
	TOKEN_PREFIX_BUSINESS   = "hei:auth:" + string(enums.LoginTypeBusiness) + ":token:"
	SESSION_PREFIX_BUSINESS = "hei:auth:" + string(enums.LoginTypeBusiness) + ":session:"

	TOKEN_PREFIX_CONSUMER   = "hei:auth:" + string(enums.LoginTypeConsumer) + ":token:"
	SESSION_PREFIX_CONSUMER = "hei:auth:" + string(enums.LoginTypeConsumer) + ":session:"
)

// No-repeat request prevention Redis key prefix
const NO_REPEAT_PREFIX = "norepeat:"
