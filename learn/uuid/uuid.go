/*
@Author : YaoKun
@Time : 2023/5/4 15:37
*/

package main

import (
	"fmt"
	"github.com/google/uuid"
)

// v1,v4都是每次生成一个唯一的ID.
// v1同一时刻的输出非常相似，v1末尾nodeID部分用的都是mac地址，前面time的mid,high以及clock序列都是一样的，只有time-low部分不同。
func testv1() {
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("testv1: %v,%v\n", id, id.Version().String())
}

// v4加入了随机数，对各个部分都进行了随机处理，同一时刻的输出差别很大。
func testv4() {
	id := uuid.New()
	fmt.Printf("testv4: %v,%v\n", id, id.Version().String())
}

// v2 NewDCEGroup()根据os.Getgid取到的用户组ID来生成uuid,同一时刻的输出是相同的。
func testv2G() {
	id, err := uuid.NewDCEGroup()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("testv2G: %v,%v\n", id, id.Version().String())
}

// v2 NewDCEPerson()根据os.Getuid取到的用户ID来生成uuid,同一时刻的输出也是相同的。
func testv2P() {
	id, err := uuid.NewDCEPerson()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("testv2P: %v,%v\n", id, id.Version().String())
}

// v3 NewMD5(space UUID, data []byte)是根据参数传入的UUID结构体和[]byte再重新转换一次。只要传入参数相同则任意时刻的输出也相同。
func testv3() {
	id0, err := uuid.NewDCEPerson()
	if err != nil {
		fmt.Println(err)
		return
	}
	id := uuid.NewMD5(id0, []byte("fssds32"))
	fmt.Printf("testv3: %v,%v\n", id, id.Version().String())
}

// v5 NewSHA1(space UUID, data []byte)是根据参数传入的UUID结构体和[]byte再重新转换一次。只要传入参数相同则任意时刻的输出也相同。
func testv5() {
	id0, err := uuid.NewDCEPerson()
	if err != nil {
		fmt.Println(err)
		return
	}
	id := uuid.NewSHA1(id0, []byte("fssds32"))
	fmt.Printf("testv5: %v,%v\n", id, id.Version().String())
}

func main() {
	for i := 0; i < 2; i++ {
		testv1()
	}

	fmt.Println("--------------------")

	for i := 0; i < 2; i++ {
		testv4()
	}

	fmt.Println("--------------------")

	for i := 0; i < 2; i++ {
		testv2G()
	}

	fmt.Println("--------------------")

	for i := 0; i < 2; i++ {
		testv2P()
	}

	fmt.Println("--------------------")

	for i := 0; i < 2; i++ {
		testv3()
	}

	fmt.Println("--------------------")

	for i := 0; i < 2; i++ {
		testv5()
	}

	fmt.Println("--------------------")
}
