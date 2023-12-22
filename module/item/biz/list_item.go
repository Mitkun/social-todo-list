package biz

import (
	"context"
	"g09-social-todo-list/common"
	"g09-social-todo-list/module/item/model"
)

type ListItemStorage interface {
	ListItem(
		ctx context.Context,
		filter *model.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]model.TodoItem, error)
}

type listItemBiz struct {
	store ListItemStorage
}

func NewListItemBiz(store ListItemStorage) *listItemBiz {
	return &listItemBiz{store: store}
}

func (biz *listItemBiz) ListItem(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]model.TodoItem, error) {
	data, err := biz.store.ListItem(ctx, filter, paging)
	if err != nil {
		return nil, err
	}

	return data, nil
}