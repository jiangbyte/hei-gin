package user

import "hei-gin/sdk/utils"

// ClientUserToClientUserVO 将 user.ClientUser 映射到 user.ClientUserVO
func ClientUserToClientUserVO(src *ClientUser) *ClientUserVO {
	if src == nil {
		return nil
	}

	dst := &ClientUserVO{}

	dst.ID = src.ID
	dst.Username = src.Username
	dst.Nickname = src.Nickname
	dst.Avatar = src.Avatar
	dst.Motto = src.Motto
	dst.Gender = src.Gender
	dst.Email = src.Email
	dst.Github = src.Github
	dst.Phone = src.Phone
	dst.Status = src.Status
	dst.LastLoginIP = src.LastLoginIP
	dst.LoginCount = src.LoginCount

	// *time.Time → string manual conversion
	dst.LastLoginAt = utils.FormatDateTimePtr(src.LastLoginAt)
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)

	return dst
}

// ClientUserVOToClientUser 将 user.ClientUserVO 映射到 user.ClientUser
func ClientUserVOToClientUser(src *ClientUserVO) *ClientUser {
	if src == nil {
		return nil
	}

	dst := &ClientUser{}

	dst.ID = src.ID
	dst.Username = src.Username
	dst.Nickname = src.Nickname
	dst.Avatar = src.Avatar
	dst.Motto = src.Motto
	dst.Gender = src.Gender
	dst.Email = src.Email
	dst.Github = src.Github
	dst.Phone = src.Phone
	dst.Status = src.Status
	dst.LastLoginIP = src.LastLoginIP
	dst.LoginCount = src.LoginCount

	return dst
}
