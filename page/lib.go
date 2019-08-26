package page

import (
	"github.com/no-bibi/bugs/page/options"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Page)
)

// Register curd
func Register(name string, driver Page) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("curd: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("curd: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func New(name string, db interface{}) Page {
	driversMu.RLock()
	defer driversMu.RUnlock()
	driver, ok := drivers[name]
	if !ok {
		panic("curd: not  driver " + name)
	}

	return driver.New(db)
}

type Page interface {
	New(db interface{}) Page
	Opt(opt *options.Options) Page                     //分页大小，页数，排序
	Where(query interface{}, args ...interface{}) Page //设置查询参数
	Select(columns ...string) Page                     //查询字段
	Page(data interface{}) Result                      //执行结果
}

//分页返回
type Result struct {
	Count       int         `json:"count"`
	TotalPage   int         `json:"total_page"`
	CurrentPage int         `json:"current_page"`
	Limit       int         `json:"limit"`
	Records     interface{} `json:"records"`
}

func init() {

}
