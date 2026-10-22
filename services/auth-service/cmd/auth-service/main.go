package main

import (
	"auth-service/internal/config"
	authv1 "auth-service/internal/gen/auth/v1"
	"auth-service/internal/kafka"
	"auth-service/internal/relay"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"
	"auth-service/internal/telemetry"
	"auth-service/internal/transport/grpcserver"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgx-contrib/pgxotel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {

	cfg, err := config.CreateConfig(os.Getenv("GRPC_ADDR"), os.Getenv("POSTGRES_URL"), os.Getenv("JWT_SECRET"), os.Getenv("KAFKA_BROKERS"), os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("OTEL_SERVICE_NAME"))
	if err != nil {
		log.Fatal("err load .env")
	}

	//Context
	ctx := context.Background()
	//OpenTelemetry
	shutdown, err := telemetry.Init(ctx, cfg.OtelExporterOtlpEndpoint, cfg.OtelServiceName)
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(context.Background())
	//pgxPool
	poolConfig, err := pgxpool.ParseConfig(cfg.PostgresURL)
	if err != nil {
		log.Fatal("failed to parse postgres url:", err)
	}

	poolConfig.ConnConfig.Tracer = &pgxotel.QueryTracer{
		Name:             cfg.OtelServiceName, // Имя вашего сервиса в Jaeger
		IncludeStatement: true,                // Обязательно true, если хотите видеть текст SQL-запросов
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal("failed to create pgx pool:", err)
	}

	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//outbox
	outboxRepo := postgres.CreateOutboxRepo(pool)
	ctx2, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	producer := kafka.NewProducer(cfg.KafkaBroker)
	defer producer.Close()
	rel := relay.NewRelay(outboxRepo, 2*time.Second, producer)
	defer stop()
	go rel.Run(ctx2)

	//grpc
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	userrepo := postgres.CreateUserRepo(pool)
	sessionrepo := postgres.CreateSessionRepo(pool)

	registerrepo := postgres.CreateRegisterRepo(pool)
	authSvc := service.NewAuthService(userrepo, sessionrepo, registerrepo, cfg.JWTSecret)
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
