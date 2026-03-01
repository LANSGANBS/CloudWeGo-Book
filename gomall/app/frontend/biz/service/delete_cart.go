package service

import (
	"context"
	"strconv"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/infra/rpc"
	frontendutils "github.com/cloudwego/biz-demo/gomall/app/frontend/utils"
	rpccart "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteCartService(Context context.Context, RequestContext *app.RequestContext) *DeleteCartService {
	return &DeleteCartService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteCartService) Run(productId string) error {
	productIdInt, err := strconv.ParseUint(productId, 10, 32)
	if err != nil {
		return err
	}
	
	_, err = rpc.CartClient.AddItem(h.Context, &rpccart.AddItemReq{
		UserId: frontendutils.GetUserIdFromCtx(h.Context),
		Item: &rpccart.CartItem{
			ProductId: uint32(productIdInt),
			Quantity:  0,
		},
	})
	return err
}
