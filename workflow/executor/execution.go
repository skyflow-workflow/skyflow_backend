package executor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"sort"
	"time"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// Execution  StateMachine 执行实例
type Execution struct {
	StateMachine *states.StateMachine
	Data         *po.Execution
	Executor     *Executor
	// state map
	States           map[string]Step
	ExecutionService ExecutionService
}

// ExecutionInfo  内置的ExecutionInfo Struct
type ExecutionInfo struct {
	UUID    string `json:"uuid"`    // execution uuid
	URLPath string `json:"urlpath"` // 访问路径
	URL     string `json:"url"`     // 完整访问地址

}

// NewExecutionFromID Create Execution by execution id
func NewExecutionFromID(id int, svc ExecutionService) (*Execution, error) {

	var err error
	var dbExe *po.Execution

	// 使用最小的数据集来来初始化
	dbExe, err = svc.QueryExecutionByID(id, ExecutionFields.L1, nil)
	if err != nil {
		return nil, err
	}

	exe, err := NewExecutionFromData(dbExe, svc)
	return exe, err
}

// NewExecutionFromData  根据DB中的数据创建一个ServiceExecution 类型解析
/* data : po.Execution 类型， Execution 信息
* datastates : Execution 下所有的DB节点
* StepGroups : Execution 下所有节点的关联关系.
 */
func NewExecutionFromData(data *po.Execution, svc ExecutionService) (exe *Execution, err error) {

	exe = &Execution{
		Data:             data,
		States:           map[string]Step{},
		ExecutionService: svc,
	}

	// delay init
	// var header = grammer.StateMachineHeader{}
	// if data.Header != "" {
	// 	err = json.Unmarshal([]byte(data.Header), &header)
	// 	if err != nil {
	// 		return
	// 	}

	// } else {
	// 	var sm *grammer.StateMachine
	// 	var wf flow.WorkFlow
	// 	wf, err = parser.ParseWorkflow(data.FlowDefinition)
	// 	if err != nil {
	// 		return
	// 	}
	// 	sm = wf.GetNode()
	// 	header = sm.StateMachineHeader
	// }
	// exe.StateMachineHeader = &header
	return
}

// FullInit 全量初始化
// 从数据库中查询到Execution 的所有信息， 解析StateMachine， 并初始化StateMachine 结构
func (exe *Execution) FullInit() error {
	// 如果已经初始化， 则忽略
	var err error
	if exe.StateMachine != nil {
		return nil
	}

	var dbExecution *po.Execution
	id := exe.Data.ID

	// 查询到全量的数据
	dbExecution, err = exe.ExecutionService.QueryExecutionByID(id, []string{}, nil)

	if err != nil {
		return err
	}
	// 后面需要计算input
	exe.Data = dbExecution

	sm, err := parser.ParseStateMachine(exe.Data.Definition)
	if err != nil {
		return err
	}

	exe.StateMachine = sm

	return nil
}

// ProcessEvent 处理事件
func (exe *Execution) ProcessEvent(msgBody queue.InnerMessageBody) error {

	var err error

	lock := exe.ExecutionService.LockService.LockExecution(msgBody.ExecutionID)
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()
	tx := lock.GetTx()

	var dbExe *po.Execution

	dbExe, err = exe.ExecutionService.QueryExecutionByID(msgBody.ExecutionID, []string{"id", "status"}, tx)
	if err != nil {
		return err
	}

	checkstatus, ok := ExecutionEventCheckStatus[msgBody.Type]
	// 如果状态在规则中， 则检查状态，如果不在， 则忽略不检查
	if ok {
		// 如果不在预期状态中， 忽略事件
		if !slices.Contains(checkstatus, dbExe.Status) {
			msg := fmt.Sprintf("execution '%d' process event '%s' current status '%s' not match ",
				msgBody.ExecutionID, msgBody.Type, dbExe.Status)
			slog.Error(msg)
			return nil
		}
	}

	switch msgBody.Type {
	case MessageType.ExecutionInit:
		err = exe.ProcessInit()
	case MessageType.ExecutionTimout:
		err = exe.ProcessTimeout()
	case MessageType.ExecutionAbortTimout:
		err = exe.ProcessAbortTimeout()
	case MessageType.ExecutionFailed:
		err = exe.ProcessAbortTimeout()
	// 不会有 MessageType.ExecutionSucceed: 类型的消息，只有findnext 的时候， 直接调用ProcessSucceed()
	// case MessageType.ExecutionSucceed:
	// 	err = exe.ProcessSucceed()
	default:
		err = fmt.Errorf("unrecognized event '%s' ", msgBody.Type)
	}

	if err != nil {
		err = exe.ExecutionService.ExecutionErrorProcess(err, msgBody)
	}
	return err
}

