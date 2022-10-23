package apijson

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

// AccessVerify 是否权限验证
var AccessVerify = false

// AccessCondition 自定义权限限制条件
var AccessCondition AccessConditionFunc

type AccessConditionFunc func(ctx context.Context, table string, req g.Map, needRole []string) (g.Map, error)
