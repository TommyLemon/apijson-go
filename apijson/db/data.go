package db

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"my-apijson/apijson/config"
	"strings"
)

var AccessMap = map[string]Access{}

type Access struct {
	Debug     int8
	Name      string
	Alias     string
	Get       []string
	Head      []string
	Gets      []string
	Heads     []string
	Post      []string
	Put       []string
	Delete    []string
	CreatedAt *gtime.Time
	Detail    string

	// ext

	RowKey string
}

var RequestMap = map[string]Request{}

type Request struct {
	Debug   int8
	Version int16
	Method  string
	Tag     string
	// https://github.com/Tencent/APIJSON/blob/master/APIJSONORM/src/main/java/apijson/orm/Operation.java
	Structure g.Map
	Detail    string
	CreatedAt *gtime.Time
}

func Init() {
	getAccessMap()
	getRequestMap()
}

func getAccessMap() {
	accessMap := make(map[string]Access)

	var accessList []Access
	g.DB().Model(config.TableAccess).Scan(&accessList)

	for _, access := range accessList {
		name := access.Alias
		if name == "" {
			name = access.Name
		}
		accessMap[name] = access
	}

	AccessMap = accessMap
}

func getRequestMap() {
	requestMap := make(map[string]Request)

	var requestList []Request
	g.DB().Model(config.TableRequest).Scan(&requestList)

	for _, item := range requestList {

		tag := item.Tag
		if strings.ToLower(tag) != tag {
			// 本身大写, 如果没有外层, 则套一层
			// https://github.com/Tencent/APIJSON/issues/115#issuecomment-565733254
			if _, ok := item.Structure[tag]; !ok {
				item.Structure = g.Map{
					tag: item.Structure,
				}
			}
		}

		requestMap[item.Method+"@"+item.Tag] = item
	}

	RequestMap = requestMap
}
