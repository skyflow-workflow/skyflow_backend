package executor

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"trpc.group/trpc-go/tnet/log"
)

// Execution  StateMachine 执行实例
type Execution struct {
	StateMachine *grammar.StateMachine
	Data         *po.Execution
	// state map
	States           map[string]Step
	ExecutionService ExecutionService
}

// ExecutionInfo  内置的ExecutionInfor Struct
type ExecutionInfo struct {
	UUID    string `json:"uuid"`    // execution uuid
	URLPath string `json:"urlpath"` // 访问路径
	URL     string `json:"url"`     // 完整访问地址

}

// NewExecutionFromID Create Execution by execution id
func NewExecutionFromID(id int, svc ExecutionService) (*Execution, error) {

	var err error
	var dbexe po.Execution

	// 使用最小的数据集来来初始化
	dbexe, err = svc.QueryExecutionByID(id, ExecutionFields.L1, nil)
	if err != nil {
		return nil, err
	}

	exe, err := NewExecutionFromData(&dbexe, svc)
	return exe, err
}

// NewExecutionFromData  根据DB中的数据创建一个ServiceExecution 类型解析
/* data : po.Execution 类型， Execution 信息
* datastates : Execution 下所有的DB节点
* StepGroups : Execution 下所有节点的关联关系.
 */
