// 加锁服务， 可以使用DB /Redis or 其他资源服务完成加锁需求
package lock

import "github.com/mmtbak/microlibrary/rdb"

// LockService lock service .
type LockService interface {
	LockExecution(int) Lock
	LockStep(int) Lock
}

type Lock interface {
	Lock() error
	Unlock()
	GetTx() rdb.Tx
}

var LockTypes = struct {
	Execution string
	Step      string
}{
	Execution: "Execution",
	Step:      "Step",
}
