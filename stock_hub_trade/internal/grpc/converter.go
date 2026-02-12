package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"stock_hub_trade/api/proto"
	"stock_hub_trade/pkg/model"
)

// ToProtoUser converts domain User to proto User
func ToProtoUser(user *model.User) *proto.User {
	if user == nil {
		return nil
	}
	return &proto.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

// FromProtoUser converts proto User to domain User
func FromProtoUser(protoUser *proto.User) *model.User {
	if protoUser == nil {
		return nil
	}
	return &model.User{
		ID:        protoUser.Id,
		Username:  protoUser.Username,
		Email:     protoUser.Email,
		CreatedAt: protoUser.CreatedAt.AsTime(),
	}
}

// ToProtoMessage converts domain Message to proto Message
func ToProtoMessage(msg *model.Message) *proto.Message {
	if msg == nil {
		return nil
	}
	return &proto.Message{
		Id:        msg.ID,
		Content:   msg.Content,
		ClientId:  msg.ClientID,
		CreatedAt: timestamppb.New(msg.CreatedAt),
	}
}

// FromProtoMessage converts proto Message to domain Message
func FromProtoMessage(protoMsg *proto.Message) *model.Message {
	if protoMsg == nil {
		return nil
	}
	return &model.Message{
		ID:        protoMsg.Id,
		Content:   protoMsg.Content,
		ClientID:  protoMsg.ClientId,
		CreatedAt: protoMsg.CreatedAt.AsTime(),
	}
}

// ToProtoStockOrder converts domain StockOrder to proto StockOrder
func ToProtoStockOrder(order *model.StockOrder) *proto.StockOrder {
	if order == nil {
		return nil
	}

	var orderType proto.OrderType
	switch order.OrderType {
	case model.OrderTypeBid:
		orderType = proto.OrderType_ORDER_TYPE_BID
	case model.OrderTypeAsk:
		orderType = proto.OrderType_ORDER_TYPE_ASK
	default:
		orderType = proto.OrderType_ORDER_TYPE_UNSPECIFIED
	}

	return &proto.StockOrder{
		Id:        order.ID,
		UserId:    order.UserID,
		Username:  order.Username,
		Symbol:    order.Symbol,
		OrderType: orderType,
		Price:     order.Price,
		Quantity:  int32(order.Quantity),
		Status:    order.Status,
		CreatedAt: timestamppb.New(order.CreatedAt),
	}
}

// FromProtoStockOrder converts proto StockOrder to domain StockOrder
func FromProtoStockOrder(protoOrder *proto.StockOrder) *model.StockOrder {
	if protoOrder == nil {
		return nil
	}

	var orderType model.OrderType
	switch protoOrder.OrderType {
	case proto.OrderType_ORDER_TYPE_BID:
		orderType = model.OrderTypeBid
	case proto.OrderType_ORDER_TYPE_ASK:
		orderType = model.OrderTypeAsk
	}

	return &model.StockOrder{
		ID:        protoOrder.Id,
		UserID:    protoOrder.UserId,
		Username:  protoOrder.Username,
		Symbol:    protoOrder.Symbol,
		OrderType: orderType,
		Price:     protoOrder.Price,
		Quantity:  int(protoOrder.Quantity),
		Status:    protoOrder.Status,
		CreatedAt: protoOrder.CreatedAt.AsTime(),
	}
}

// ToProtoUsers converts slice of domain Users to proto Users
func ToProtoUsers(users []*model.User) []*proto.User {
	result := make([]*proto.User, len(users))
	for i, user := range users {
		result[i] = ToProtoUser(user)
	}
	return result
}

// ToProtoMessages converts slice of domain Messages to proto Messages
func ToProtoMessages(messages []*model.Message) []*proto.Message {
	result := make([]*proto.Message, len(messages))
	for i, msg := range messages {
		result[i] = ToProtoMessage(msg)
	}
	return result
}

// ToProtoStockOrders converts slice of domain StockOrders to proto StockOrders
func ToProtoStockOrders(orders []*model.StockOrder) []*proto.StockOrder {
	result := make([]*proto.StockOrder, len(orders))
	for i, order := range orders {
		result[i] = ToProtoStockOrder(order)
	}
	return result
}
