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

package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/infra/mq"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	checkout "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/checkout"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/email"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/payment"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
)

type CheckoutService struct {
	ctx context.Context
} // NewCheckoutService new CheckoutService
func NewCheckoutService(ctx context.Context) *CheckoutService {
	return &CheckoutService{ctx: ctx}
}

/*
Run

1. get cart
2. calculate cart & check stock
3. deduct stock
4. create order
5. empty cart
6. pay
7. change order result
8. increment sales
9. finish (restore stock if failed)
*/
func (s *CheckoutService) Run(req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	// Finish your business logic.
	// Idempotent
	// get cart
	cartResult, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{UserId: req.UserId})
	if err != nil {
		klog.Error(err)
		err = fmt.Errorf("GetCart.err:%v", err)
		return
	}
	if cartResult == nil || cartResult.Cart == nil || len(cartResult.Cart.Items) == 0 {
		err = errors.New("cart is empty")
		return
	}
	var (
		oi              []*order.OrderItem
		total           float32
		productSalesMap = make(map[uint32]int32)
	)
	for _, cartItem := range cartResult.Cart.Items {
		productResp, resultErr := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{Id: cartItem.ProductId})
		if resultErr != nil {
			klog.Error(resultErr)
			err = resultErr
			return
		}
		if productResp.Product == nil {
			continue
		}
		p := productResp.Product
		
		// Check stock availability
		if p.Stock < int64(cartItem.Quantity) {
			err = fmt.Errorf("insufficient stock for product %s: available %d, requested %d", p.Name, p.Stock, cartItem.Quantity)
			return
		}
		
		cost := p.Price * float32(cartItem.Quantity)
		total += cost
		oi = append(oi, &order.OrderItem{
			Item: &cart.CartItem{ProductId: cartItem.ProductId, Quantity: cartItem.Quantity},
			Cost: cost,
		})
		productSalesMap[cartItem.ProductId] = cartItem.Quantity
	}
	
	// Deduct stock before creating order
	for productId, quantity := range productSalesMap {
		deductResp, deductErr := rpc.ProductClient.DeductStock(s.ctx, &product.DeductStockReq{
			ProductId: productId,
			Quantity:  quantity,
		})
		if deductErr != nil || !deductResp.Success {
			// Restore already deducted stock
			for pid, qty := range productSalesMap {
				if pid == productId {
					break
				}
				_, _ = rpc.ProductClient.RestoreStock(s.ctx, &product.RestoreStockReq{
					ProductId: pid,
					Quantity:  qty,
				})
			}
			err = fmt.Errorf("DeductStock.err: productId=%d, quantity=%d, err=%v, msg=%s", 
				productId, quantity, deductErr, deductResp.GetErrorMessage())
			return
		}
		klog.Infof("Deducted stock: productId=%d, quantity=%d, remaining=%d", productId, quantity, deductResp.RemainingStock)
	}
	
	// create order
	orderReq := &order.PlaceOrderReq{
		UserId:       req.UserId,
		UserCurrency: "USD",
		OrderItems:   oi,
		Email:        req.Email,
	}
	if req.Address != nil {
		addr := req.Address
		zipCodeInt, _ := strconv.Atoi(addr.ZipCode)
		orderReq.Address = &order.Address{
			StreetAddress: addr.StreetAddress,
			City:          addr.City,
			Country:       addr.Country,
			State:         addr.State,
			ZipCode:       int32(zipCodeInt),
		}
	}
	orderResult, err := rpc.OrderClient.PlaceOrder(s.ctx, orderReq)
	if err != nil {
		// Restore stock if order creation fails
		for productId, quantity := range productSalesMap {
			_, _ = rpc.ProductClient.RestoreStock(s.ctx, &product.RestoreStockReq{
				ProductId: productId,
				Quantity:  quantity,
			})
		}
		err = fmt.Errorf("PlaceOrder.err:%v", err)
		return
	}
	klog.Info("orderResult", orderResult)
	// empty cart
	emptyResult, err := rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{UserId: req.UserId})
	if err != nil {
		err = fmt.Errorf("EmptyCart.err:%v", err)
		return
	}
	klog.Info(emptyResult)
	// charge
	var orderId string
	if orderResult != nil || orderResult.Order != nil {
		orderId = orderResult.Order.OrderId
	}
	payReq := &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: orderId,
		Amount:  total,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
		},
	}
	paymentResult, err := rpc.PaymentClient.Charge(s.ctx, payReq)
	if err != nil {
		// Restore stock if payment fails
		for productId, quantity := range productSalesMap {
			_, _ = rpc.ProductClient.RestoreStock(s.ctx, &product.RestoreStockReq{
				ProductId: productId,
				Quantity:  quantity,
			})
		}
		err = fmt.Errorf("Charge.err:%v", err)
		return
	}
	data, _ := proto.Marshal(&email.EmailReq{
		From:        "from@example.com",
		To:          req.Email,
		ContentType: "text/plain",
		Subject:     "You just created an order in CloudWeGo shop",
		Content:     "You just created an order in CloudWeGo shop",
	})
	msg := &nats.Msg{Subject: "email", Data: data, Header: make(nats.Header)}

	// otel inject
	otel.GetTextMapPropagator().Inject(s.ctx, propagation.HeaderCarrier(msg.Header))

	_ = mq.Nc.PublishMsg(msg)

	klog.Info(paymentResult)
	// change order state
	klog.Info(orderResult)
	_, err = rpc.OrderClient.MarkOrderPaid(s.ctx, &order.MarkOrderPaidReq{UserId: req.UserId, OrderId: orderId})
	if err != nil {
		klog.Error(err)
		return
	}

	// increment sales for each product
	for productId, quantity := range productSalesMap {
		_, salesErr := rpc.ProductClient.IncrementSales(s.ctx, &product.IncrementSalesReq{
			Id:       productId,
			Quantity: quantity,
		})
		if salesErr != nil {
			klog.Errorf("IncrementSales.err: productId=%d, quantity=%d, err=%v", productId, quantity, salesErr)
		}
	}

	resp = &checkout.CheckoutResp{
		OrderId:       orderId,
		TransactionId: paymentResult.TransactionId,
	}
	return
}
