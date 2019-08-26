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
	db interface{} //各个库的db session
	*options.Options
}

func (this *Xorm) New(db interface{}) page.Page {
	switch this.db.(type) {
	case *xorm.EngineGroup:
		this.db = this.db.(*xorm.EngineGroup).NewSession()
	case *xorm.Engine:
		this.db = this.db.(*xorm.Engine).NewSession()
	case *xorm.Session:
		this.db = db
	default:
		panic(`db type is not support`)
	}

	obj := &Xorm{
		db: db,
	}
	return obj
}

func (this *Xorm) Opt(opt *options.Options) page.Page {

	if len(opt.OrderBy) > 0 {
		for _, o := range opt.OrderBy {
			//处理排序字符
			if strings.Index(o, "-") == 0 {
				o = fmt.Sprintf("`%s` desc", strings.TrimLeft(o, "-"))
			} else {
				o = fmt.Sprintf("`%s`", o)
			}
			this.db.(*xorm.Session).OrderBy(o)
		}
	}

	this.db.(*xorm.Session).Limit(opt.Limit, opt.Offset)
	this.Options = opt
	return this
}

func (this *Xorm) Where(query interface{}, args ...interface{}) page.Page {
	this.db = this.db.(*xorm.Session).Where(query, args...)
	return this
}

func (this *Xorm) Select(columns ...string) page.Page {
	this.db = this.db.(*xorm.Session).Cols(columns...)
	return this
}

func (this *Xorm) Page(data interface{}) page.Result {

	var (
		count int64
		err   error
	)
	if count, err = this.db.(*xorm.Session).FindAndCount(data); err != nil {
		panic(err)
	}

	//always make sure data is [] not null
	data = fun.MakeClone(data)

	return page.Result{
		Count:       int(count),
		Records:     data,
		CurrentPage: this.Options.Page,
		Limit:       this.Limit,
		TotalPage:   int(math.Ceil(float64(count) / float64(this.Limit))),
	}

}

func init() {
	page.Register(`xorm`, &Xorm{})
}
