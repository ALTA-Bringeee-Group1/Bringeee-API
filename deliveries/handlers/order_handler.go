package handlers

import (
	orderService "bringeee-capstone/services/order"
	userService "bringeee-capstone/services/user"
)

type OrderHandler struct {
	orderService orderService.OrderServiceInterface
	userService userService.UserServiceInterface
}

func NewOrderHandler(
	service orderService.OrderServiceInterface,
	userService userService.UserServiceInterface,
) *OrderHandler {
	return &OrderHandler{
		orderService: service,
		userService: userService,
	}
}