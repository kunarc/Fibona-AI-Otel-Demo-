package handler

import (
	"net/http"
	server "otel-demo/grpc-server"
	"otel-demo/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func SendChat(ctx *gin.Context) {
	var server server.ChatServer
	if err := ctx.ShouldBind(&server); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": err.Error()},
		)
		return
	}
	// 创建 gRPC 客户端连接---用户连接远程大模型服务
	res, fa_ctx, err := server.SendChat(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": err.Error()},
		)
		return
	}
	// 开启监控
	tracer := otel.Tracer("save-tarcer")
	_, span := tracer.Start(fa_ctx, "save-span")
	span.AddEvent("将获取到大模型中的数据保存到数据库")
	defer span.End()
	// 模拟阻塞
	time.Sleep(3 * time.Second)
	err = models.DB.Model(&models.Chat{}).Create(&models.Chat{
		Prompt: server.Prompt,
		Reply:  res.Response,
	}).Error
	if err != nil {
		span.SetStatus(codes.Error, "saveToDatabase failed")
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("isTrue", false))
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": err.Error()},
		)
	}
	span.SetAttributes(attribute.Bool("isTrue", true))
	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": res.Response},
	)
}
