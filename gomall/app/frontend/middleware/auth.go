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

package middleware

import (
	"context"
	"fmt"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sessions"
)

func GlobalAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.Default(c)
		userId := session.Get("user_id")
		if userId == nil {
			c.Next(ctx)
			return
		}
		ctx = context.WithValue(ctx, utils.UserIdKey, userId)
		c.Next(ctx)
	}
}

func Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.Default(c)
		userId := session.Get("user_id")
		if userId == nil {
			byteRef := c.GetHeader("Referer")
			ref := string(byteRef)
			next := "/sign-in"
			if ref != "" {
				if utils.ValidateNext(ref) {
					next = fmt.Sprintf("%s?next=%s", next, ref)
				}
			}
			c.Redirect(302, []byte(next))
			c.Abort()
			c.Next(ctx)
			return
		}
		ctx = context.WithValue(ctx, utils.UserIdKey, userId)
		c.Next(ctx)
	}
}

func AdminAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.Default(c)
		userId := session.Get("user_id")
		if userId == nil {
			hlog.CtxInfof(ctx, "Admin access denied: not logged in")
			c.Redirect(302, []byte("/sign-in?next=/admin"))
			c.Abort()
			return
		}

		uid, ok := userId.(float64)
		if !ok {
			hlog.CtxErrorf(ctx, "Invalid user_id type in session")
			c.Redirect(302, []byte("/sign-in?next=/admin"))
			c.Abort()
			return
		}

		isAdmin, err := rpc.CheckAdmin(ctx, uint32(uid))
		if err != nil {
			hlog.CtxErrorf(ctx, "Failed to check admin status: %v", err)
			c.HTML(500, "error", map[string]interface{}{
				"title":   "错误",
				"message": "服务器内部错误",
			})
			c.Abort()
			return
		}

		if !isAdmin {
			hlog.CtxInfof(ctx, "Admin access denied for user %d: not admin", uint32(uid))
			c.HTML(403, "error", map[string]interface{}{
				"title":   "访问被拒绝",
				"message": "您没有权限访问此页面",
			})
			c.Abort()
			return
		}

		ctx = context.WithValue(ctx, utils.UserIdKey, userId)
		ctx = context.WithValue(ctx, utils.IsAdminKey, true)
		c.Next(ctx)
	}
}
