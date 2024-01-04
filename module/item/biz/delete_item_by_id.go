package biz

import (
	"context"
	"errors"
	"social-todo-list/common"
	"social-todo-list/module/item/model"
)

type DeleteItemStorage interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	DeleteItem(ctx context.Context, cond map[string]interface{}) error
}

type deleteItemBiz struct {
	store     DeleteItemStorage
	requester common.Requester
}

func NewDeleteItemBiz(store DeleteItemStorage, requester common.Requester) *deleteItemBiz {
	return &deleteItemBiz{store: store, requester: requester}
}

func (biz *deleteItemBiz) DeleteItemById(ctx context.Context, id int) error {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return common.ErrCannotGetEntity(model.EntityName, err)
	}

	if data.Status == "Deleted" {
		return model.ErrItemIsDeleted
	}

	isOwner := biz.requester.GetUserId() == data.UserId
	if !isOwner && !common.IsAdmin(biz.requester) {
		return common.ErrNoPermission(errors.New("no permission"))
	}

	if err := biz.store.DeleteItem(ctx, map[string]interface{}{"id": id}); err != nil {
		return common.ErrCannotDeleteEntity(model.EntityName, err)
	}

	return nil
}
