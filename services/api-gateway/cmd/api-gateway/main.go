package main

import (
	"api-gateway/internal/config"
	authv1 "api-gateway/internal/gen/auth/v1"
	"api-gateway/internal/router"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.CreateConfig(os.Getenv("HTTP_ADDR"), os.Getenv("GRPC_ADDR_AUTH"), os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.NewClient(cfg.GRPCAddrAuth, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	authclient := authv1.NewAuthServiceClient(conn)

	err = http.ListenAndServe(cfg.HttpAddr, router.NewRouter(authclient, cfg.JWTSecret))
	if err != nil {
		log.Fatal(err)
	}
}
