/*
@Author : YaoKun
@Time : 2022/6/16 9:35
*/

package main

import (
	"fmt"
	"github.com/issue9/cnregion"
)

func main() {
	v, _ := cnregion.LoadFile("regions.db", "-", 2021)
	p := v.Provinces()            // 返回所有省列表
	cities := p[0].Items()        // 返回该省下的所有市
	counties := cities[0].Items() // 返回该市下的所有县
	towns := counties[0].Items()  // 返回所有镇
	villages := towns[0].Items()  // 所有村和街道信息

	//d := v.Districts() // 按以前的行政大区进行划分
	//provinces := d[0].Items() // 该大区下的所有省份

	//list := v.Search("温州", nil) // 按索地名中带温州的区域列表
}
