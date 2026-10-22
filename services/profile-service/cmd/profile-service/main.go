package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"profile-service/internal/config"
	"profile-service/internal/consumer"
	profilev1 "profile-service/internal/gen/profile/v1"
	"profile-service/internal/repository/postgres"
	"profile-service/internal/service"
	"profile-service/internal/telemetry"
	"profile-service/internal/transport/grpcserver"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgx-contrib/pgxotel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
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
	profilerepo := postgres.CreateProfileRepo(pool)
	c := consumer.CreateConsumer(profilerepo, "profile-service", cfg.KafkaBroker)
	ctx2, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go c.Run(ctx2)
	defer c.Close()
	profileService := service.CreateProfileService(profilerepo)
	srv := &grpcserver.Server{Profiles: profileService}
	//grpc
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	profilev1.RegisterProfileServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening", cfg.GrpcAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
