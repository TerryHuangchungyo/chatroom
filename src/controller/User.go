package controller

import (
	"crypto/sha256"
	"fmt"
	"io"
	"local/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Err error
}

func (u *UserController) Create(context *gin.Context) {
	userId := context.PostForm("userId")
	userName := context.PostForm("userName")
	password := context.PostForm("password")
	passwordAgain := context.PostForm("passwordAgain")

	if password != passwordAgain {
		context.HTML(http.StatusOK, "signup.html", "密碼與再次輸入的密碼不符")
		return
	}

	hash := sha256.New()
	io.WriteString(hash, password)
	password = fmt.Sprintf("%x", hash.Sum(nil))

	err := model.User.Create(userId, userName, password)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	context.HTML(http.StatusOK, "signupSuccess.html", nil)
}
