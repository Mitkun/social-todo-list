package storage

import (
	"context"
	"g09-social-todo-list/module/item/model"
)

func (s *sqlStore) CreateItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := s.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}