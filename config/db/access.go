package db

import (
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/samber/lo"
	"net/http"
	"strings"
)

var accessMap = map[string]Access{}

type FieldsGetValue struct {
	In  map[string][]string
	Out map[string]string
}

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

	RowKeyGen string // 主键生成策略
	RowKey    string
	FieldsGet map[string]FieldsGetValue
	Executor  string
}

func (a *Access) GetFieldsGetOutByRole(role string) []string {
	var fieldsMap map[string]string

	if val, exists := a.FieldsGet[role]; exists {
		fieldsMap = val.Out
	} else {
		fieldsMap = a.FieldsGet["default"].Out
	}
	return lo.Keys(fieldsMap)
}

func (a *Access) GetFieldsGetInByRole(role string) map[string][]string {
	var inFieldsMap map[string][]string

	if val, exists := a.FieldsGet[role]; exists {
		inFieldsMap = val.In
	} else {
		inFieldsMap = a.FieldsGet["default"].In
	}

	return inFieldsMap
}

func loadAccessMap() {
	_accessMap := make(map[string]Access)

	var accessList []Access

	db := g.DB()

	err := db.Model(config.TableAccess).Scan(&accessList)
	if err != nil {
		panic(err)
	}

	type AccessExt struct {
		RowKey    string
		FieldsGet map[string]FieldsGetValue
	}

	for _, access := range accessList {
		name := access.Alias
		if name == "" {
			name = access.Name
		}

		if len(access.Get) > 0 {
			access.Get = strings.Split(access.Get[0], ",")
		}
		if len(access.Head) > 0 {
			access.Head = strings.Split(access.Head[0], ",")
		}
		if len(access.Gets) > 0 {
			access.Gets = strings.Split(access.Gets[0], ",")
		}
		if len(access.Heads) > 0 {
			access.Heads = strings.Split(access.Heads[0], ",")
		}
		if len(access.Post) > 0 {
			access.Post = strings.Split(access.Post[0], ",")
		}
		if len(access.Put) > 0 {
			access.Put = strings.Split(access.Put[0], ",")
		}
		if len(access.Delete) > 0 {
			access.Delete = strings.Split(access.Delete[0], ",")
		}

		_accessMap[name] = access
	}

	accessMap = _accessMap
}

func GetAccess(tableAlias string, accessVerify bool) (*Access, error) {
	tableAlias, _ = util.ParseNodeKey(tableAlias)
	access, ok := accessMap[tableAlias]

	if !ok {
		if accessVerify {
			return nil, gerror.Newf("access[%s]: 404", tableAlias)
		}
		return &Access{
			Debug: 0,
			Name:  tableAlias,
			Alias: tableAlias,
		}, nil
	}

	return &access, nil
}

func GetAccessRole(table string, method string) ([]string, string, error) {
	access, ok := accessMap[table]

	if !ok {
		return nil, "", gerror.Newf("access[%s]: 404", table)
	}

	switch method {
	case http.MethodGet:
		return access.Get, access.Name, nil
	case http.MethodHead:
		return access.Head, access.Name, nil
	case http.MethodPost:
		return access.Post, access.Name, nil
	case http.MethodPut:
		return access.Put, access.Name, nil
	case http.MethodDelete:
		return access.Delete, access.Name, nil
	}

	return []string{}, access.Name, nil
}

func Init() {

	initTable()

	Reload()
}

// Reload 重载刷新配置
func Reload() {
	loadAccessMap()
	loadRequestMap()
	loadTableMeta()
}

// initTable 暂时先这样吧
func initTable() {
	// maybe need the sql file when not use table-sync

}
