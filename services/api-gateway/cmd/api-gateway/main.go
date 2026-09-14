package main

import (
	"api-gateway/internal/config"
	authv1 "api-gateway/internal/gen/auth/v1"
	profilev1 "api-gateway/internal/gen/profile/v1"
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
	cfg, err := config.CreateConfig(os.Getenv("HTTP_ADDR"), os.Getenv("GRPC_ADDR_AUTH"), os.Getenv("GRPC_ADDR_PROFILE"), os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatal(err)
	}
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

	authclient := authv1.NewAuthServiceClient(connAuth)
	profileclient := profilev1.NewProfileServiceClient(connProfile)
	err = http.ListenAndServe(cfg.HttpAddr, router.NewRouter(authclient, profileclient, cfg.JWTSecret))
	if err != nil {
		log.Fatal(err)
	}
}
