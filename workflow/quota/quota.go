// package quota  defined cloudflow quota here
package quota

// hard limit

// MaxWorkflowSize 最大的工作流模板大小 512KB
var MaxWorkflowSize = 512 * 1024 // 512KB

// MaxStartExecutionInputSize 最大的工作流执行输入大小 128KB
var MaxStartExecutionInputSize = 128 * 1024

// MaxInputSize 最大的工作流执行输入大小 128KB
var MaxInputSize = 128 * 1024

// MaxOutputSize 最大的工作流执行输出大小 128KB
var MaxOutputSize = 128 * 1024

// MaxTaskInputSize 最大的Task Input大小 32KB
var MaxTaskInputSize = 32 * 1024

// MaxTaskOutputSize 最大的Task Submit Output大小 32KB
var MaxTaskOutputSize = 32 * 1024

// MaxTaskDataSize 最大的StoreTaskData大小 32KB
var MaxTaskDataSize = 32 * 1024

// MaxParallelBranchNumber 最大的并行分支数
var MaxParallelBranchNumber = 500

// MaxMapBranchNumber 最大的Map分支数
var MaxMapBranchNumber = 500

// MaxMapConncurrencyLimit 最大的Map并发数
var MaxMapConncurrencyLimit = 100

// MaxExecuteTimes 最大的执行次数
var MaxExecuteTimes = 10000

// MaxStepNameSizeLimit 最大的步骤名称长度
var MaxStepNameSizeLimit = 200

// MaxRunningExecutionLimit 最大的运行中的工作流执行数
var MaxRunningExecutionLimit = 1000

// MaxTitleSizeLimit 最大的标题长度
var MaxTitleSizeLimit = 200

// MaxWorkflowURISizeLimit 最大的工作流URI长度
var MaxWorkflowURISizeLimit = 200

// MaxActivityURISizeLimit 最大的活动URI长度
var MaxActivityURISizeLimit = 200

// MaxTaskResourceLengthLimit 最大的Task Resource 定义长度
var MaxTaskResourceLengthLimit = 200

// MaxExecutionUUIDSize 最大的执行UUID长度
var MaxExecutionUUIDSize = 200
