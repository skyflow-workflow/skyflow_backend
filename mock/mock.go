package mock

import (
	"database/sql"

	"github.com/mmtbak/microlibrary/mq"
	"github.com/mmtbak/microlibrary/rdb"
)

var (
	MockDB       *sql.DB
	MockDBClient *rdb.DBClient
	MockKafkaMQ  *mq.KafkaMessageQueue
)
var LocalUnitTestMySQLConfig = rdb.Config{
	DSN:          "mysql://root:rootpassword@tcp(127.0.0.1:3306)/testdb?charset=utf8&parseTime=True&loc=Local",
	LogLevel:     "info",
	MaxOpenConns: 200,
	MaxIdleConns: 200,
}

var LocalUnitTestKafkaDSN = "kafka://localhost:9092/?" +
	"topics=my-event-test-topic" +
	"&numpartition=2&numreplica=1&autocommitsecond=1" +
	"initial=oldest&version=1.1.1"

func init() {
	err := InitMockDB()
	if err != nil {
		panic(err)
	}
}

func GetMockDBClient() *rdb.DBClient {
	return MockDBClient
}

// InitMockDB initialize the test database
func InitMockDB() error {

	var err error
	MockDBClient, err = rdb.NewDBClient(&LocalUnitTestMySQLConfig)
	if err != nil {
		return err
	}
	return nil
}

func GetMockKafkaDSN() string {
	return LocalUnitTestKafkaDSN
}

func InitMockKafka() error {
	var err error
	MockKafkaMQ, err = mq.NewKafkaMessageQueue(GetMockKafkaDSN())
	if err != nil {
		return err
	}
	return nil
}

func GetMockKafkaMQ() *mq.KafkaMessageQueue {
	return MockKafkaMQ
}
