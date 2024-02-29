package controllers

import (
	"EasyGin/app/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func (u UserController) GetUser(ctx *gin.Context) {
	//url参数获取
	//name:= ctx.DefaultQuery("name", "枯藤")
	//表单获取
	//types := c.DefaultPostForm("type", "post")
	//username := c.PostForm("username")
	user := services.UsersService{}.GetUser(ctx, "12345")
	if user == nil {
		ctx.JSON(500, gin.H{
			"message": "is not user",
		})
	} else {
		jsonData, err := json.Marshal(user)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
		}
		ctx.JSON(200, gin.H{
			"message": fmt.Sprintf("%s", jsonData),
		})
	}

}
