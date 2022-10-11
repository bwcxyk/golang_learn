/*
   Copyright (c) [2021] IT.SOS
   golibs is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package reflects

import "reflect"

// tagType 的值
const (
	YAML = "yaml"
	JSON = "json"
)

func TagToValueFlip(t reflect.Value, tagType string) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Type().Field(i)
		if f.Type.Name() != "string" {
			continue
		}
		tag := f.Tag.Get(tagType)
		s := t.FieldByName(f.Name)
		if !s.CanSet() {
			panic(f.Name + " => not set value. " + tag)
		}
		s.SetString(tag)
	}
}
