package dbm

import (
	"sync"
)

type IDriver interface {
	Go2Db(t string) string
	Db2Go(t string) string
	ToSql(tab *Table) string
}

var driverMap = map[string]IDriver{}
var driverLock sync.RWMutex

func Register(driver string, parser IDriver) {
	driverLock.Lock()
	defer driverLock.Unlock()
	driverMap[driver] = parser
}

func GetDriver(driver string) IDriver {
	driverLock.RLock()
	defer driverLock.RUnlock()
	return driverMap[driver]
}
