package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"social-graph-service/internal/config"
	socialgraphv1 "social-graph-service/internal/gen/social-graph/v1"
	"social-graph-service/internal/kafka"
	"social-graph-service/internal/relay"
	"social-graph-service/internal/repository/postgres"
	"social-graph-service/internal/service"
	"social-graph-service/internal/transport/grpcserver"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {

	cfg, err := config.CreateConfig(os.Getenv("GRPC_ADDR"), os.Getenv("POSTGRES_URL"), os.Getenv("KAFKA_BROKERS"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	graphRepo := postgres.CreateGraphRepo(pool)
	outboxRepo := postgres.CreateOutboxRepo(pool)
	producer := kafka.NewProducer(cfg.KafkaBroker)
	defer producer.Close()

	ctxRelay, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go relay.NewRelay(outboxRepo, 2*time.Second, producer).Run(ctxRelay)
	svc := service.NewGraphService(graphRepo)
	server := &grpcserver.Server{GraphService: svc}
	//grpc
	grpcServer := grpc.NewServer()

	socialgraphv1.RegisterSocialGraphServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening", cfg.GRPCAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
