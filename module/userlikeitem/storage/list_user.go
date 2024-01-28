package storage

import (
	"context"
	"social-todo-list/common"
	"social-todo-list/module/userlikeitem/model"
)

func (s *sqlStore) ListUsers(ctx context.Context, itemId int, paging *common.Paging) ([]common.SimpleUser, error) {
	var result []model.Like

	db := s.db.Where("item_id = ?", itemId)

	if err := db.Table(model.Like{}.TableName()).Select("user_id").Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	if err := db.Select("*").
		Order("created_at desc").
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Preload("User").
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	users := make([]common.SimpleUser, len(result))

	for i := range users {
		users[i] = *result[i].User
		users[i].CreatedAt = result[i].CreatedAt
		users[i].UpdatedAt = nil
	}

	return users, nil
}