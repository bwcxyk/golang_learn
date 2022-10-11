/*
   Copyright (c) [2021] IT.SOS
   kn is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package utils

import (
	"testing"
	"time"
)

func getTime(s string) time.Time {
	e := time.Now()
	d, _ := time.ParseDuration(s)
	return e.Add(d)
}

func Test_GetDuration(t *testing.T) {
	timeRange := []time.Time{
		getTime("-10s"),
		getTime("-2m"),
		getTime("-2h"),
		getTime("-24h"),
		getTime("-48h"),
		getTime("-72h"),
		getTime("-168h"),
		getTime("-336h"),
		getTime("-672h"),
		getTime("-720h"),
		getTime("-8640h"),
		getTime("-8760h"),
		getTime("-26280h"),
	}
	for _, v := range timeRange {
		t.Log(TimeDuration(v))
	}
}
