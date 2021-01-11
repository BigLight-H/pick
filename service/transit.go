package service

import (
	"github.com/davecgh/go-spew/spew"
	"pick/conf"
)

//根据链接转发源
func WhereGo(domain string)  {
	rid := conf.GetKeys(domain)
	spew.Dump(rid)
	switch rid {
		case 1:
			BookLists(domain)
			break
		case 2:
			BookTwoLists(domain)
			break
	}
}