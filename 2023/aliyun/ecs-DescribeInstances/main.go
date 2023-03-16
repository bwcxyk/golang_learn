// This file is auto-generated, don't edit it. Thanks.
package main

import (
	env "github.com/alibabacloud-go/darabonba-env/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	string_ "github.com/alibabacloud-go/darabonba-string/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v2/client"
	console "github.com/alibabacloud-go/tea-console/client"
	"github.com/alibabacloud-go/tea/tea"
	"os"
)

func _main(args []*string) (_err error) {
	// 1. 初始化配置
	config := &openapi.Config{}
	// 您的AccessKey ID
	config.AccessKeyId = env.GetEnv(tea.String("ACCESS_KEY_ID"))
	// 您的AccessKey Secret
	config.AccessKeySecret = env.GetEnv(tea.String("ACCESS_KEY_SECRET"))
	//设置请求地址
	config.Endpoint = tea.String("ecs.aliyuncs.com")
	// 设置连接超时为5000毫秒
	config.ConnectTimeout = tea.Int(5000)
	// 设置读超时为5000毫秒
	config.ReadTimeout = tea.Int(5000)
	// 2. 初始化客户端
	client, _err := ecs.NewClient(config)
	if _err != nil {
		return _err
	}

	regionIds := string_.Split(args[0], tea.String(","), tea.Int(50))
	for _, regionId := range regionIds {
		describeInstancesRequest := &ecs.DescribeInstancesRequest{
			PageSize: tea.Int32(100),
			RegionId: regionId,
		}
		resp, _err := client.DescribeInstances(describeInstancesRequest)
		if _err != nil {
			return _err
		}

		instances := resp.Body.Instances.Instance
		console.Log(tea.String(tea.StringValue(regionId) + " 下 ECS 实例列表:"))
		for _, instance := range instances {
			console.Log(tea.String("  " + tea.StringValue(instance.HostName) + " 实例ID " + tea.StringValue(instance.InstanceId) + " CPU:" + tea.ToString(tea.Int32Value(instance.Cpu)) + "  内存:" + tea.ToString(tea.Int32Value(instance.Memory)) + " MB 规格：" + tea.StringValue(instance.InstanceType) + " 系统:" + tea.StringValue(instance.OSType) + "(" + tea.StringValue(instance.OSName) + ") 状态：" + tea.StringValue(instance.Status)))
		}
	}
	return _err
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}