// GetInput GetInput
// 计算Execution 的input, 计算ExecutionPath 之后 的input
func (e *Execution) GetInput() (interface{}, error) {
	var inputdata, newinputdata interface{}
	var err error

	err = json.Unmarshal([]byte(e.Data.Input), &inputdata)
	if err != nil {
		return nil, err
	}

	fullurl := e.ExecutionService.DomainService.QueryExecutionURL(e.Data.UUID)

	var info = ExecutionInfo{
		UUID: e.Data.UUID,
		URL:  fullurl,
	}
	newinputdata, err = e.StateMachine.GetInput(inputdata, info)
	return newinputdata, err
}

// GetBone 获得execution bone
// NOCC:golint/fnsize("设计如此")
func (e *Execution) GetBone() (ExecutionBone, error) {

	var bone ExecutionBone
	var err error
	var dbSteps []po.Step
	var dbStepGroups []po.StepGroup
	//
	var dbStepMap = map[string]*po.Step{}

	// 查询数据
	tx := e.ExecutionService.MetaDB.DB()
	err = tx.Where(po.Step{ExecutionID: e.Data.ID}).Find(&dbSteps).Error
	if err != nil {
		return bone, err
	}

	err = tx.Where(po.StepGroup{ExecutionID: e.Data.ID}).Find(&dbStepGroups).Error
	if err != nil {
		return bone, err
	}

	// levelconnector  是子流程连接的节点，主要是Parallel/Map 节点
	var levelconnector = []*po.Step{}
	// 构建bone
	StartGroupID := states.StartGroupID
	// GroupBoneMap  group 组内的bone map，
	// 两层map结构， 第一层是 group_id， 第二层是statename ,都可以从DB钟获得
	// 通过初始化， GroupBoneMap 按照分层存储所有的节点的Bone。 节点都是 StateBone 类型
	var GroupBoneMap = map[int]map[string]StepBone{
		// StartGroupID: map[string]StateBone{},
	}
	// StateBoneMap  按照step_id 生成的 Map
	var StateMap = map[int]Step{}

	for di := range dbSteps {
		dp := &dbSteps[di]
		// 存成map
		dbStepMap[dp.Name] = dp
		var groupid int
		newstate, err := NewStepFromData(dp, e.ExecutionService.StandardExecutor)
		if err != nil {
			return bone, err
		}
		// 如果GroupBoneMap 中新的GroupID 不存在， 则创建该GroupID的map， 存入该GroupID 子流程下所有的节点。
		newbone := newstate.GetBone()
		groupid = dp.GroupID
		GroupBone, ok := GroupBoneMap[groupid]
		if !ok {
			newssb := map[string]StepBone{}
			newssb[dp.Name] = newbone
			GroupBoneMap[groupid] = newssb
		} else {
			GroupBone[dp.Name] = newbone
		}
		StateMap[dp.ID] = newstate
		// 如果是parallel 或者map 节点， 记录下来
		if dp.Type == string(states.StateTypes.Map) || dp.Type == string(states.StateTypes.Parallel) {
			levelconnector = append(levelconnector, dp)
		}
	}

	// 处理StepGroup, 把StepGroup 按照step_id 分成不同的分组,  方便计算出一个State下有多少个子流程
	// sgmap  StepGroup 按照start_id 进行分组
	var sgmap = map[int][]*po.StepGroup{}
	var step_id int
	for index := range dbStepGroups {
		step_id = dbStepGroups[index].StepID
		if sgg, ok := sgmap[step_id]; ok {
			sgmap[step_id] = append(sgg, &(dbStepGroups[index]))
		} else {
			sgmap[step_id] = []*po.StepGroup{&(dbStepGroups[index])}
		}
	}
	// 排序， 每个 子流程内的多个group 按照index 进行排序
	for _, arraygroup := range sgmap {
		sort.SliceStable(arraygroup, func(i, j int) bool { return arraygroup[i].GroupIndex < arraygroup[j].GroupIndex })
	}
	// 处理连接器节点，把 GroupBoneMap 中的不同层级的节点连起来。
	for _, lc := range levelconnector {

		// 处理Parallel/Map 类型的连接器
		// parallel state process
		sb := GroupBoneMap[lc.GroupID][lc.Name]
		sb.Branches = []StepBone{}
		sg, ok := sgmap[lc.ID]
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
			subgroupid := subgroup.SubGroupID
			startatnode, ok := dbStepMap[subgroup.StartAt]
			if !ok {
				return bone, fmt.Errorf("state data error: startat state not found ")
			}
			newseb := StepBone{
				StateMachineBone: &StateMachineBone{
					StartAt: startatnode.Name,
					States:  GroupBoneMap[subgroupid],
				},
			}
			sb.Branches = append(sb.Branches, newseb)
		}

		GroupBoneMap[lc.GroupID][lc.Name] = sb
		continue

	}
	toplevelStateBone := ExecutionBone{
		StartAt: e.StateMachine.StartAt,
		States:  GroupBoneMap[StartGroupID],
	}
	return toplevelStateBone, nil
}

