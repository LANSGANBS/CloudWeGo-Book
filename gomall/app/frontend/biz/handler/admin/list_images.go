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

package admin

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type ImageInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
	URL  string `json:"url"`
}

func ListImages(ctx context.Context, c *app.RequestContext) {
	imageDir := "static/image"

	var images []ImageInfo

	err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
		if !validExts[ext] {
			return nil
		}

		relPath, err := filepath.Rel(imageDir, path)
		if err != nil {
			return nil
		}

		images = append(images, ImageInfo{
			Name: info.Name(),
			Path: relPath,
			URL:  "/static/image/" + relPath,
		})
		return nil
	})

	if err != nil {
		c.JSON(consts.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Failed to read image directory",
		})
		return
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].Name < images[j].Name
	})

	c.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"images":  images,
	})
}
