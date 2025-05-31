package exporter

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/mmtbak/dsnparser"
)

func TestParseDSN(t *testing.T) {

	var testcases = []struct {
		dsn       string
		wantError bool
	}{
		{
			dsn:       "kafka://username:pasword@tcp(ip1:9093,ip2:9093,ip3:9093)?topic=vsulblog",
			wantError: false,
		},
		{
			dsn:       "mysql://user:password@tcp(127.0.0.1:3306)/mydb?charset=utf8",
			wantError: false,
		},
		{
			dsn:       "user:password@tcp(127.0.0.1:3306)/mydb?charset=utf8",
			wantError: false,
		},
		{
			dsn:       "cls://sercetid:sercetkey@cls.ap-chongqing.tencentcloudapi.com",
			wantError: false,
		},
	}

	for _, t := range testcases {
		fmt.Println(t.dsn)
		fmt.Println("parse result ..")
		u, err := url.Parse(t.dsn)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(u.Scheme, u.Host, u.User.Username(), u.User.String(), u.Path)
			fmt.Println(u.User.Password())
			fmt.Println(u.Port(), u.RawFragment)

		}
		d := dsnparser.Parse(t.dsn)
		fmt.Println(d)
	}

}
