package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"profile-service/internal/config"
	"profile-service/internal/consumer"
	"profile-service/internal/repository/postgres"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	//Config
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.CreateConfig(os.Getenv("GRPC_ADDR"), os.Getenv("POSTGRES_URL"), os.Getenv("KAFKA_BROKERS"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	//PgxPool
	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	profilerepo := postgres.CreateProfileRepo(pool)
	c := consumer.CreateConsumer(profilerepo, "profile-service", cfg.KafkaBroker)
	ctx2, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go c.Run(ctx2)
	defer c.Close()

	//grpc
	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening", cfg.GrpcAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
