package middleware

import (
	"encoding/json"
	"github.com/i4ba1/CustomerOrderAPI/helper"
	"github.com/i4ba1/CustomerOrderAPI/order"
	"log"
	"net/http"
	"time"
)

func CreateOrder(w http.ResponseWriter, r *http.Request){
	newOrder := &order.OrderDto{}
	err := json.NewDecoder(r.Body).Decode(&newOrder)

	if err != nil {
		ErrorResponse(http.StatusUnprocessableEntity, "Invalid JSON", w)
		return
	}

	processNewOrder(newOrder)
}


func processNewOrder(order *order.OrderDto) (int64, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()
	orderNumber := "PO-123"+"/"+helper.Roman(int(time.Now().Month()))+"/"+string(time.Now().Year())
	orderId := helper.GenerateId()
	insertQueryOrder := "INSERT INTO tbl_order (order_id, customer_id, order_number, payment_method_id) VALUES (?,?,?,?)"
	result, err := db.Exec(
		insertQueryOrder,
		orderId,
		order.CustomerId,
		orderNumber,
		"select payment_method_id from tbl_payment_method where payment_method_id="+order.PaymentMethodId)

	if err != nil {
		log.Fatalf("Failed to execute the query insert Order. %v", err)
	}
	row, err := result.RowsAffected()

	if row > 0 {
		insertQueryOrderDetail := "INSERT INTO tbl_order_detail (order_detail_id, order_id, product_id, quantity, created_date) VALUES (?,?,?,?,?)"
		orders := order.Orders
		for i:=0; i<len(orders); i++{
			orderDetailId := helper.GenerateId()
			result, err = db.Exec(insertQueryOrderDetail,
				orderDetailId,
				orderId,
				"select product_id from tbl_product where product_id="+orders[i].ProductId,
				orders[i].Quantity,
				time.Now())
		}

		row, err = result.RowsAffected()
		if err != nil {
			log.Fatalf("Failed to execute the query insert Order Detail. %v", err)
		}
	}

	return row, err
}