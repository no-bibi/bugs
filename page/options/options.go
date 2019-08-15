package options

import (
	"net/http"
	"strconv"
	"strings"
)

// Param 分页参数
type Options struct {
	Page    int
	Limit   int
	Offset  int
	OrderBy []string
}

// Param 分页参数
type Setting struct {
	DefaultLimit int
	PageField    string
	LimitField   string
	OrderField   string
}

var setting = &Setting{12, `page`, `limit`, `sort`}

func Change(s *Setting) {
	setting = s
}

//生成分页选项
func New(r *http.Request) *Options {

	params := r.URL.Query()
	o := &Options{}

	o.Page, _ = strconv.Atoi(params.Get(setting.PageField))
	if o.Page < 1 {
		o.Page = 1
	}

	o.Limit, _ = strconv.Atoi(params.Get(setting.LimitField))
	if o.Limit < 1 {
		o.Limit = setting.DefaultLimit
	}

	if o.Page != 1 {
		o.Offset = (o.Page - 1) * o.Limit
	}

	sort := params.Get(setting.OrderField)
	if len(sort) > 0 {
		o.OrderBy = strings.Split(sort, ",")
	}

	return o
}
