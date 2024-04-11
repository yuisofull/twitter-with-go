package userbiz

import (
	"context"
	"twitter/common"
	"twitter/component/appctx"
	"twitter/component/tokenprovider"
	"twitter/modules/user/usermodel"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
}

//type TokenConfig interface {
//	GetAtExp() int
//	GetRtExp() int
//}

type loginBusiness struct {
	//appCtx    appctx.AppContext
	userStore LoginStorage
	//tokenConfig   TokenConfig
	expiry        int // expiry will replace for type TokenConfig
	tokenProvider tokenprovider.Provider
	hasher        Hasher
}

func NewLoginBusiness(appCtx appctx.AppContext,
	userStore LoginStorage,
	//tokenConfig TokenConfig,
	expiry int,
	tokenProvider tokenprovider.Provider,
	hasher Hasher,
) *loginBusiness {
	return &loginBusiness{
		//appCtx:        appCtx,
		userStore:     userStore,
		expiry:        expiry,
		tokenProvider: tokenProvider,
		hasher:        hasher,
		//tokenConfig:   tokenConfig,
	}
}

// 1. Find user, email
// 2. Hash pass from input & compare with pass in db
// 3. Provider: issue JWT token for Client
// 3.1 Access token & Refresh token
// 4. Return token(s)

func (biz *loginBusiness) Login(ctx context.Context, data *usermodel.UserLogin) (*tokenprovider.Token, error) {
	user, err := biz.userStore.FindUser(ctx, map[string]interface{}{"email": data.Email})

	if err != nil {
		return nil, usermodel.ErrEmailOrPasswordInvalid
	}

	passwordHashed := biz.hasher.Hash(data.Password + user.Salt)

	if user.Password != passwordHashed {
		return nil, usermodel.ErrEmailOrPasswordInvalid
	}

	payload := tokenprovider.TokenPayload{
		UserId: user.Id,
		Role:   user.Role,
	}

	//business.tokenConfig.GetAtExp() ===> business.expiry
	accessToken, err := biz.tokenProvider.Generate(payload, biz.expiry)

	if err != nil {
		return nil, common.ErrInternal(err)
	}

	//refreshToken, err := business.tokenProvider.Generate(payload, business.expiry)
	//
	//if err != nil {
	//	return nil, common.ErrInternal(err)
	//}
	//
	//account := model.NewAccount(accessToken, refreshToken)

	return accessToken, nil
}
