package service

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type FlashSaleService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewFlashSaleService(Context context.Context, RequestContext *app.RequestContext) *FlashSaleService {
	return &FlashSaleService{RequestContext: RequestContext, Context: Context}
}

func (s *FlashSaleService) Run() (res map[string]any, err error) {
	ctx := s.Context
	
	resp, err := rpc.ProductClient.ListProducts(ctx, &product.ListProductsReq{DiscountFilter: 3})
	if err != nil {
		hlog.CtxErrorf(ctx, "FlashSaleService: ListProducts RPC error: %v", err)
		return utils.H{
			"title": "限时特惠",
			"items": []*product.Product{},
		}, nil
	}

	hlog.CtxInfof(ctx, "FlashSaleService: found %d flash sale products", len(resp.Products))

	return utils.H{
		"title": "限时特惠",
		"items": resp.Products,
	}, nil
}
