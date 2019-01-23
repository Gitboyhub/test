package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"newsPublish/models"
	"math"
)

type ArticleController struct {
	beego.Controller
}

//展示首页内容
func (this *ArticleController) ShowIndex()  {
	//登录校验
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login",302)
		return
	}
	this.Layout = "layout.html"
	this.TplName="index.html"
	//获取orm对象
	o:=orm.NewOrm()
	//指定要查询的表
	qs:=o.QueryTable("article")
	//定义一个容器来接收查询内容
	var artilces []models.Article
	//查询所有的表内容
	//qs.All(&artilce)
	//实现分页功能
	pageSize := 2
	//处理首页末页
	pageIndex,err:=this.GetInt("pageIndex")
	if err != nil {
		pageIndex =1
	}
	start:=pageSize * (pageIndex - 1)

	//下拉框改变的时候，获取不同类型的文章数据
	var count int64
	typeName := this.GetString("select")
	if typeName == ""{
		count,_= qs.RelatedSel("ArticleType").Count()
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&artilces)
	}else {
		count,_= qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&artilces)
	}
	//获取到类型数据,根据这个数据获取相应文章

	//默认多表查询是惰性查询

	pageCount := math.Ceil(float64(count)/float64(pageSize))
	this.Data["pageIndex"] = pageIndex
	this.Data["typeName"] = typeName
	//将数据传递给前端
	this.Data["articles"]=artilces

	//获取所有类型数据并传递给前段展示
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"] = articleTypes


	this.Data["count"] = count
	this.Data["pageCount"] = pageCount

}

//显示添加新闻
func (this *ArticleController) ShowAdd()  {
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"] = articleTypes
	this.Layout = "layout.html"
	this.TplName="add.html"
}
//处理添加新闻
func (this *ArticleController) HandleArticle()  {
	//获取数据
	articleName :=this.GetString("articleName")
	content :=this.GetString("content")
	//file
	file,head,err :=this.GetFile("uploadname")

	//获取数据
	if articleName == "" || content == "" || err != nil{
		beego.Error("获取用户添加数据失败",err)
		this.TplName = "add.html"
		return
	}
	defer file.Close()
	//需要判断大小
	if head.Size > 5000000{
		beego.Error("图片太大，我不收")
		this.TplName = "add.html"
		return
	}
	//需要校验格式
	ext :=path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		beego.Error("文件格式不正确")
		this.TplName = "add.html"
		return
	}

	//防止重名
	//beego.Info("time.now = ",time.Now().Format("2006-01-02 15:04:05"))
	fileName := time.Now().Format("20060102150405")
	//操作数据
	this.SaveToFile("uploadname","./static/img/"+fileName+ext)

	//把数据插入到数据库
	//获取orm对象
	o := orm.NewOrm()
	//获取插入对象
	var article models.Article
	//给插入对象赋值
	article.Title = articleName
	article.Content = content
	article.Img = "/static/img/"+fileName+ext

	//获取类型数据
	typeName := this.GetString("select")

	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType,"TypeName")

	article.ArticleType = &articleType

	//插入
	o.Insert(&article)

	//返回数据
	this.Redirect("/article/index",302)
}

//显示详情页内容
func (this *ArticleController) ShowArticle()  {
	//获取数据
	id,err:=this.GetInt("id")
	//校验数据
	if err != nil {
		beego.Error("获取数据错误",err)
		return
	}
	//处理数据
	//从数据库查询数据
	//定义orm对象
	o:=orm.NewOrm()
	//定义查询对象
	var article models.Article
	//定义查询条件
	article.Id2 = id
	//查询
	o.Read(&article)
	//每读取一次数据，阅读量加1
	article.ReadCount += 1
	o.Update(&article)
	//返回数据
	this.Data["article"] = article
	//指定视图
	this.Data["title"] = "文章详情"
	this.Layout = "layout.html"
	this.TplName = "content.html"
}

