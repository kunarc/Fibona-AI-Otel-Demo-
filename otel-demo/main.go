package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"

	"otel-demo/handler"

	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("go-client-otel"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
func main() {
	r := gin.Default()
	// 测试服务是否正常运行
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/chat", handler.SendChat)
	ctx := context.Background()
	// 配置 OpenTelemetry导出器
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("localhost:4317"))
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	// 确保在程序退出时关闭导出器
	defer func() {
		if err := exporter.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown exporter: %v", err)
		}
	}()

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		// 配置资源属性
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown tracer provider: %v", err)
		}
	}()
	// 设置全局跟踪器提供程序
	otel.SetTracerProvider(tp)
	r.Run(":8081")
}
