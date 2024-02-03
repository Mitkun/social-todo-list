package biz

import (
	"context"
	"log"
	"social-todo-list/common"
	"social-todo-list/module/userlikeitem/model"
)

type UserULikeItemStoreLikeItemStore interface {
	Find(ctx context.Context, userId, itemId int) (*model.Like, error)
	Delete(ctx context.Context, userId, itemId int) error
}

type DecreaseLikeCount interface {
	DecreaseLikeCount(ctx context.Context, id int) error
}

type userUnlikeItemBiz struct {
	store     UserULikeItemStoreLikeItemStore
	itemStore DecreaseLikeCount
}

func NewUserUnlikeItemBiz(store UserULikeItemStoreLikeItemStore, decreaseLike DecreaseLikeCount) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{store: store, itemStore: decreaseLike}
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

	go func() {
		defer common.Recovery()

		if err := biz.itemStore.DecreaseLikeCount(ctx, itemId); err != nil {
			log.Println(err)
		}

	}()
	return nil
}
