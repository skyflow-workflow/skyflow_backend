package cache

import "github.com/skyflow-workflow/skyflow_backbend/pkg/cachemap"

// taskCacheService 任务缓存服务的具体实现
// 使用 cachemap.CacheMap 作为底层存储
type taskCacheService struct {
	CacheMap *cachemap.CacheMap
}

func NewTaskCacheService(cm *cachemap.CacheMap) TaskCacheService {
	return &taskCacheService{
		CacheMap: cm,
	}
}

func (s *taskCacheService) Get(key string) (int, bool) {

	value := s.CacheMap.Pop(key)
	if value == cachemap.ZeroValue {
		return 0, false
	}
	return value, true
}

func (s *taskCacheService) Refresh(data map[string][]int) error {
	return s.CacheMap.Refresh(data)
}
