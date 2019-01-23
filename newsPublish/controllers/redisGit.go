package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

type RedisGit struct {
	beego.Controller
}

func (this *RedisGit) ShowRedis()  {
	conn,err:=redis.Dial("tcp",":6379")
	if err != nil {
		beego.Error("连接错误")
		return
	}

	//操作函数
	resp,err:=conn.Do("mget","kk","ll")
	//conn.Send("set","kk","vv")
	//conn.Flush()
	//conn.Receive()

	//回复助手函数  类型转换
	result,_:=redis.Values(resp,err)
	var v1 string
	var v3 int
	redis.Scan(result,&v1,&v3)
}