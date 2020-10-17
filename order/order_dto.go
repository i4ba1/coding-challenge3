package order

import "time"

type OrderDto struct {
	CustomerId string
	OrderNumber string
	PaymentMethodId string
	Orders []OrderDetailDto
}

type OrderDetailDto struct {
	OrderDetailId string
	OrderId string
	ProductId string
	Quantity int
	CreatedDate time.Time
}
