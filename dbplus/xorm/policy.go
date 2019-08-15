package xorm

import (
	"github.com/go-xorm/xorm"
	"sync"
	"time"
)

type NotWork struct {
	List map[int]bool
	sync.Mutex
}

//轮询方案：保证可用性 + 自动恢复 + 无从读主
func AliveRoundRobinPolicy() xorm.GroupPolicyHandler {
	var pos = -1
	var one = sync.Once{}
	var x *xorm.EngineGroup

	var done = make(chan bool, 0)
	notWorking := NotWork{List: make(map[int]bool)}
	go func() {
		<-done
		for {
			notWorking.Lock()
			for index, s := range x.Slaves() {
				if err := s.Ping(); err != nil {
					notWorking.List[index] = true
				} else {
					delete(notWorking.List, index)
				}
			}
			notWorking.Unlock()
			time.Sleep(time.Second * 3)
		}

	}()

	return func(group *xorm.EngineGroup) *xorm.Engine {

		var slaves = group.Slaves()

		one.Do(func() {
			x = group
			done <- true
		})

		notWorking.Lock()
		defer notWorking.Unlock()
		counter := 0

		for {
			pos++
			counter++

			if pos >= len(slaves) {
				pos = 0
			}

			if len(notWorking.List) == len(slaves) {
				return group.Engine
			}

			_, exist := notWorking.List[pos]
			if !exist || counter >= len(slaves) {
				break
			}
		}
		return slaves[pos]
	}
}
