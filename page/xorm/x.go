package xorm

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/no-bibi/bugs/fun"
	"github.com/no-bibi/bugs/page"
	"github.com/no-bibi/bugs/page/options"
	"math"
	"strings"
)

type Xorm struct {
	DB    interface{}   //各个库的db
	Query []interface{} //查询参数
	opt   *options.Options
}

func (this *Xorm) New(db interface{}) page.Page {
	obj := &Xorm{
		DB:    db,
		Query: make([]interface{}, 0),
	}
	return obj
}

func (this *Xorm) Opt(opt *options.Options) page.Page {
	this.opt = opt
	return this
}

func (this *Xorm) Where(query interface{}, args ...interface{}) page.Page {

	this.Query = append(this.Query, query)
	this.Query = append(this.Query, args...)
	return this
}

func (this *Xorm) Page(data interface{}) page.Result {

	var session *xorm.Session

	switch this.DB.(type) {
	case *xorm.EngineGroup:
		session = this.DB.(*xorm.EngineGroup).Where(this.Query[0], this.Query[1:]...)
	case *xorm.Session:
		session = this.DB.(*xorm.Session).Where(this.Query[0], this.Query[1:]...)
	default:
		panic(`db type is not support`)
	}

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
	data = fun.MakeClone(data)

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
