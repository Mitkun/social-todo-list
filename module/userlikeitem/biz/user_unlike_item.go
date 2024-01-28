package biz

import (
	"context"
	"social-todo-list/common"
	"social-todo-list/module/userlikeitem/model"
)

type UserULikeItemStoreLikeItemStore interface {
	Find(ctx context.Context, userId, itemId int) (*model.Like, error)
	Delete(ctx context.Context, userId, itemId int) error
}

type userUnlikeItemBiz struct {
	store UserULikeItemStoreLikeItemStore
}

func NewUserUnlikeItemBiz(store UserULikeItemStoreLikeItemStore) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{store: store}
}

func (biz *userUnlikeItemBiz) UnlikeItem(ctx context.Context, userId, itemId int) error {
	_, err := biz.store.Find(ctx, userId, itemId)

	// Delete if data existed
	if err == common.RecordNotFound {
		return model.ErrDisNotLikeItem(err)
	}

	if err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	if err := biz.store.Delete(ctx, userId, itemId); err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	return nil
}