// ProcessInit  初始化execution
// 初始化，是将 Execution 解析成多个Step，并且持久化到 po.Step 中。
func (e *Execution) ProcessInit() error {
	var err error
	var dbExe *po.Execution

	err = e.FullInit()
	if err != nil {
		return err
	}

	var sm = e.StateMachine
	dbExe = e.Data

	starttime := time.Now()

	// 计算第一个Execution 经过Header的input
	stateinput, err := e.GetInput()
	if err != nil {
		return err
	}
	stateinputstr, err := toolkit.ToString(stateinput)
	if err != nil {
		return err
	}
	// 获得Header
	header := e.StateMachine.StateMachineHeader

	// 先开启session， 所有的事情都应该在一个session中
	// 开始事务
	tx, maker := e.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	now := time.Now()
	updateExecution := po.Execution{
		Status:    string(ExecutionStatus.Running),
		StartTime: &now,
		Input:     stateinputstr,
		Header:    header.GetDefinition(),
	}

	err = tx.Where(po.Execution{ID: dbExe.ID}).Updates(&updateExecution).Error
	if err != nil {
		return err
	}

	inSmResp, err := e.InsertStateMachine(sm.StateMachineBody, InsertStateMachineOption{
		StartDeIndex: -1,
		StartDepth:   states.StartDepth,
		StartGroupID: states.StartGroupID,
	}, tx)
	if err != nil {
		return err
	}
	updatetask := po.Step{
		Input:  stateinputstr,
		Status: string(StepStatus.WaitInit),
	}

	err = tx.Where(po.Step{ID: inSmResp.StartStepID}).Updates(&updatetask).Error
	if err != nil {
		return err
	}

	tx.Commit()

	TaskCreateTime := time.Now()
	// Create Event
	event1 := vo.ExecutionEvent{
		ExecutionID: dbExe.ID,
		StepName:    ExecutionEventStateName.Start,
		StartTime:   starttime,
		FinishTime:  TaskCreateTime,
		Data: EventContent_ExecutionStart{
			Input: dbExe.Input,
			URI:   dbExe.URI,
		},
	}
	e.ExecutionService.SendExecutionEvents(event1)

	// 发送消息
	// 处理下一个节点
	msg := StepExecuteMessage{
		Block: false,
	}
	message := NewStepMessage(dbExe.ID, MessageType.StateNewTurn, inSmResp.StartStepID, msg)
	err = e.ExecutionService.InnerQueue.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}
	timeout := e.StateMachine.GetTimeout()
	if timeout.Timeout > 0 {
		// 发送超时事件
		message = NewExecutionMessage(dbExe.ID, MessageType.ExecutionTimout, nil)
		eventTime := time.Now().Add(timeout.Timeout)
		err = e.ExecutionService.InnerQueue.SendInnerMessage(message, &eventTime)
		if err != nil {
			return err
		}
	}
	if timeout.AbortTimeout > 0 {
		// 发送Abort事件
		message = NewExecutionMessage(dbExe.ID, MessageType.ExecutionAbortTimout, nil)
		eventTime := time.Now().Add(timeout.AbortTimeout)
		err = e.ExecutionService.InnerQueue.SendInnerMessage(message, &eventTime)
		if err != nil {
			return err
		}
	}
	return nil
}

