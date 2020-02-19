package main

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/boltdb/bolt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"main/database"
	"main/env"
	"main/rcache"
	"main/reply"
	"main/resources"
	"regexp"
	"strconv"
	"strings"
)

var resource = regexp.MustCompile(`#(\d*)`)

func MuxHandler(ctx *context.Context) {
	input := ctx.Input.Query("input")

	if strings.Index(input, "找书") >= 0 {
		temp := strings.Split(strings.TrimSpace(input), " ")
		if len(temp) != 2 {
			rb := reply.NewBaseMessage("text", reply.HELPBook)
			ctx.ResponseWriter.Write(rb)
			return
		}

		var content []byte
		err := rcache.GetDb().View(func(tx *bolt.Tx) error {
			b := tx.Bucket(rcache.DBResource)
			if b != nil {
				content = b.Get([]byte(temp[1]))
				if content == nil {
					return errors.New("key not exists")
				}
			}
			return nil
		})

		if err == nil {
			ctx.ResponseWriter.Write(content)
			return
		}

		log.Println("进入数据库")

		ds := database.NewSessionStore()
		defer ds.Close()
		con := ds.C("resources")

		var res []resources.Resources
		err = con.Find(bson.M{"file_name": bson.RegEx{Pattern: temp[1], Options: "i"}}).All(&res)
		if err != nil || len(res) == 0 {
			rb := reply.NewBaseMessage("text", "没有找到相关资源")
			ctx.ResponseWriter.Write(rb)
			rcache.ResourceChan <- rcache.Resource{Key: []byte(temp[1]), Value: rb}
			return
		}

		var buf strings.Builder
		buf.WriteString("查询结果\n")
		for _, k := range res {
			buf.WriteString(k.GetFileName())
			buf.WriteString("\n")
		}
		buf.WriteString("回复编号例如「#1」获取📚资源")

		rb := reply.NewBaseMessage("text", buf.String())
		ctx.ResponseWriter.Write(rb)
		rcache.ResourceChan <- rcache.Resource{Key: []byte(temp[1]), Value: rb}
		return
	}

	if resource.MatchString(input) { // 查询资源

		res := resource.FindStringSubmatch(input)
		if len(res) < 2 {
			rb := reply.NewBaseMessage("text", "指令有误，回复编号例如「#1」获取1号资源")
			ctx.ResponseWriter.Write(rb)
			return
		}

		i, err := strconv.Atoi(res[1])
		if err != nil {
			rb := reply.NewBaseMessage("text", "指令有误，回复编号例如「#1」获取1号资源")
			ctx.ResponseWriter.Write(rb)
			return
		}

		var content []byte = nil
		err = rcache.GetDb().View(func(tx *bolt.Tx) error {
			b := tx.Bucket(rcache.DBResource)

			content = b.Get([]byte(res[1]))
			if content == nil {
				return errors.New("key not exists")
			}
			return nil
		})

		if err == nil {
			ctx.ResponseWriter.Write(content)
			return
		}
		fmt.Println("进入数据库")

		ds := database.NewSessionStore()
		defer ds.Close()
		var r resources.Resources
		err = ds.C("resources").Find(bson.M{"id": i}).One(&r)
		if err != nil {
			rb := reply.NewBaseMessage("text", "未找到资源或资源出错，将报告杂货铺，随后为您提供更优质的服务")
			ctx.ResponseWriter.Write(rb)
			return
		}

		rb := reply.NewBaseMessage("text", strings.Join([]string{
			"找到了资源，请查收\n\n", r.GetFileName(), "\n下载链接\n", r.Link,
		}, ""))
		ctx.ResponseWriter.Write(rb)
		rcache.ResourceChan <- rcache.Resource{Key: []byte(res[1]), Value: rb,}
		return
	}

	ctx.ResponseWriter.Write(reply.NewBaseMessage("text", reply.HELPBook))
	return
}
func main() {

	go rcache.Saver()
	beego.Get("/get", MuxHandler)
	beego.Get("/flush", rcache.FlushAll)
	beego.Get("/test_cache", rcache.TestCache)
	beego.Run(fmt.Sprintf(":%s", env.Port))
}
