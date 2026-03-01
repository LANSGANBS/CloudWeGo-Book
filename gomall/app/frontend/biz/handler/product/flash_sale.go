package product

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/biz/service"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/biz/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func FlashSale(ctx context.Context, c *app.RequestContext) {
	resp, err := service.NewFlashSaleService(ctx, c).Run()
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	c.HTML(consts.StatusOK, "flash-sale", utils.WarpResponse(ctx, c, resp))
}
