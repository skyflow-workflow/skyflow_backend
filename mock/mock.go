package mock

import (
	"database/sql"

	"github.com/mmtbak/microlibrary/config"
	"github.com/mmtbak/microlibrary/rdb"
)

var (
	MockDB       *sql.DB
	MockDBClient *rdb.DBClient
)
var LocalUnitTestMySQLConfig = config.AccessPoint{
	Source: "mysql://root:root@tcp(127.0.0.1:3306)/testdb?charset=utf8&parseTime=True&loc=Local",
	Options: map[string]interface{}{
		"sqllevel":    "info",
		"maxopenconn": 100,
		"maxidleconn": 100,
	},
}

var LocalUnitTestKafka = config.AccessPoint{
	Source: "kafka://localhost:9092/?" +
		"topics=my-event-test-topic" +
		"&numpartition=2&numreplica=1&autocommitsecond=1" +
		"initial=oldest&version=1.1.1",
}

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
	config, err := rdb.ParseConfig(LocalUnitTestMySQLConfig)
	if err != nil {
		return err
	}
	MockDBClient, err = rdb.NewDBClient(config)
	if err != nil {
		return err
	}
	return nil
}

func GetKakfkaConfig() config.AccessPoint {
	return LocalUnitTestKafka
}
