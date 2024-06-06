package server

import (
	"context"
	"log"
	pb "otel-demo/proto"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatServer struct {
	Prompt string `json:"prompt"`
}

func (server *ChatServer) SendChat(ctx context.Context) (*pb.ChatResponse, context.Context, error) {
	// 创建 gRPC 客户端连接---用户连接远程大模型服务
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, nil, err
	}
	defer conn.Close()
	// 创建 gRPC 客户端
	client := pb.NewSendChatClient(conn)
	// 调用服务
	tracer := otel.Tracer("call-tarcer")
	new_ctx, span := tracer.Start(ctx, "call-span")
	// 设置事件和属性
	span.SetAttributes(attribute.String("prompt", server.Prompt))
	span.AddEvent("调用grpc大模型服务")
	defer span.End()
	req := &pb.ChatRequest{Message: server.Prompt}
	time.Sleep(2 * time.Second)
	res, err := client.SendChat(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, "operationThatCouldFail failed")
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("isSuccessCall", false))
		return nil, new_ctx, err
	}
	span.SetAttributes(attribute.Bool("isSuccessCall", true))
	return res, new_ctx, nil
}
