package executor

/*
 * @Author: mumangtao@gmail.com
 * @Date: 2020-08-06 17:29:43
 * @Last Modified by: mumangtao@gmail.com
 * @Last Modified time: 2020-08-06 17:33:49
 */

import (
	"log/slog"
	"time"

	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
)

// RefreshCacheMap refresh cache
func (svc *executionService) RefreshCacheMap() {

	var err error
	for {
		err = svc.RefreshCacheMapOnce()
		if err != nil {
			slog.Error("refresh execution service cachemap failed ", "error", err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

func (svc *executionService) RefreshCacheMapOnce() error {
	var err error
	var cm = svc.CacheService
	// QueryResult 查询结果
	type QueryResult struct {
		ID       int    `json:"id"`
		Resource string `json:"resource"`
	}
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	var results = []QueryResult{}

	err = tx.Model(new(po.ActivityTask)).Select("id", "resource").Order("id asc").Find(&results).Error
	if err != nil {
		return err
	}
	var urimap = map[string][]int{}
	for _, item := range results {
		if v, ok := urimap[item.Resource]; ok {
			v = append(v, item.ID)
			urimap[item.Resource] = v
		} else {
			urimap[item.Resource] = []int{item.ID}
		}
	}
	err = cm.Refresh(urimap)
	if err != nil {
		return err
	}
	return nil
}
