package storage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"social-todo-list/common"
	"social-todo-list/module/userlikeitem/model"
)

func (s *sqlStore) Find(ctx context.Context, userId, itemId int) (*model.Like, error) {
	var data model.Like
	if err := s.db.Where("user_id = ? and item_id = ?", userId, itemId).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &data, nil
}
