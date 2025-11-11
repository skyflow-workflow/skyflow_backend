package dispatcher

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/goodaye/wire"
	"github.com/mmtbak/microlibrary/limiter"
	"github.com/panjf2000/ants/v2"
	"github.com/skyflow-workflow/skyflow_backbend/workflow"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/executor"
)

// DispatcherService  message dispatcher
type DispatcherService struct {
	wire.BaseService
	workflowService  workflow.WorkflowService
	ExecutionService executor.ExecutionService
	ctx              context.Context
	cancelfunc       context.CancelFunc
	// 消息接收器的waitgroup
	receiverWg sync.WaitGroup
	// 用于等待所有worker完成
	eventWg    sync.WaitGroup
	workerPool *ants.PoolWithFunc
	option     Option
	// event counter for performance
	// 记录已经处理过的事件数量
	EventCounter uint64
	// 内存限制器
	memlimiter *limiter.MemoryLimiter
}

// Option dispatcher option
type Option struct {
	Concurrency int
	Debug       bool
}

// DefaultOption default option
var DefaultOption = Option{
	Concurrency: 100,
	Debug:       false,
}

// NewDispatcher create a new dispatcher
func NewDispatcher(workflowSvc workflow.WorkflowService, option Option) (*DispatcherService, error) {

	var err error
	// check config reasonability
	if option.Concurrency > workflowSvc.DBClient.GetConfig().MaxOpenConns {
		err = fmt.Errorf("dispatcher config invalid, concurrency should less than MetaDB maxopenconn")
		return nil, err
	}

	// Context
	ctx, cancel := context.WithCancel(context.Background())

	dispather := &DispatcherService{
		workflowService:  workflowSvc,
		ExecutionService: workflowSvc.ExecutionService,
		receiverWg:       sync.WaitGroup{},
		eventWg:          sync.WaitGroup{},
		ctx:              ctx,
		cancelfunc:       cancel,
		option:           option,
		EventCounter:     0,
	}
	return dispather, nil
}

// Start start schedular worker
func (svc *DispatcherService) Start() error {
	slog.Info("start run dispatcher ")
	//
	var err error
	var workerPool *ants.PoolWithFunc
	// Pool
	workerPool, err = ants.NewPoolWithFunc(svc.option.Concurrency, svc.ProcessMessage, ants.WithNonblocking(false))
	if err != nil {
		return err
	}
	svc.workerPool = workerPool

	// 启动Limit
	if svc.memlimiter != nil {
		svc.memlimiter.Start()
	}

	svc.StartSchedularWorkerManager()

	slog.Info("start run dispatcher success")
	return nil
}

// Stop return stopfinish chan
func (svc *DispatcherService) Stop() error {
	slog.Info("try to stop dispatcher ")
	svc.cancelfunc()
	svc.receiverWg.Wait()
	slog.Info("stop event receiver success")
	svc.eventWg.Wait()
	slog.Info("stop event worker success")

	// 关闭Limiter
	svc.memlimiter.Close()
	svc.workerPool.Release()
	slog.Info("stop dispatcher success")
	return nil
}

// IncreaseEventCounter increase event counter
// 性能计数器， 用于统计事件数量
func (svc *DispatcherService) IncreaseEventCounter() {
	atomic.AddUint64(&svc.EventCounter, 1)
}

// GetEventCounter get event counter
func (svc *DispatcherService) GetEventCounter() uint64 {
	return atomic.LoadUint64(&svc.EventCounter)
}

// SetLimiter set memory limiter
func (svc *DispatcherService) SetLimiter(limiter *limiter.MemoryLimiter) {
	svc.memlimiter = limiter
}
