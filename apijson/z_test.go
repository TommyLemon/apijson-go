package apijson

import (
	"context"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"my-apijson/apijson/consts"
	"my-apijson/apijson/db"
	"my-apijson/apijson/util"
	"testing"
)

func TestTwoTableGet(t *testing.T) {
	req := `
{
 "User":{
        "id@":"Todo/userId"
    },
    "Todo":{
        "id":1627794043692
    }
   
}
`
	ctx := context.TODO()
	reqMap := gjson.New(req).Map()
	out, err := Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

func TestTowTableGetList(t *testing.T) {
	req := `
{
 	"[]":{
	"User":{
        "id@":"/Todo/userId"
    },
    "Todo":{
       
    }
	}
   
}
`
	ctx := context.TODO()
	reqMap := gjson.New(req).Map()
	out, err := Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

func TestCheckRequest(t *testing.T) {

	db.Init()

	req := `
{
    "Todo":{
        
		"title":"asdasda"
    },
	"tag":"Todo"
}
`
	ctx := context.TODO()
	reqMap := gjson.New(req).Map()

	out, err := Post(ctx, reqMap)
	if err != nil {
		g.Dump(err)
	}
	g.Dump(out)
}

func TestAccess(t *testing.T) {

	db.Init()

	req := `
{
 "User":{
        "id@":"Todo/userId"
    },
    "Todo":{
        "id":1
    }
   
}
`
	ctx := context.TODO()

	ctx = context.WithValue(ctx, "ajg.userId", "2")
	ctx = context.WithValue(ctx, consts.RoleKey, []string{consts.LOGIN, consts.OWNER})

	AccessCondition = accessCondition
	AccessVerify = true

	reqMap := gjson.New(req).Map()
	out, err := Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

func accessCondition(ctx context.Context, table string, req g.Map, needRole []string) (g.Map, error) {

	userRole := ctx.Value(consts.RoleKey).([]string)

	// 可改成switch方式

	if util.Contains(needRole, consts.UNKNOWN) {
		return nil, nil
	}

	if util.Contains(needRole, consts.LOGIN) && util.Contains(userRole, consts.LOGIN) { // 登录后公开资源
		return nil, nil
	}

	if util.Contains(needRole, consts.OWNER) && util.Contains(userRole, consts.OWNER) {
		if table == "User" {
			return g.Map{
				"id": ctx.Value("ajg.userId"),
			}, nil
		} else {
			return g.Map{
				"userId": ctx.Value("ajg.userId"),
			}, nil
		}
	}

	return nil, nil
}
