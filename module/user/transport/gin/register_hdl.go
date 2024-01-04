package ginuser

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"social-todo-list/common"
	"social-todo-list/module/user/biz"
	"social-todo-list/module/user/model"
	"social-todo-list/module/user/storage"
)

func Register(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var data model.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		store := storage.NewSQLStore(db)
		md5 := common.NewMd5Hash()
		biz := biz.NewRegisterBusiness(store, md5)

		if err := biz.Register(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
