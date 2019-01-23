package models

import ("github.com/astaxie/beego/orm"
		_"github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int
	Name string `orm:"unique"`
	PassWord string
	Articles []*Article `orm:"rel(m2m)"`
}

type Article struct {
	Id2 int  `orm:"pk;auto"`
	Title string `orm:"size(40)"`
	Content string `orm:"size(100)"`
	ReadCount int 	`orm:"default(0)"`
	Time time.Time	`orm:"type(datetime);auto_now_add"`
	Img string	`orm:"null"`
	ArticleType *ArticleType `orm:"rel(fk)"`
	Users []*User `orm:"reverse(many)"`
}

type ArticleType struct {
	Id int
	TypeName string
	Articles []*Article `orm:"reverse(many)"`
}

func init()  {
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/NewsPublish?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}

