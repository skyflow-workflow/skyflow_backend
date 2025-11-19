package domain

import (
	"fmt"
)

// domainService 生成域名相关服务
type DomainService = *domainService

type domainService struct {
	Domain string
}

var DefaultDomainService = NewDomainService("http://127.0.0.1")

func NewDomainService(domain string) DomainService {

	return &domainService{
		Domain: domain,
	}

}
func (d *domainService) QueryExecutionURL(uuid string) string {
	var urlpath = fmt.Sprintf("%s/#/task_detail?UUID=%s", d.Domain, uuid)
	return urlpath
}
