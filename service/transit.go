package service

import (
	"github.com/davecgh/go-spew/spew"
	"pick/util"
)

//根据链接转发源
func WhereGo(domain string, num int)  {
	rid := util.GetKeys(domain, num)
	spew.Dump(rid, domain)
	switch rid {
		case 1:
			BookLists(domain, rid)
			break
		case 2:
			BookTwoLists(domain, rid)
			break
		case 3:
			BookThreeLists(domain, rid)
			break
	}
}