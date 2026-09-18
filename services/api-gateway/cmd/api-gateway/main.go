package main

import (
	"api-gateway/internal/config"
	"api-gateway/internal/denylist"
	authv1 "api-gateway/internal/gen/auth/v1"
	profilev1 "api-gateway/internal/gen/profile/v1"
	"api-gateway/internal/middleware"
	"api-gateway/internal/router"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	redis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	cfg, err := config.CreateConfig(os.Getenv("HTTP_ADDR"), os.Getenv("GRPC_ADDR_AUTH"), os.Getenv("GRPC_ADDR_PROFILE"), os.Getenv("REDIS_ADDR"), os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatal(err)
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	connAuth, err := grpc.NewClient(cfg.GRPCAddrAuth, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer connAuth.Close()
	connProfile, err := grpc.NewClient(cfg.GRPCAddrProfile, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer connProfile.Close()
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	redisDenyList := denylist.CreateDenyList(redisClient)

	authclient := authv1.NewAuthServiceClient(connAuth)
	profileclient := profilev1.NewProfileServiceClient(connProfile)
	err = http.ListenAndServe(cfg.HttpAddr, middleware.RequestIDToContext(middleware.Logging(router.NewRouter(authclient, profileclient, cfg.JWTSecret, redisDenyList))))
	if err != nil {
		log.Fatal(err)
	}

}
