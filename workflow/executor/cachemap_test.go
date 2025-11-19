package executor

import (
	"fmt"
	"testing"
)

func TestCachemap(t *testing.T) {

	// myExecutionService.RefreshCacheMapOnce()
	cm := myExecutionService.CacheService
	fmt.Println(cm)

}
