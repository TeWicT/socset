package main

import (
	"api-gateway/internal/config"
	"api-gateway/internal/denylist"
	authv1 "api-gateway/internal/gen/auth/v1"
	profilev1 "api-gateway/internal/gen/profile/v1"
	"api-gateway/internal/middleware"
	"api-gateway/internal/router"
	"api-gateway/internal/telemetry"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	redis "github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//Config
	cfg, err := config.CreateConfig(os.Getenv("HTTP_ADDR"), os.Getenv("GRPC_ADDR_AUTH"), os.Getenv("GRPC_ADDR_PROFILE"), os.Getenv("REDIS_ADDR"), os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("OTEL_SERVICE_NAME"), os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	//Logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	//Conn
	connAuth, err := grpc.NewClient(cfg.GRPCAddrAuth, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Fatal(err)
	}
	defer connAuth.Close()
	connProfile, err := grpc.NewClient(cfg.GRPCAddrProfile, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Fatal(err)
	}
	defer connProfile.Close()

	//Context
	ctx := context.Background()

	//Redis
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	redisDenyList := denylist.CreateDenyList(redisClient)

	//gRPC client
	authclient := authv1.NewAuthServiceClient(connAuth)
	profileclient := profilev1.NewProfileServiceClient(connProfile)
	//OpenTelemetry
	shutdown, err := telemetry.Init(ctx, cfg.OtelExporterOtlpEndpoint, cfg.OtelServiceName)
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(context.Background())

	//Http
	err = http.ListenAndServe(cfg.HttpAddr, otelhttp.NewHandler(middleware.RequestIDToContext(middleware.Logging(router.NewRouter(authclient, profileclient, cfg.JWTSecret, redisDenyList))), cfg.OtelServiceName))
	if err != nil {
		log.Fatal(err)
	}

}