func NewExecutionFromData(data *po.Execution, svc ExecutionService) (exe *Execution, err error) {

	exe = &Execution{
		Data:   data,
		States: map[string]Step{},
		Bone: ExecutionBone{
			StartAt: "",
			States:  map[string]StepBone{},
		},
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
func (exe *Execution) FullInit() error {
	// 如果已经初始化， 则忽略
	var err error
	if exe.StateMachine != nil {
		return nil
	}

	var dbexecution po.Execution
	id := exe.Data.ID

	// 查询到全量的数据
	dbexecution, err = exe.ExecutionService.QueryExecutionByID(id, []string{}, nil)

	if err != nil {
		return err
	}
	// 后面需要计算input
	exe.Data = &dbexecution

	wf, err := parser.ParseWorkflow(dbexecution.FlowDefinition)
	if err != nil {
		return err
	}
	sm := wf.GetNode()

	exe.StateMachine = sm

	return nil
}

// ProcessEvent 处理事件
func (exe *Execution) ProcessEvent(msg queue.InnerMessage) error {

	var err error

	lock := exe.ExecutionService.lockservice.LockExecution(msg.ExecutionID)
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()

	var dbexe po.Execution

	dbexe, err = exe.ExecutionService.QueryExecutionByID(msg.ExecutionID, []string{"id", "status"}, nil)
	if err != nil {
		return err
	}

	checkstatus, ok := ExecutionEventCheckStatus[msg.Type]
	// 如果状态在规则中， 则检查状态，如果不在， 则忽略不检查
	if ok {
		// 如果不在预期状态中， 忽略事件
		if !toolkit.StringInSlice(checkstatus, dbexe.Status) {
			msg := fmt.Sprintf("execution '%d' process event '%s' current status '%s' not match ",
				msg.ExecutionID, msg.Type, dbexe.Status)
			log.Error(msg)
			return nil
		}
	}

	switch msg.Type {
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
		err = fmt.Errorf("unrecognize event '%s' ", msg.Type)
	}

	if err != nil {
		err = exe.ExecutionService.ExecutionErrorProcess(err, exe.Data.ID)
	}
	return err
}

// GetInput GetInput
func (e *Execution) GetInput() (interface{}, error) {
	var inputdata, newinputdata interface{}
	var err error

	err = json.Unmarshal([]byte(e.Data.Input), &inputdata)
	if err != nil {
		return nil, err
	}

	fullurl := e.ExecutionService.domainService.QueryExecutionURL(e.Data.UUID)

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
	var dbsteps []po.Step
	var dbstepgroups []po.StepGroup
	//
	var dbstepmap = map[string]*po.Step{}

	// 查询数据
	tx := e.ExecutionService.metadb.DB()
	err = tx.Where(po.Step{ExecutionID: e.Data.ID}).Find(&dbsteps).Error
	if err != nil {
		return bone, err
	}

	err = tx.Where(po.StepGroup{ExecutionID: e.Data.ID}).Find(&dbstepgroups).Error
	if err != nil {
		return bone, err
	}

	// levelconnector  是子流程连接的节点，主要是Parallel/Map 节点
	var levelconnector = []*po.Step{}
	// 构建bone
	StartGroupID := grammar.StartGroupID
	// GroupBoneMap  group 组内的bone map，
	// 两层map结构， 第一层是 group_id， 第二层是statename ,都可以从DB钟获得
	// 通过初始化， GroupBoneMap 按照分层存储所有的节点的Bone。 节点都是 StateBone 类型
	var GroupBoneMap = map[int]map[string]StepBone{
		// StartGroupID: map[string]StateBone{},
	}
	// StateBoneMap  按照step_id 生成的 Map
	var StateMap = map[int]Step{}

	for di := range dbsteps {
		dp := &dbsteps[di]
		// 存成map
		dbstepmap[dp.Name] = dp
		var groupid int
		newstate, err := NewStepFromData(dp, e.ExecutionService, e.Data.FlowType)
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
		if dp.Type == grammar.StateType.Map || dp.Type == grammar.StateType.Parallel {
			levelconnector = append(levelconnector, dp)
		}
	}

	// 处理StepGroup, 把StepGroup 按照step_id 分成不同的分组,  方便计算出一个State下有多少个子流程
	// sgmap  StepGroup 按照start_id 进行分组
	var sgmap = map[int][]*po.StepGroup{}
	var step_id int
	for index := range dbstepgroups {
		step_id = dbstepgroups[index].StepID
		if sgg, ok := sgmap[step_id]; ok {
			sgmap[step_id] = append(sgg, &(dbstepgroups[index]))
		} else {
			sgmap[step_id] = []*po.StepGroup{&(dbstepgroups[index])}
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
			startatnode, ok := dbstepmap[subgroup.StartAt]
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
	e.Bone = toplevelStateBone
	return e.Bone, nil
}

// ProcessInit  初始化execution
// 初始化，是将 Execution 解析成多个Step，并且持久化到 po.Step 中。
// NOCC:golint/fnsize("设计如此")
func (e *Execution) ProcessInit() error {
	var err error
	var dbexecution *po.Execution

	err = e.FullInit()
	if err != nil {
		return err
	}

	var sm = e.StateMachine
	dbexecution = e.Data

	starttime := time.Now()

	// 计算第一个Execution 经过Header的input
	stateinput, err := e.GetInput()
	if err != nil {
		return err
	}
	stateinputstr, err := grammar.ToString(stateinput)
	if err != nil {
		return err
	}
	// 获得Header
	header := e.StateMachine.StateMachineHeader

	// 先开启session， 所有的事情都应该在一个session中
	// 开始事务
	tx, maker := e.ExecutionService.metadb.NewSessionMaker(nil)
	defer maker.Close(&err)

	now := time.Now()
	updateExecution := po.Execution{
		Status:    string(ExecutionStatus.Running),
		StartTime: &now,
		Input:     stateinputstr,
		Header:    header.GetDefinition(),
	}

	err = tx.Where(po.Execution{ID: dbexecution.ID}).Updates(&updateExecution).Error
	if err != nil {
		return err
	}

	insmresp, err := e.InsertStateMachine(&sm.StateMachineBody, InsertStateMachineOption{
		StartDeindex: -1,
		StartDepth:   grammar.StartDepth,
		StartGroupID: grammar.StartGroupID,
	}, tx)
	if err != nil {
		return err
	}
	updatetask := po.Step{
		Input:  stateinputstr,
		Status: string(StepStatus.WaitInit),
	}

	err = tx.Where(po.Step{ID: insmresp.StartStepID}).Updates(&updatetask).Error
	if err != nil {
		return err
	}

	tx.Commit()

	TaskCreateTime := time.Now()
	// Create Event
	event1 := ExecutionEvent{
		ExecutionID: dbexecution.ID,
		StepName:    ExecutionEventStateName.Start,
		StartTime:   starttime,
		FinishTime:  TaskCreateTime,
		Data: EventContent_ExecutionStart{
			Input: dbexecution.Input,
			URI:   dbexecution.URI,
		},
	}
	e.ExecutionService.SendExecutionEvents(event1)

	// 发送消息
	// 处理下一个节点
	msg := StateExecuteMessage{
		Unblock: false,
	}
	message := NewStateMessage(dbexecution.ID, MessageType.StateNewTurn, insmresp.StartStepID, msg)
	err = e.ExecutionService.innerqueue.SendInnerMessage(message, time.Now())
	if err != nil {
		return err
	}
	timeout := e.StateMachine.GetTimeout()
	if timeout.Timeout > 0 {
		// 发送超时事件
		message = NewExecutionMessage(dbexecution.ID, MessageType.ExecutionTimout, nil)
		err = e.ExecutionService.innerqueue.SendInnerMessage(message, time.Now().Add(timeout.Timeout))
		if err != nil {
			return err
		}
	}
	if timeout.AbortTimeout > 0 {
		// 发送Abort事件
		message = NewExecutionMessage(dbexecution.ID, MessageType.ExecutionAbortTimout, nil)
		err = e.ExecutionService.innerqueue.SendInnerMessage(message, time.Now().Add(timeout.AbortTimeout))
		if err != nil {
			return err
		}
	}
	return nil
}

// InsertStateMachine insert statemachine to db
func (e *Execution) InsertStateMachine(smb *grammar.StateMachineBody, opt InsertStateMachineOption,
	tx rdb.Session) (resp InsertStateMachineResponse, err error) {

	var deindex = opt.StartDeindex
	// offset_groupid GroupID偏移量
	var offset_groupid = opt.StartGroupID - 1

	var dbexecution_id = e.Data.ID
	var max_group_id int

	// 计算GroupState
	groupStates, err := smb.GetGroupStates()
	if err != nil {
		return
	}
	type tmpgroup struct {
		po.StepGroup
		*grammar.SubGroupState
	}
	var groups []tmpgroup

	// 写入所有的步骤
	for _, groupState := range groupStates {
		state := groupState.State
		statetype := state.GetType()
		task := po.Step{
			ExecutionID:  dbexecution_id,
			ExecuteIndex: deindex,
			GroupID:      groupState.GroupID + offset_groupid,
			Name:         groupState.Name,
			GroupIndex:   groupState.GroupIndex,
			Type:         statetype,
			ExecuteCount: 0, //初始化执行次数为0
			Definition:   state.GetDefinition(),
			Depth:        groupState.Depth,
			Status:       string(StepStatus.Created),
			Data:         "{}",
			References:   "{}",
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
			newg := tmpgroup{
				po.StepGroup{
					StepID:      task.ID,
					ExecutionID: dbexecution_id,
					SubGroupID:  groupState.SubGroup.SubGroupID,
					StartAt:     groupState.SubGroup.SubStartAt,
					GroupIndex:  groupState.GroupIndex,
				},
				groupState.SubGroup,
			}
			groups = append(groups, newg)
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
		deindex--
	}

	// 从写入的step中查找相关信息, 写入stepgroup
	var masterstep po.Step
	for _, g := range groups {
		masterstep = po.Step{
			ExecutionID: dbexecution_id,
			Name:        g.SubGroupState.MasterName,
			GroupID:     g.SubGroupState.MasterGroupID + offset_groupid,
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
		ExecutionID: dbexecution_id,
		GroupID:     grammar.StartGroupID + offset_groupid,
	}

	err = tx.Where(exestartstep).Select(StepFields.L1).Take(&exestartstep).Error
	if err != nil {
		return
	}

	resp.StartStepID = exestartstep.ID
	resp.MinDeindex = deindex
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
	event1 := ExecutionEvent{
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
	event1 := ExecutionEvent{
		ExecutionID: e.Data.ID,
		StartTime:   starttime,
		FinishTime:  now,
		Data: EventContent_ExecutionAbort{
			Error: InnerStateError.StatesAbortTimeout,
		},
	}
	e.ExecutionService.SendExecutionEvents(event1)
	// 清理过期的消息
	err = e.ExecutionService.innerqueue.CleanExecutionMessage(e.Data.ID)
	return err
}

// ProcessExecutionFailed  强制 Execution失败
func (e *Execution) ProcessExecutionFailed(message queue.InnerMessage) error {

	var err error
	starttime := time.Now()

	err = e.ChangeExecutionStatus(ExecutionStatus.Failed, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	event1 := ExecutionEvent{
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
func (e *Execution) ProcessExecutionSuspend(message queue.InnerMessage) error {

	var err error
	starttime := time.Now()

	err = e.ChangeExecutionStatus(ExecutionStatus.Suspending, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	event1 := ExecutionEvent{
		ExecutionID: e.Data.ID,
		StartTime:   starttime,
		FinishTime:  now,
		Data:        EventContent_ExecutionSuspend{},
	}
	e.ExecutionService.SendExecutionEvents(event1)
	return nil
}

// ProcessExecutionBlocked  Execution Suspend
func (e *Execution) ProcessExecutionBlocked(message queue.InnerMessage) error {

	var err error
	starttime := time.Now()

	err = e.ChangeExecutionStatus(ExecutionStatus.Blocked, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	event1 := ExecutionEvent{
		ExecutionID: e.Data.ID,
		StartTime:   starttime,
		FinishTime:  now,
		Data:        EventContent_ExecutionBlocked{},
	}
	e.ExecutionService.SendExecutionEvents(event1)
	return nil
}

// ChangeExecutionStatus 修改Execution状态
func (e *Execution) ChangeExecutionStatus(status _ExecutionStatusType, session rdb.Session) error {

	var err error
	// 开始事务
	tx, maker := e.ExecutionService.metadb.NewSessionMaker(session)
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
	var dbexecution po.Execution
	// 开始事务
	tx, maker := e.ExecutionService.metadb.NewSessionMaker(nil)
	defer maker.Close(&err)

	// 查询 output
	dbexecution, err = e.ExecutionService.QueryExecutionByID(e.Data.ID, []string{"id", "output"}, tx)
	if err != nil {
		return err
	}
	// 更新状态为Success, 同时更新Finsishtime
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
	event1 := ExecutionEvent{
		ExecutionID: dbexecution.ID,
		StepName:    ExecutionEventStateName.End,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_ExecutionSucceeded{
			Output: dbexecution.Output,
		},
	}
	e.ExecutionService.SendExecutionEvents(event1)
	// 清理过期的消息
	err = e.ExecutionService.innerqueue.CleanExecutionMessage(dbexecution.ID)
	return err
}

// StopExecution stop certain execution
// 终止Execution ，最高优先级，不论什么状态， 都修改成Aborted 状态。
func (e *Execution) StopExecution(errorcode string, cause string) error {

	// 需要注意， 和其他操作同时并发的时候， 需要加锁，防止并发操作下的Stop失败。
	var err error
	// var has bool

	// 加锁
	lock := e.ExecutionService.lockservice.LockExecution(e.Data.ID)
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()
	// 开始事务
	tx, maker := e.ExecutionService.metadb.NewSessionMaker(nil)
	defer maker.Close(&err)

	resp, err := e._StopExecution(errorcode, cause, tx)
	if err != nil {
		return err
	}
	tx.Commit()

	e.ExecutionService.SendExecutionEvents(resp.Events...)
	// 清理过期的消息
	err = e.ExecutionService.innerqueue.CleanExecutionMessage(e.Data.ID)
	if err != nil {
		return err
	}
	return nil
}

// _StopExecution 内部使用的state
func (e *Execution) _StopExecution(errorcode string, cause string, tx rdb.Session) (
	resp StopExecutionResponse, err error) {

	var dbexecution_id = e.Data.ID
	var dbexecution_uuid = e.Data.UUID
	starttime := time.Now()
	//避免重复stop
	if e.Data.Status == string(ExecutionStatus.Aborted) {
		err = fmt.Errorf("%w : execution [ %s ] has been stopped", ErrorExecutionStatus, dbexecution_uuid)
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

	err = tx.Where(po.Execution{ID: dbexecution_id}).Updates(&updateexecution).Error
	if err != nil {
		return
	}

	finishtime := time.Now()
	event := ExecutionEvent{
		ExecutionID: dbexecution_id,
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
