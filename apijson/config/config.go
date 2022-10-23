package config

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

type AccessConditionFunc func(ctx context.Context, table string, req g.Map, needRole []string) (g.Map, error)

var (
	// AccessVerify 是否权限验证
	AccessVerify = false
	// AccessCondition 自定义权限限制条件
	AccessCondition AccessConditionFunc
)

var (
	RoleKey   = "ajg.role" // ctx 中role 的key
	UserIdKey = "ajg.userId"
)

var (
	TableAccess  = "_access"
	TableRequest = "_request"
)
