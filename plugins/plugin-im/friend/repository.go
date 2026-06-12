package friend

import (
	"context"

	cliUser "hei-gin/plugins/plugin-client/user"
	imModel "hei-gin/plugins/plugin-im/model"
	sysUser "hei-gin/plugins/plugin-sys/user"
	"hei-gin/sdk/infra/db"
)

func CountFriendship(ctx context.Context, userID, userType, friendID, friendType string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&imModel.Friendship{}).
		Where("(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?) OR "+
			"(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?)",
			userID, userType, friendID, friendType,
			friendID, friendType, userID, userType).
		Count(&count)
	return count
}

func CountPendingRequest(ctx context.Context, senderID, senderType, receiverID, receiverType string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("sender_id = ? AND sender_type = ? AND receiver_id = ? AND receiver_type = ? AND status = ?",
			senderID, senderType, receiverID, receiverType, "pending").
		Count(&count)
	return count
}

func CreateRequest(ctx context.Context, req *imModel.FriendRequest) error {
	return db.DB.WithContext(ctx).Create(req).Error
}

func FindPendingRequestForReceiver(ctx context.Context, requestID, receiverID, receiverType string) (*imModel.FriendRequest, error) {
	var req imModel.FriendRequest
	if err := db.DB.WithContext(ctx).First(&req, "id = ? AND receiver_id = ? AND receiver_type = ? AND status = ?",
		requestID, receiverID, receiverType, "pending").Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func AcceptRequest(ctx context.Context, req *imModel.FriendRequest, pair1, pair2 *imModel.Friendship) error {
	tx := db.DB.WithContext(ctx).Begin()
	if err := tx.Model(req).Update("status", "accepted").Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(pair1).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(pair2).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func RejectRequest(ctx context.Context, requestID, receiverID, receiverType string) (int64, error) {
	res := db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("id = ? AND receiver_id = ? AND receiver_type = ? AND status = ?",
			requestID, receiverID, receiverType, "pending").
		Update("status", "rejected")
	return res.RowsAffected, res.Error
}

func ListFriendships(ctx context.Context, userID, userType string) []imModel.Friendship {
	var rows []imModel.Friendship
	db.DB.WithContext(ctx).Model(&imModel.Friendship{}).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Find(&rows)
	return rows
}

func ListSysUsers(ctx context.Context, ids []string) []sysUser.SysUser {
	var users []sysUser.SysUser
	db.DB.WithContext(ctx).Model(&sysUser.SysUser{}).Where("id IN ?", ids).Find(&users)
	return users
}

func ListClientUsers(ctx context.Context, ids []string) []cliUser.ClientUser {
	var users []cliUser.ClientUser
	db.DB.WithContext(ctx).Model(&cliUser.ClientUser{}).Where("id IN ?", ids).Find(&users)
	return users
}

func ListPendingIncoming(ctx context.Context, userID, userType string) []imModel.FriendRequest {
	var rows []imModel.FriendRequest
	db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("receiver_id = ? AND receiver_type = ? AND status = ?", userID, userType, "pending").
		Order("created_at DESC").Find(&rows)
	return rows
}

func ListPendingOutgoing(ctx context.Context, userID, userType string) []imModel.FriendRequest {
	var rows []imModel.FriendRequest
	db.DB.WithContext(ctx).Model(&imModel.FriendRequest{}).
		Where("sender_id = ? AND sender_type = ? AND status = ?", userID, userType, "pending").
		Order("created_at DESC").Find(&rows)
	return rows
}

func RemoveFriendshipPair(ctx context.Context, userID, userType, friendID, friendType string) (int64, int64, error) {
	tx := db.DB.WithContext(ctx).Begin()
	r1 := tx.Where("user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?",
		userID, userType, friendID, friendType).Delete(&imModel.Friendship{})
	if r1.Error != nil {
		tx.Rollback()
		return 0, 0, r1.Error
	}
	r2 := tx.Where("user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?",
		friendID, friendType, userID, userType).Delete(&imModel.Friendship{})
	if r2.Error != nil {
		tx.Rollback()
		return 0, 0, r2.Error
	}
	if err := tx.Commit().Error; err != nil {
		return 0, 0, err
	}
	return r1.RowsAffected, r2.RowsAffected, nil
}

func SearchSysUsers(ctx context.Context, like string, limit int) []sysUser.SysUser {
	var users []sysUser.SysUser
	db.DB.WithContext(ctx).Model(&sysUser.SysUser{}).
		Where("username LIKE ? OR nickname LIKE ?", like, like).
		Limit(limit).Find(&users)
	return users
}

func SearchClientUsers(ctx context.Context, like string, limit int) []cliUser.ClientUser {
	var users []cliUser.ClientUser
	db.DB.WithContext(ctx).Model(&cliUser.ClientUser{}).
		Where("username LIKE ? OR nickname LIKE ?", like, like).
		Limit(limit).Find(&users)
	return users
}

func CountBlocks(ctx context.Context, userID, userType, blockedID, blockedType string) int64 {
	var count int64
	db.DB.WithContext(ctx).Model(&imModel.FriendBlock{}).
		Where("user_id = ? AND user_type = ? AND blocked_id = ? AND blocked_type = ?",
			userID, userType, blockedID, blockedType).Count(&count)
	return count
}

func CreateBlock(ctx context.Context, block *imModel.FriendBlock) error {
	return db.DB.WithContext(ctx).Create(block).Error
}

func DeleteFriendshipForBlock(ctx context.Context, userID, userType, blockedID, blockedType string) {
	db.DB.WithContext(ctx).
		Where("(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?) OR "+
			"(user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?)",
			userID, userType, blockedID, blockedType,
			blockedID, blockedType, userID, userType).
		Delete(&imModel.Friendship{})
}

func DeleteBlock(ctx context.Context, userID, userType, blockedID, blockedType string) (int64, error) {
	res := db.DB.WithContext(ctx).
		Where("user_id = ? AND user_type = ? AND blocked_id = ? AND blocked_type = ?",
			userID, userType, blockedID, blockedType).
		Delete(&imModel.FriendBlock{})
	return res.RowsAffected, res.Error
}

func ListBlocks(ctx context.Context, userID, userType string) []imModel.FriendBlock {
	var rows []imModel.FriendBlock
	db.DB.WithContext(ctx).Model(&imModel.FriendBlock{}).
		Where("user_id = ? AND user_type = ?", userID, userType).
		Find(&rows)
	return rows
}

func UpdateRemark(ctx context.Context, userID, userType, friendID, friendType, remark string) error {
	return db.DB.WithContext(ctx).Model(&imModel.Friendship{}).
		Where("user_id = ? AND user_type = ? AND friend_id = ? AND friend_type = ?",
			userID, userType, friendID, friendType).
		Update("remark", remark).Error
}
