package exporter

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/coocood/freecache"
	"github.com/mmtbak/microlibrary/rdb"
	"github.com/panjf2000/ants/v2"
	"github.com/skyflow-workflow/skyflow_backbend/config"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

const (
	// DefaultPoolSize default pool size
	DefaultPoolSize = 1000
	// EventContentPrefix event content prefix
	EventContentPrefix = "EventContent_"
)

type ExporterService = *exporterService

type exporterService struct {
	PoolSize   int
	WorkerPool *ants.PoolWithFunc
	MetaDB     *rdb.DBClient
	// DBListener is special listener, it can query history events
	DBListener *DBListener
	Listeners  []Listener
	// Execution Info Cache
	ExecutionCache *freecache.Cache
	config         *config.ExporterConfig
}

// NewExporterService create a new exporter instance
func NewExporterService(dbClient *rdb.DBClient, conf *config.ExporterConfig) (*exporterService, error) {
	var err error
	dbListener := NewDBListener(dbClient)
	exporter := &exporterService{
		MetaDB:     dbClient,
		WorkerPool: nil,
		DBListener: dbListener,
		Listeners:  []Listener{dbListener},
		config:     conf,
	}
	cache := freecache.NewCache(conf.CacheSizeMB * 1024 * 1024)
	exporter.ExecutionCache = cache

	// 创建异步执行
	workerPool, err := ants.NewPoolWithFunc(DefaultPoolSize, exporter.AsyncSendExecutionEvents, ants.WithNonblocking(false))
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
	if len(events) == 0 {
		return
	}
	svc.FullFillExecutionEvent(&events)
	err := svc.WorkerPool.Invoke(events)
	if err != nil {
		slog.Error(err.Error())
	}
}

// AsyncSendExecutionEvents async send execution events
func (svc *exporterService) AsyncSendExecutionEvents(i interface{}) {
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
	tx := svc.MetaDB.NewTx()
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
	tx := svc.MetaDB.NewTx()
	defer tx.Commit()
	// 查询总数
	tx = tx.Model(new(po.ExecutionEvent)).Where(po.ExecutionEvent{StepID: req.StepID})
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

// SendExecutionEvents 发送event
func (svc *exporterService) FullFillExecutionEvent(events *[]vo.ExecutionEvent) {

	for idx := range *events {

		// 不能是空类型
		if (*events)[idx].Data == nil {
			slog.Error("event content data  is nil")
			continue
		}
		exeID := (*events)[idx].ExecutionID

		// 根据data 获得 event类型和event值
		ot := reflect.TypeOf((*events)[idx].Data).Name()
		t := strings.TrimPrefix(ot, EventContentPrefix)
		exeInfo := svc.FindCacheExecution(exeID)
		(*events)[idx].ExecutionUUID = exeInfo.UUID
		(*events)[idx].ExecutionURI = exeInfo.URI
		(*events)[idx].EventType = t
		// (*events)[idx].NanoSeconds = fmt.Sprintf("%020d", time.Now().UnixNano())
		(*events)[idx].NanoSeconds = time.Now().UnixNano()
	}

}

// ExecutionCacheInfo  缓存的Execution内容
type CacheExecutionInfo struct {
	URI         string
	UUID        string
	ExecutionID int
}

// FindCacheExecution findexecution uri
func (svc *exporterService) FindCacheExecution(execution_id int) CacheExecutionInfo {

	var err error
	var exeInfo CacheExecutionInfo
	key := []byte(fmt.Sprintf("execution_%d", execution_id))
	data, err := svc.ExecutionCache.Get(key)
	if err == nil {
		// 找到了， 解析数据
		err = json.Unmarshal(data, &exeInfo)
		if err != nil {
			return exeInfo
		}
		return exeInfo
	}
	var dbExe po.Execution
	tx := svc.MetaDB.NewTx()
	defer tx.Commit()
	err = tx.Where(po.Execution{ID: execution_id}).Select("id", "uuid", "uri").Take(&dbExe).Error
	if err != nil {
		return exeInfo
	}
	exeInfo = CacheExecutionInfo{
		URI:         dbExe.URI,
		UUID:        dbExe.UUID,
		ExecutionID: dbExe.ID,
	}
	byteExeInfo, _ := json.Marshal(exeInfo)
	svc.ExecutionCache.Set(key, byteExeInfo, svc.config.CacheTTLSecond)
	return exeInfo
}
