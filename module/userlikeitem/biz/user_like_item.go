package biz

import (
	"context"
	"log"
	"social-todo-list/common"
	"social-todo-list/module/userlikeitem/model"
)

type UserLikeItemStore interface {
	Create(ctx context.Context, data *model.Like) error
}

type IncreaseItemStorage interface {
	IncreaseLikeCount(ctx context.Context, id int) error
}

type userLikeItemBiz struct {
	store     UserLikeItemStore
	itemStore IncreaseItemStorage
}

func NewUserLikeItemBiz(store UserLikeItemStore, increaseLike IncreaseItemStorage) *userLikeItemBiz {
	return &userLikeItemBiz{store: store, itemStore: increaseLike}
}

func (biz *userLikeItemBiz) LikeItem(ctx context.Context, data *model.Like) error {
	if err := biz.store.Create(ctx, data); err != nil {
		return model.ErrCannotLikeItem(err)
	}

	go func() {
		defer common.Recovery()

		if err := biz.itemStore.IncreaseLikeCount(ctx, data.ItemId); err != nil {
			log.Println(err)
		}
	}()

	return nil
}
