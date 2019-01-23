package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsPublish/models"
)

type UserController struct {
	beego.Controller
}

//展示注册界面
func (this *UserController) ShowRigister() {
	//指定视图
	this.TplName="register.html"
}

//处理注册数据
func (this *UserController) HandleRigister()  {
	//获取数据
	userName:=this.GetString("userName")
	passWord:=this.GetString("password")
	//处理数据
	if userName==""||passWord==""{
		beego.Error("用户名或者密码不能为空")
		this.TplName="register.html"
		return
	}
	//获取orm对象
	o:=orm.NewOrm()
	//获取插入对象
	var user models.User
	//给插入对象赋值
	user.Name=userName
	user.PassWord=passWord
	//插入数据
	_,err:=o.Insert(&user)
	if err != nil {
		beego.Error("用户注册失败")
		this.TplName="register.html"
		return
	}
	//this.Ctx.WriteString("恭喜您，注册成功")
	this.Redirect("/login",302)
}

//展示登录界面
func (this *UserController) ShowLogin()  {

	//实现记住用户名功能
	userName:=this.Ctx.GetCookie("userName")
	if userName == ""{
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}else {
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	}
	this.TplName="login.html"
}

//处理登录操作
func (this *UserController) HandleLogin() {
	userName:=this.GetString("userName")
	passWord:=this.GetString("password")
	if userName==""||passWord==""{
		beego.Error("用户名或密码为空")
		this.TplName="login.html"
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	err:=o.Read(&user,"Name")
	if err != nil {
		beego.Error("用户不存在")
		this.TplName="login.html"
		return
	}
	if user.PassWord !=passWord{
		beego.Error("用户密码输入错误")
		this.TplName="login.html"
		return
	}
	//登录成功的情况下，选中复选框把用户名存储到cookie里面
	remeber := this.GetString("remember")
	if remeber == "on"{
		this.Ctx.SetCookie("userName",userName,60 * 60 * 24)
	}else{
		this.Ctx.SetCookie("userName",userName,-1)
	}
	//this.Ctx.WriteString("登录成功")

	this.SetSession("userName",userName)

	this.Redirect("article/index",302)
}

//处理退出操作
func (this *UserController) Logout() {
	//删除登录状态（删除session数据）
	this.DelSession("userName")
	this.Redirect("/login",302)
}