// InsertStateMachine 将StateMachine包含的Step 插入到数据库中
// 根据StateMachineBody 生成Steps, 并且插入到数据库中
func (e *Execution) InsertStateMachine(smb *states.StateMachineBody, opt InsertStateMachineOption,
	tx rdb.Tx) (resp InsertStateMachineResponse, err error) {

	var deIndex = opt.StartDeIndex
	// offsetGroupId GroupID偏移量
	var offsetGroupId = opt.StartGroupID - 1

	var dbExeId = e.Data.ID
	var max_group_id int

	// 计算GroupState
	groupStates, err := smb.GetGroupStates(opt.StartGroupID)
	if err != nil {
		return
	}
	type tmpGroup struct {
		po.StepGroup
		*states.SubGroupState
	}
	var groups []tmpGroup

	var stateDefinition string

	// 写入所有的步骤
	for _, groupState := range groupStates {
		state := groupState.State
		statetype := state.GetType()

		stateDefinition, err = state.GetDefinition()
		if err != nil {
			return resp, err
		}
		task := po.Step{
			ExecutionID:  dbExeId,
			ExecuteIndex: deIndex,
			GroupID:      groupState.GroupID + offsetGroupId,
			Name:         groupState.Name,
			GroupIndex:   groupState.GroupIndex,
			Type:         statetype,
			ExecuteCount: 0, //初始化执行次数为0
			Definition:   stateDefinition,
			Depth:        groupState.Depth,
			Status:       string(StepStatus.Created),
			Data:         "{}",
		}
		err = tx.Create(&task).Error
		if err != nil {
			return
		}

		// 跟踪max_gruop_id
		if groupState.GroupID+opt.StartGroupID > max_group_id {
			max_group_id = groupState.GroupID + opt.StartGroupID
		}

		// 新增stepgroup
		if groupState.SubGroup != nil {
			neSg := tmpGroup{
				po.StepGroup{
					StepID:      task.ID,
					ExecutionID: dbExeId,
					SubGroupID:  groupState.SubGroup.SubGroupID,
					StartAt:     groupState.SubGroup.SubStartAt,
					GroupIndex:  groupState.GroupIndex,
				},
				groupState.SubGroup,
			}
			groups = append(groups, neSg)
		}

		// for idx, g := range groupState.SubGroup {
		// 	newg := tmpgroup{
		// 		po.StepGroup{
		// 			ExecutionID: dbexecution.ID,
		// 			SubGroupID:  g.GroupID,
		// 			StepID:      task.ID,
		// 			GroupIndex:  idx,
		// 			// Status:           string(TaskStatus.Created),
		// 		},
		// 		g,
		// 	}
		// 	groups = append(groups, newg)

		// }
		deIndex--
	}

	// 从写入的step中查找相关信息, 写入stepgroup
	var masterstep po.Step
	for _, g := range groups {
		masterstep = po.Step{
			ExecutionID: dbExeId,
			Name:        g.SubGroupState.MasterName,
			GroupID:     g.SubGroupState.MasterGroupID + offsetGroupId,
		}
		err = tx.Where(masterstep).Select(StepFields.L1).Take(&masterstep).Error
		if err != nil {
			return
		}
		newsg := g.StepGroup
		newsg.MasterStepID = masterstep.ID
		//
		err = tx.Create(&newsg).Error
		if err != nil {
			return
		}
	}

	// 获得第一个步骤
	var exestartstep = po.Step{
		Name:        smb.StartAt,
		ExecutionID: dbExeId,
		GroupID:     states.StartGroupID + offsetGroupId,
	}

	err = tx.Where(exestartstep).Select(StepFields.L1).Take(&exestartstep).Error
	if err != nil {
		return
	}

	resp.StartStepID = exestartstep.ID
	resp.MinDeIndex = deIndex
	resp.StartGroupID = exestartstep.GroupID
	resp.MaxGroupID = max_group_id

	return
}

