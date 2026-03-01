package main

import (
	"context"
	product "github.com/cloudwego/biz-demo/gomall/kitex_gen/product"
	cart "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart/cart"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product/product"
)

// ProductCatalogServiceImpl implements the last service interface defined in the IDL.
type ProductCatalogServiceImpl struct{}

// ListProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	// TODO: Your code here...
	return
}

// GetProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) GetProduct(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	// TODO: Your code here...
	return
}

// SearchProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	// TODO: Your code here...
	return
}

// CreateProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) CreateProduct(ctx context.Context, req *product.CreateProductReq) (resp *product.CreateProductResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) UpdateProduct(ctx context.Context, req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	// TODO: Your code here...
	return
}

// AddItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	// TODO: Your code here...
	return
}

// GetCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	// TODO: Your code here...
	return
}

// EmptyCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) DeleteItem(ctx context.Context, req *cart.DeleteItemReq) (resp *cart.DeleteItemResp, err error) {
	// TODO: Your code here...
	return
}

// IncrementSales implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) IncrementSales(ctx context.Context, req *product.IncrementSalesReq) (resp *product.IncrementSalesResp, err error) {
	// TODO: Your code here...
	return
}

// DeductStock implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) DeductStock(ctx context.Context, req *product.DeductStockReq) (resp *product.DeductStockResp, err error) {
	// TODO: Your code here...
	return
}

// RestoreStock implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) RestoreStock(ctx context.Context, req *product.RestoreStockReq) (resp *product.RestoreStockResp, err error) {
	// TODO: Your code here...
	return
}

// GetStock implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) GetStock(ctx context.Context, req *product.GetStockReq) (resp *product.GetStockResp, err error) {
	// TODO: Your code here...
	return
}
