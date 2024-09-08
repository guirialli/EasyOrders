package main

import (
	"database/sql"
	"fmt"
	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/guirialli/go-pos/clean_arch/configs"
	"github.com/guirialli/go-pos/clean_arch/internals/event/handler"
	"github.com/guirialli/go-pos/clean_arch/internals/infra/graph"
	"github.com/guirialli/go-pos/clean_arch/internals/infra/grpc/pb"
	"github.com/guirialli/go-pos/clean_arch/internals/infra/grpc/service"
	"github.com/guirialli/go-pos/clean_arch/internals/infra/web/webserver"
	"github.com/guirialli/go-pos/clean_arch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

func main() {
	conf, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	urlConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName)
	db, err := sql.Open(conf.DBDriver, urlConnection)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	eventDispatcher.Register("OrderListed", &handler.OrdersListedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUserCase := NewCreateOrderUseCase(db, eventDispatcher)
	findAllOrderUserCase := NewFindAllOrdersUseCase(db, eventDispatcher)

	server := webserver.NewWebServer(conf.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	server.AddHandler("/order", webOrderHandler.Create)
	server.AddHandler("/find", webOrderHandler.FindAll)
	fmt.Println("Starting web server on port", conf.WebServerPort)
	go server.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUserCase, *findAllOrderUserCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)
	fmt.Println("Starting gRPC server on port", conf.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase:   *createOrderUserCase,
		FindAllOrdersUseCase: *findAllOrderUserCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", conf.GraphQLServerPort)
	http.ListenAndServe(":"+conf.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