// ProcessTimeout 执行超时
func (e *Execution) ProcessTimeout() error {

	var err error
	starttime := time.Now()

	err = e.ChangeExecutionStatus(ExecutionStatus.Failed, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	event1 := vo.ExecutionEvent{
		ExecutionID: e.Data.ID,
		StartTime:   starttime,
		FinishTime:  now,
		Data: EventContent_ExecutionFailed{
			Error: InnerStateError.StatesTimeout,
		},
	}
	e.ExecutionService.SendExecutionEvents(event1)
	return nil
}

// ProcessAbortTimeout 超时关闭
func (e *Execution) ProcessAbortTimeout() error {
	var err error
	starttime := time.Now()

	err = e.ChangeExecutionStatus(ExecutionStatus.Aborted, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	event1 := vo.ExecutionEvent{
		ExecutionID: e.Data.ID,
		StartTime:   starttime,
		FinishTime:  now,
		Data: EventContent_ExecutionAbort{
			Error: InnerStateError.StatesAbortTimeout,
		},
	}
	e.ExecutionService.SendExecutionEvents(event1)
	// 清理过期的消息
	err = e.ExecutionService.InnerQueue.CleanExecutionMessage(e.Data.ID)
	return err
}

// ProcessExecutionFailed  强制 Execution失败
func (e *Execution) ProcessExecutionFailed(message queue.InnerMessageBody) error {

	var err error
	starttime := time.Now()

	err = e.ChangeExecutionStatus(ExecutionStatus.Failed, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	event1 := vo.ExecutionEvent{
		ExecutionID: e.Data.ID,
		StartTime:   starttime,
		FinishTime:  now,
		Data: EventContent_ExecutionFailed{
			Cause: string(message.Data),
		},
	}
	e.ExecutionService.SendExecutionEvents(event1)
	return nil
}

// ProcessExecutionSuspend  Execution Suspend
func (e *Execution) ProcessExecutionSuspend(message queue.InnerMessageBody) error {

	var err error
	starttime := time.Now()

	err = e.ChangeExecutionStatus(ExecutionStatus.Suspending, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	event1 := vo.ExecutionEvent{
		ExecutionID: e.Data.ID,
		StartTime:   starttime,
		FinishTime:  now,
		Data:        EventContent_ExecutionSuspend{},
	}
	e.ExecutionService.SendExecutionEvents(event1)
	return nil
}

// ProcessExecutionBlocked  Execution Suspend
func (e *Execution) ProcessExecutionBlocked(message queue.InnerMessageBody) error {

	var err error
	starttime := time.Now()

	err = e.ChangeExecutionStatus(ExecutionStatus.Blocked, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	event1 := vo.ExecutionEvent{
		ExecutionID: e.Data.ID,
		StartTime:   starttime,
		FinishTime:  now,
		Data:        EventContent_ExecutionBlocked{},
	}
	e.ExecutionService.SendExecutionEvents(event1)
	return nil
}

// ChangeExecutionStatus 修改Execution状态
func (e *Execution) ChangeExecutionStatus(status _ExecutionStatusType, session rdb.Tx) error {

	var err error
	// 开始事务
	tx, maker := e.ExecutionService.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)

	updateexecution := po.Execution{
		Status: string(status),
	}

	err = tx.Where(po.Execution{ID: e.Data.ID}).Updates(&updateexecution).Error
	if err != nil {
		return err
	}
	return err
}

// ProcessSucceed  处理 Execution Successs
func (e *Execution) ProcessSucceed(output string) error {

	var err error
	starttime := time.Now()
	var dbExe *po.Execution
	// 开始事务
	tx, maker := e.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	// 查询 output
	dbExe, err = e.ExecutionService.QueryExecutionByID(e.Data.ID, []string{"id", "output"}, tx)
	if err != nil {
		return err
	}
	// 更新状态为Success, 同时更新Finish Time
	now := time.Now()
	updateexecution := po.Execution{
		FinishTime: &now,
		Output:     output,
		Status:     string(ExecutionStatus.Success),
	}

	err = tx.Where(po.Execution{ID: e.Data.ID}).Updates(&updateexecution).Error
	if err != nil {
		return err
	}

	tx.Commit()

	finishtime := time.Now()
	event1 := vo.ExecutionEvent{
		ExecutionID: dbExe.ID,
		StepName:    ExecutionEventStateName.End,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_ExecutionSucceeded{
			Output: dbExe.Output,
		},
	}
	e.ExecutionService.SendExecutionEvents(event1)
	// 清理过期的消息
	err = e.ExecutionService.InnerQueue.CleanExecutionMessage(dbExe.ID)
	return err
}

// StopExecution stop certain execution
// 终止Execution ，最高优先级，不论什么状态， 都修改成Aborted 状态。
func (e *Execution) StopExecution(errorcode string, cause string) error {

	// 需要注意， 和其他操作同时并发的时候， 需要加锁，防止并发操作下的Stop失败。
	var err error
	// var has bool

	// 加锁
	lock := e.ExecutionService.LockService.LockExecution(e.Data.ID)
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()
	// 开始事务
	tx, maker := e.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	resp, err := e._StopExecution(errorcode, cause, tx)
	if err != nil {
		return err
	}
	tx.Commit()

	e.ExecutionService.SendExecutionEvents(resp.Events...)
	// 清理过期的消息
	err = e.ExecutionService.InnerQueue.CleanExecutionMessage(e.Data.ID)
	if err != nil {
		return err
	}
	return nil
}

// _StopExecution 内部使用的state
func (e *Execution) _StopExecution(errorcode string, cause string, tx rdb.Tx) (
	resp StopExecutionResponse, err error) {

	var dbExeId = e.Data.ID
	var dbExecutionUuid = e.Data.UUID
	starttime := time.Now()
	//避免重复stop
	if e.Data.Status == string(ExecutionStatus.Aborted) {
		err = fmt.Errorf("%w : execution [ %s ] has been stopped", ErrorExecutionStatus, dbExecutionUuid)
		return
	}
	var exceptData = ExceptionData{
		Error: errorcode,
		Cause: cause,
	}

	exceptDatabyte, _ := json.Marshal(exceptData)
	exceptDataString := string(exceptDatabyte)
	updateexecution := po.Execution{
		Status:    string(ExecutionStatus.Aborted),
		Exception: exceptDataString,
	}

	err = tx.Where(po.Execution{ID: dbExeId}).Updates(&updateexecution).Error
	if err != nil {
		return
	}

	finishtime := time.Now()
	event := vo.ExecutionEvent{
		ExecutionID: dbExeId,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_ExecutionAbort{
			Error: errorcode,
			Cause: cause,
		},
	}
	resp.Events = append(resp.Events, event)

	return

}
