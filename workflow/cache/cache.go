package cache

// TaskCacheService 任务缓存服务接口
// 提供任务相关的缓存操作功能，包括任务数据的弹出和刷新
type TaskCacheService interface {
	// Get 从缓存中弹出指定TaskResource对应的任务ID
	// 参数:
	//   - key: TaskResource，用于标识特定的任务类型
	// 返回值:
	//   - int: 弹出的任务ID，如果键不存在则返回0
	//   - bool: 表示是否成功弹出，true表示成功，false表示键不存在
	Get(key string) (int, bool)

	// Refresh 刷新缓存数据，用新的数据替换现有缓存
	// 参数:
	//   - data: 新的缓存数据映射，键为TaskResource名字，值为任务ID列表
	// 说明:
	//   - 此操作会完全替换现有缓存内容
	//   - 适用于批量更新缓存场景
	Refresh(data map[string][]int) error
}
