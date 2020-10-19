package controller

import (
	"chatroom/model"
	"chatroom/service/websocket"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Err error
}

/***
 * 創建新的使用者，密碼用sha256雜湊加密
 */
func (u *UserController) Signup(context *gin.Context) {
	userId := context.PostForm("userId")
	userName := context.PostForm("userName")
	password := context.PostForm("password")
	passwordAgain := context.PostForm("passwordAgain")

	// 檢查密碼與再次輸入的密碼是否吻合
	if password != passwordAgain {
		context.HTML(http.StatusOK, "signup.html", "密碼與再次輸入的密碼不符")
		return
	}

	// 密碼加密
	hash := sha256.New()
	io.WriteString(hash, password)
	password = fmt.Sprintf("%x", hash.Sum(nil))

	err := model.User.Create(userId, userName, password)
	if err != nil {
		context.HTML(http.StatusOK, "signup.html", "此帳號已被註冊過")
		return
	}

	context.HTML(http.StatusOK, "signupSuccess.html", nil)
}

/***
 * 登入，檢查密碼是否吻合
 */
func (u *UserController) Login(context *gin.Context) {
	userId := context.PostForm("userId")
	inputPassword := context.PostForm("password")

	password, err := model.User.GetPassword(userId)
	if err != nil {
		context.HTML(http.StatusOK, "login.html", "帳號或密碼錯誤")
		return
	}

	// 加密
	hash := sha256.New()
	io.WriteString(hash, inputPassword)
	inputPassword = fmt.Sprintf("%x", hash.Sum(nil))
	if inputPassword != password {
		context.HTML(http.StatusOK, "login.html", "帳號或密碼錯誤")
		return
	}

	context.Redirect(http.StatusMovedPermanently, "/chatroom?userId="+userId)
}

/***
 * 登出，將service client刪除，終止client所執行的goroutines
 */
func (u *UserController) Logout(context *gin.Context) {
	userId := context.Param("userId")
	websocket.Destroy(userId)
	context.Redirect(http.StatusFound, "/login")
}
