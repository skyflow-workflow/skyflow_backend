package exporter

import (
	"log/slog"

	"github.com/panjf2000/ants/v2"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

const (
	// DefaultPoolSize default pool size
	DefaultPoolSize = 1000
)

type ExporterService = *exporterService

type exporterService struct {
	PoolSize   int
	WorkerPool *ants.PoolWithFunc
	// DBListener is special listener, it can query history events
	DBListener *DBListener
	Listeners  []Listener
}

// NewExporterService create a new exporter instance
func NewExporterService(lis *DBListener) (*exporterService, error) {
	var err error
	exporter := &exporterService{
		WorkerPool: nil,
		DBListener: lis,
		Listeners:  []Listener{},
	}
	workerPool, err := ants.NewPoolWithFunc(DefaultPoolSize, exporter.AyncSendExecutionEvents, ants.WithNonblocking(false))
	if err != nil {
		return nil, err
	}
	exporter.WorkerPool = workerPool
	return exporter, nil
}

// AddListener add a event listener
func (svc *exporterService) AddListener(listener Listener) {
	svc.Listeners = append(svc.Listeners, listener)
}

func (svc *exporterService) SyncSchema() error {
	var err error
	if svc.DBListener != nil {
		err = svc.DBListener.SyncSchema()
		return err
	}
	for _, lis := range svc.Listeners {
		err = lis.SyncSchema()
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *exporterService) SendExecutionEvents(events []vo.ExecutionEvent) {

	err := svc.WorkerPool.Invoke(events)
	if err != nil {
		slog.Error(err.Error())
	}
}

// AyncSendExecutionEvents async send execution events
func (svc *exporterService) AyncSendExecutionEvents(i interface{}) {
	var events, ok = i.([]vo.ExecutionEvent)
	if !ok {
		return
	}

	for _, lis := range svc.Listeners {
		lis.SendEvents(events)
	}
}

// ListExecutionEvents
func (svc *exporterService) ListExecutionEvents(req vo.ListExecutionEventsRequest) (vo.ListExecutionEventsResponse, error) {

	var err error
	var count int64
	var events = []po.ExecutionEvent{}
	limit, offset := req.PageRequest.Limit()
	// 新建事务
	tx := svc.DBListener.Client().NewTx()
	defer tx.Commit()
	// 查询总数
	tx = tx.Model(new(po.ExecutionEvent)).Where(po.ExecutionEvent{ExecutionID: req.ExecutionID})
	err = tx.Count(&count).Error
	if err != nil {
		return vo.ListExecutionEventsResponse{}, err
	}
	// 查询数据
	err = tx.Limit(limit).Offset(offset).Order("nano_seconds asc").Find(&events).Error
	if err != nil {
		return vo.ListExecutionEventsResponse{}, err
	}
	resp := vo.ListExecutionEventsResponse{
		Events:       events,
		PageResponse: req.PageRequest.Response(count),
	}
	return resp, nil
}

func (svc *exporterService) ListStepEvents(req vo.ListStepEventsRequest) (vo.ListExecutionEventsResponse, error) {
	var count int64
	var events = []po.ExecutionEvent{}
	limit, offset := req.PageRequest.Limit()
	// 新建事务
	tx := svc.DBListener.Client().NewTx()
	defer tx.Commit()
	// 查询总数
	tx = tx.Model(new(po.ExecutionEvent)).Where(po.ExecutionEvent{StateID: req.StepID})
	err := tx.Count(&count).Error
	if err != nil {
		return vo.ListExecutionEventsResponse{}, err
	}
	// 查询数据
	err = tx.Limit(limit).Offset(offset).Order("nano_seconds asc").Find(&events).Error
	if err != nil {
		return vo.ListExecutionEventsResponse{}, err
	}

	resp := vo.ListExecutionEventsResponse{
		Events:       events,
		PageResponse: req.PageRequest.Response(count),
	}
	return resp, nil
}
