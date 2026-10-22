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
	"social-graph-service/internal/telemetry"
	"social-graph-service/internal/transport/grpcserver"
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

	cfg, err := config.CreateConfig(os.Getenv("GRPC_ADDR"), os.Getenv("POSTGRES_URL"), os.Getenv("KAFKA_BROKERS"), os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("OTEL_SERVICE_NAME"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	//OpenTelemetry
	shutdown, err := telemetry.Init(ctx, cfg.OtelExporterOtlpEndpoint, cfg.OtelServiceName)
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(context.Background())

	//PgxPool

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
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

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
