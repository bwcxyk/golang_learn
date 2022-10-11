/*
   Copyright (c) [2021] IT.SOS
   golibs is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package minio

import (
	"golibs/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"sync"
)

var once sync.Once
var minioNew *minio.Client

func NewMinio() *minio.Client {
	once.Do(func() {
		var err error
		endpoint := config.GetMinio().GetEndpoint()
		accessKeyID := config.GetMinio().GetAccessKeyID()
		secretAccessKey := config.GetMinio().GetSecretAccessKey()
		useSSL := config.GetMinio().GetUseSSL()

		// Initialize minio client object.
		minioNew, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			panic(err)
		}
	})
	return minioNew
}
