package service

import (
	"pick/conf"
)

//根据链接转发源
func WhereGo(domain string)  {
	rid := conf.GetKeys(domain)
	switch rid {
		case 1:
			BookLists(domain)
			break
		case 2:
			BookTwoLists(domain)
			break
	}
}