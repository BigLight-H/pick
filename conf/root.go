package conf

import "github.com/beego/beego/core/config"

//源信息
var RootLinks map[int]string

func Init() {
	sone, _ := config.String("source_root")
	stwo, _ := config.String("source_root")
	RootLinks[1] = sone
	RootLinks[2] = stwo
}

func GetKeys(domain string) int {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率较高
	j := 1
	keys := 0
	for k := range RootLinks {
		if domain == RootLinks[j] {
			keys = k
		}
		j++
	}
	return keys
}
