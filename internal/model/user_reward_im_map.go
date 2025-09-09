package model

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type UserRewardImMap struct {
	Id           int64  `gorm:"column:id"`
	OldMessageId int64  `gorm:"column:old_message_id"` // 老的消息id
	NewMessageId string `gorm:"column:new_message_id"` // 新的消息id
	Status       int8   `gorm:"column:status"`         // 是否已经发过,0否1是
	CreatedTime  int64  `gorm:"column:created_time"`
	UpdatedTime  int64  `gorm:"column:updated_time"`
	IsDeleted    int8   `gorm:"column:is_deleted"`
}

func (a UserRewardImMap) TableName() string {
	return "im_message.user_reward_im_map"
}

// DescribeUserRewardImMapById 根据mid找到变动分组的数据
func DescribeUserRewardImMapById(log *logrus.Entry, db *gorm.DB, q *UserRewardImMap) (*UserRewardImMap, error) {
	a := &UserRewardImMap{}
	var userRewardImMapList []*UserRewardImMap

	qs := db.Table(a.TableName())

	if q.OldMessageId > 0 {
		qs = qs.Where("old_message_id = ? ", q.OldMessageId)
	}

	if q.NewMessageId != "" {
		qs = qs.Where("new_message_id = ? ", q.NewMessageId)
	}

	if err := qs.Find(&userRewardImMapList).Error; err != nil {
		log.Errorf("err %v", err)
		return nil, err
	}

	if len(userRewardImMapList) > 0 {
		return userRewardImMapList[0], nil
	}

	return nil, nil
}

// CreateOrUpdateUserRewardImMap  创建迁移的数据
func CreateOrUpdateUserRewardImMap(log *logrus.Entry, db *gorm.DB, u *UserRewardImMap) (err error) {
	info, err := DescribeUserRewardImMapById(log, db, u)
	if err != nil {
		return
	}

	if info != nil && info.Id > 0 {
		updateMap := map[string]interface{}{
			"updated_time": time.Now().Unix(),
		}

		if u.OldMessageId > 0 {
			updateMap["old_message_id"] = u.OldMessageId
		}

		if u.NewMessageId != "" {
			updateMap["new_message_id"] = u.NewMessageId
		}

		// 更新
		rows := db.Table(u.TableName()).Where("id = ?", info.Id).Updates(updateMap).RowsAffected
		if err != nil {
			log.Errorf("UpdateAnchorStrategyData err %v", err)
			return
		}

		if rows <= 0 {
			log.Error("更新失败")
			err = errors.New("更新失败")
		}
		return
	}

	if cerr := db.Table(u.TableName()).Create(u).Error; cerr != nil {
		log.Errorf("cerr %v", cerr)
		err = cerr
		return
	}
	return
}
