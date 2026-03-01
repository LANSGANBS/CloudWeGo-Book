// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"context"
	"strings"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/conf"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var userDB *gorm.DB

func InitUserDB() {
	dsn := conf.GetConf().MySQL.DSN
	dsn = strings.Replace(dsn, "/gorm?", "/user?", 1)
	var err error
	userDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		hlog.Fatalf("failed to connect user database: %v", err)
	}
}

func CheckAdmin(ctx context.Context, userId uint32) (bool, error) {
	if userDB == nil {
		InitUserDB()
	}

	var isAdmin bool
	err := userDB.Table("user").Select("is_admin").Where("id = ?", userId).Scan(&isAdmin).Error
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to check admin status for user %d: %v", userId, err)
		return false, err
	}

	return isAdmin, nil
}

func GetUserCount(ctx context.Context) (int64, error) {
	if userDB == nil {
		InitUserDB()
	}

	var count int64
	err := userDB.Table("user").Where("is_admin = ?", false).Count(&count).Error
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get user count: %v", err)
		return 0, err
	}

	return count, nil
}
