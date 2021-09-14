/*
@Author : YaoKun
@Time : 2021/9/10 9:58
*/

package IdCard

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	ocr "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ocr/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ocr/v1/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ocr/v1/region"
	"github.com/spf13/viper"

	_ "config"
)

func InitConfig() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func IdCard() {
	//ak := "xxxxxxxxxxx"
	//sk := "xxxxxxxxxxxxx"
    ak := viper.GetString(`HuaWeiOcr.AccessKey`)
    sk := viper.GetString(`HuaWeiOcr.SecretAccessKey`)

	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := ocr.NewOcrClient(
		ocr.OcrClientBuilder().
			WithRegion(region.ValueOf("cn-north-4")).
			WithCredential(auth).
			Build())

	request := &model.RecognizeIdCardRequest{}
	sideIdCardRequestBody:= "front"
	urlIdCardRequestBody:= "https://file.yuanfusc.com/group1/M00/09/B5/rBOz8V25Iu2ACDmCAACQ5EvSOMg247.jpg"
	request.Body = &model.IdCardRequestBody{
		Side: &sideIdCardRequestBody,
		Url: &urlIdCardRequestBody,
	}
	response, err := client.RecognizeIdCard(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
}