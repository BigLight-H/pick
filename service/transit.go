package service

import (
	"github.com/davecgh/go-spew/spew"
	"pick/util"
)

//根据链接转发源
func WhereGo(domain string)  {
	rid := util.GetKeys(domain)
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