package executor

import (
	"fmt"
	"sort"

	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
)

type BuildExecutionBoneRequest struct {
	Steps      *[]po.Step
	StepGroups *[]po.StepGroup
}

// BuildExecutionBone  根据DB中的数据创建一个ExecutionBone 类型解析
// 本期先解决一层流程，不考虑嵌套流程
func BuildExecutionBone(req BuildExecutionBoneRequest, executor *Executor) (*ExecutionBone, error) {
	var dbSteps = *req.Steps
	// var dbStepGroups = *req.StepGroups
	//
	var dbStepMap = map[string]*po.Step{}

	// levelconnector  是子流程连接的节点，主要是Parallel/Map 节点
	var levelconnector = []*po.Step{}
	// 构建bone
	// GroupBoneMap  group 组内的bone map，
	// 两层map结构， 第一层是 group_id， 第二层是stepname ,都可以从DB钟获得
	// 通过初始化， GroupBoneMap 按照分层存储所有的节点的Bone。 节点都是 StateBone 类型
	var GroupBoneMap = map[int]map[string]StepBone{}

	// StepIdMap  按照step_id 生成的 Map
	var StepIdMap = map[int]Step{}

	for idx := range dbSteps {
		dbStepPtr := &dbSteps[idx]
		// 存成map
		dbStepMap[dbStepPtr.Name] = dbStepPtr
		var groupId int
		newstep, err := NewStepFromData(dbStepPtr, executor)
		if err != nil {
			return nil, err
		}
		// 如果GroupBoneMap 中新的GroupID 不存在， 则创建该GroupID的map， 存入该GroupID 子流程下所有的节点。
		newbone := newstep.GetBone()
		groupId = dbStepPtr.GroupID
		GroupBone, ok := GroupBoneMap[groupId]
		if !ok {
			newSsb := map[string]StepBone{}
			newSsb[dbStepPtr.Name] = newbone
			GroupBoneMap[groupId] = newSsb
		} else {
			GroupBone[dbStepPtr.Name] = newbone
		}
		StepIdMap[dbStepPtr.ID] = newstep
		// 如果是parallel 或者map 节点， 记录下来
		if dbStepPtr.Type == string(states.StateTypes.Map) || dbStepPtr.Type == string(states.StateTypes.Parallel) {
			levelconnector = append(levelconnector, dbStepPtr)
		}
	}

	// 生成顶层的ExecutionBone
	StartGroupID := states.StartGroupID
	toplevelStateBone := &ExecutionBone{
		StartAt: "",
		States:  GroupBoneMap[StartGroupID],
	}
	return toplevelStateBone, nil

}

// BuildExecutionBoneBackup  根据DB中的数据创建一个ExecutionBone 类型解析
func BuildExecutionBoneBackup(req *BuildExecutionBoneRequest, executor *Executor) (*ExecutionBone, error) {

	var dbSteps = *req.Steps
	var dbStepGroups = *req.StepGroups
	//
	var dbStepMap = map[string]*po.Step{}

	// levelconnector  是子流程连接的节点，主要是Parallel/Map 节点
	var levelconnector = []*po.Step{}
	// 构建bone
	// GroupBoneMap  group 组内的bone map，
	// 两层map结构， 第一层是 group_id， 第二层是stepname ,都可以从DB钟获得
	// 通过初始化， GroupBoneMap 按照分层存储所有的节点的Bone。 节点都是 StateBone 类型
	var GroupBoneMap = map[int]map[string]StepBone{}

	// StepIdMap  按照step_id 生成的 Map
	var StepIdMap = map[int]Step{}

	for idx := range dbSteps {
		dbStepPtr := &dbSteps[idx]
		// 存成map
		dbStepMap[dbStepPtr.Name] = dbStepPtr
		var groupId int
		newstep, err := NewStepFromData(dbStepPtr, executor)
		if err != nil {
			return nil, err
		}
		// 如果GroupBoneMap 中新的GroupID 不存在， 则创建该GroupID的map， 存入该GroupID 子流程下所有的节点。
		newbone := newstep.GetBone()
		groupId = dbStepPtr.GroupID
		GroupBone, ok := GroupBoneMap[groupId]
		if !ok {
			newSsb := map[string]StepBone{}
			newSsb[dbStepPtr.Name] = newbone
			GroupBoneMap[groupId] = newSsb
		} else {
			GroupBone[dbStepPtr.Name] = newbone
		}
		StepIdMap[dbStepPtr.ID] = newstep
		// 如果是parallel 或者map 节点， 记录下来
		if dbStepPtr.Type == string(states.StateTypes.Map) || dbStepPtr.Type == string(states.StateTypes.Parallel) {
			levelconnector = append(levelconnector, dbStepPtr)
		}
	}

	// 处理StepGroup, 把StepGroup 按照step_id 分成不同的分组,  方便计算出一个State下有多少个子流程
	// sgMap  StepGroup 按照start_id 进行分组
	var sgMap = map[int][]*po.StepGroup{}
	var step_id int
	for index := range dbStepGroups {
		step_id = dbStepGroups[index].StepID
		if sgg, ok := sgMap[step_id]; ok {
			sgMap[step_id] = append(sgg, &(dbStepGroups[index]))
		} else {
			sgMap[step_id] = []*po.StepGroup{&(dbStepGroups[index])}
		}
	}
	// 排序， 每个 子流程内的多个group 按照index 进行排序
	for _, arraygroup := range sgMap {
		sort.SliceStable(arraygroup, func(i, j int) bool { return arraygroup[i].GroupIndex < arraygroup[j].GroupIndex })
	}
	// 处理连接器节点，把 GroupBoneMap 中的不同层级的节点连起来。
	for _, lc := range levelconnector {

		// 处理Parallel/Map 类型的连接器
		// parallel state process
		sb := GroupBoneMap[lc.GroupID][lc.Name]
		sb.Branches = []StepBone{}
		sg, ok := sgMap[lc.ID]
		if !ok {
			// 当前还没有节点数据， 说明还没有执行到这里，Map类型生成虚拟节点
			// statemachinebone
			// TODO， 这里需要重新设计。， 对于Map类型的LevelConnecter识别。
			// vstate := StateMap[lc.ID]
			// // 判断是Map 类型
			// if vstate.GetBone().Type == grammer.StateType.Map {
			// 	subsm := vstate.(*Map).State.GetIterator().GetBone()
			// 	vbone := TransformVritualBone(subsm)
			// 	sb.SubGroup = append(sb.SubGroup, vbone)
			// 	GroupBoneMap[lc.GroupID][lc.Name] = sb
			// }
			continue
		}

		for _, subgroup := range sg {
			subgroupId := subgroup.SubGroupID
			startatnode, ok := dbStepMap[subgroup.StartAt]
			if !ok {
				return nil, fmt.Errorf("state data error: startat state not found ")
			}
			newSeb := StepBone{
				ExecutionBone: &ExecutionBone{
					StartAt: startatnode.Name,
					States:  GroupBoneMap[subgroupId],
				},
			}
			sb.Branches = append(sb.Branches, newSeb)
		}

		GroupBoneMap[lc.GroupID][lc.Name] = sb
		continue

	}

	// 生成顶层的ExecutionBone
	StartGroupID := states.StartGroupID
	toplevelStateBone := &ExecutionBone{
		StartAt: "",
		States:  GroupBoneMap[StartGroupID],
	}
	return toplevelStateBone, nil

}
