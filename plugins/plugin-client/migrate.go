package plugin_client

import (
	clientUser "hei-gin/plugins/plugin-client/user"
	"hei-gin/sdk/db"
)

func init() {
	db.RegisterModel(&clientUser.ClientUser{})
}
