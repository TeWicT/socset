package main

import (
	"auth-service/internal/config"
	authv1 "auth-service/internal/gen/auth/v1"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"
	"auth-service/internal/transport/grpcserver"
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {

	// Config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err load .env")
	}
	cfg, err := config.CreateConfig(os.Getenv("GRPC_ADDR"), os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("err load .env")
	}

	//Context
	ctx := context.Background()

	//pgxPool
	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	err = pool.Ping(ctx)

	if err != nil {
		log.Fatal(err)
	}

	//grpc
	grpcServer := grpc.NewServer()
	repo := postgres.CreateUserRepo(pool)
	authSvc := service.NewAuthService(repo)
	srv := &grpcserver.Server{Auth: authSvc}
	authv1.RegisterAuthServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	//Health
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	//Listen
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening", cfg.GRPCAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

}
