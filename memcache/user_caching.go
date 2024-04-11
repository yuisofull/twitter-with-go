package memcache

import (
	"context"
	"fmt"
	"twitter/modules/user/usermodel"
)

type FindStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
}

type userCaching struct {
	store     Caching
	findStore FindStore
}

func NewUserCaching(store Caching, findStore FindStore) *userCaching {
	return &userCaching{
		store:     store,
		findStore: findStore,
	}
}

func (uc *userCaching) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error) {
	userId := conditions["id"].(int)
	key := fmt.Sprintf("user-%d", userId)

	userInCache := uc.store.Read(key)

	if userInCache != nil {
		return userInCache.(*usermodel.User), nil
	}

	user, err := uc.findStore.FindUser(ctx, conditions, moreInfo...)

	if err != nil {
		return nil, err
	}

	// Update cache
	uc.store.Write(key, user)

	return user, nil
}
