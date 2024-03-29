package ginuser

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"social-todo-list/common"
	"social-todo-list/module/user/biz"
	"social-todo-list/module/user/model"
	"social-todo-list/module/user/storage"
	"social-todo-list/plugin/tokenprovider"
)

func Login(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		tokenProvider := serviceCtx.MustGet(common.PluginJWT).(tokenprovider.Provider)

		var loginUserData model.UserLogin

		if err := c.ShouldBind(&loginUserData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewSQLStore(db)
		md5 := common.NewMd5Hash()

		business := biz.NewLoginBusiness(store, tokenProvider, md5, 60*60*24*30)
		account, err := business.Login(c.Request.Context(), &loginUserData)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