//编辑详情页
func (this *ArticleController) ShowEditArticle()  {
	//获取数据
	id,err:=this.GetInt("id")
	//校验数据
	if err != nil {
		beego.Error("获取数据错误")
		this.TplName = "index.html"
		return
	}
	//处理数据
	//查询
	//定义orm对象
	o:=orm.NewOrm()
	//定义查询对象
	var article models.Article
	//指定查询条件
	article.Id2 = id
	//查询数据
	o.Read(&article)

	//返回数据
	this.Data["article"] = article
	this.Layout = "layout.html"
	this.TplName = "update.html"
}

//封装上传文件函数 做个函数接口
func UploadFunc(this *ArticleController,fileName string) string {
	file,head,err :=this.GetFile(fileName)

	//获取数据
	if err != nil{
		beego.Error("获取用户添加数据失败",err)
		this.TplName = "add.html"
		return ""
	}
	defer file.Close()
	//需要判断大小
	if head.Size > 5000000{
		beego.Error("图片太大，我不收")
		this.TplName = "add.html"
		return ""
	}
	//需要校验格式
	ext :=path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		beego.Error("文件格式不正确")
		this.TplName = "add.html"
		return ""
	}

	//防止重名
	//beego.Info("time.now = ",time.Now().Format("2006-01-02 15:04:05"))
	filePath := time.Now().Format("20060102150405")
	//操作数据
	this.SaveToFile(fileName,"./static/img/"+filePath+ext)
	return "/static/img/"+filePath+ext
}

//处理更新详情页
func (this *ArticleController) HandleEditArticle()  {
	//获取数据
	id,err :=this.GetInt("id")
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	filePath := UploadFunc(this,"uploadname")
	//校验数据
	if err != nil || articleName == "" || content == "" || filePath == ""{
		beego.Error("获取数据错误")
		this.TplName = "update.html"
		return
	}

	//处理数据
	//更新
	//获取orm对象
	o := orm.NewOrm()
	//获取更新对象
	var article models.Article
	//给更新条件赋值
	article.Id2 = id
	//先read一下，判断要更新的数据
	err = o.Read(&article)
	//更新
	if err != nil{
		beego.Error("更新数据不存在")
		this.TplName = "update.html"
		return
	}
	article.Title = articleName
	article.Content = content
	article.Img = filePath
	o.Update(&article)

	//返回数据
	this.Redirect("/index",302)
}

//处理删除
func (this *ArticleController) HandleDelete()  {
	id,err:=this.GetInt("id")
	if err != nil {
		beego.Error("删除请求数据失败")
		this.TplName = "index.html"
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id2 = id
	_,err=o.Delete(&article)
	if err != nil {
		beego.Error("删除失败")
		this.TplName = "index.html"
		return
	}
	this.Redirect("/index",302)
}

//展示类型添加界面
func (this *ArticleController) ShowTypeAdd()  {
	o:=orm.NewOrm()
	qs:=o.QueryTable("articleType")
	var articleTypes []models.ArticleType
	qs.All(&articleTypes)
	this.Data["articleTypes"] = articleTypes
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

//添加类型
func (this *ArticleController) HandleAddType() {
	//获取数据
	TypeName := this.GetString("typeName")
	//校验数据
	if TypeName == "" {
		beego.Error("文章类型不能为空")
		//this.TplName = "addType.html"  //返回一个空的addType页面
		this.Redirect("/article/addType", 302)  //跳转一个正常的addType页面
		return
	}
	//处理数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = TypeName
	o.Insert(&articleType)
	//返回数据
	this.Redirect("/article/addType", 302)
}

//删除类型操作
func (this *ArticleController) DeleteType()  {
	//获取数据并校验
	id,err:=this.GetInt("id")
	if err != nil {
		beego.Error("获取删除数据失败",err)
		this.TplName = "addType.html"
		return
	}
	//处理数据
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	o.Delete(&articleType)
	this.Redirect("/article/addType",302)
}