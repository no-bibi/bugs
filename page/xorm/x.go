package xorm

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/no-bibi/bugs/page"
	"github.com/no-bibi/bugs/page/options"
	"log"
	"math"
	"reflect"
	"strings"
)

type Xorm struct {
	DB    interface{}   //各个库的db
	Query []interface{} //查询参数
	opt   *options.Options
}

func (this *Xorm) New(db interface{}) page.Page {
	this.DB = db
	return this
}

func (this *Xorm) Opt(opt *options.Options) page.Page {
	this.opt = opt
	return this
}

func (this *Xorm) Where(query interface{}, args ...interface{}) page.Page {

	q := make([]interface{}, 0)
	q = append(q, query)
	q = append(q, args...)
	this.Query = q
	return this
}

func (this *Xorm) Page(data interface{}) page.Result {
	db := this.DB.(*xorm.EngineGroup)

	session := db.Where(this.Query[0], this.Query[1:]...)
	if len(this.opt.OrderBy) > 0 {
		for _, o := range this.opt.OrderBy {

			//处理排序字符
			if strings.Index(o, "-") == 0 {
				o = fmt.Sprintf("`%s` desc", strings.TrimLeft(o, "-"))
			} else {
				o = fmt.Sprintf("`%s`", o)
			}
			session.OrderBy(o)
		}
	}

	var (
		count int64
		err   error
	)
	if count, err = session.Limit(this.opt.Limit, this.opt.Offset).FindAndCount(data); err != nil {
		panic(err)
	}

	//always make sure data is [] not null
	arr, ok := data.([]interface{})
	if ok {
		if len(arr) == 0 {
			data = make([]interface{}, 0)
		}
	} else {
		log.Println(reflect.TypeOf(data))
	}

	return page.Result{
		Count:       int(count),
		Records:     data,
		CurrentPage: this.opt.Page,
		Limit:       this.opt.Limit,
		TotalPage:   int(math.Ceil(float64(count) / float64(this.opt.Limit))),
	}

}

func init() {
	page.Register(`xorm`, &Xorm{})
}
