package routers

import (
	"newsPublish/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*",beego.BeforeExec,filterFUnc)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{}, "get:ShowRigister;post:HandleRigister")
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/article/index", &controllers.ArticleController{}, "get:ShowIndex")
	beego.Router("/article/add", &controllers.ArticleController{}, "get:ShowAdd;post:HandleArticle")
	beego.Router("/article/content", &controllers.ArticleController{}, "get:ShowArticle")
	beego.Router("/article/editArticle",&controllers.ArticleController{},"get:ShowEditArticle;post:HandleEditArticle")
	beego.Router("/article/DeleteArticle",&controllers.ArticleController{},"get:HandleDelete")
	beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowTypeAdd;post:HandleAddType")
	beego.Router("/article/logout",&controllers.UserController{},"get:Logout")
	beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteType")
	beego.Router("/redis",&controllers.RedisGit{},"get:ShowRedis")
}

func filterFUnc(ctx *context.Context)  {
	userName := ctx.Input.Session("userName")
	if userName == nil{
		ctx.Redirect(302,"/login")
		return
	}
}