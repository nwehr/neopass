package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v4"
)

func main() {
	runtime.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		Fatalf("%v\n", err)
	}

	defer conn.Close(ctx)

	clientUUID, ok := request.QueryStringParameters["client_uuid"]
	if !ok {
		return events.APIGatewayProxyResponse{Body: "missing client_uuid", StatusCode: 400}, nil
	}

	switch request.HTTPMethod {
	case "GET":
		name, ok := request.QueryStringParameters["name"]

		if ok {
			return getEntryHandler(ctx, conn, clientUUID, name)

		} else {
			return listEntryNameshandler(ctx, conn, clientUUID)
		}
	case "POST":
		entry := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{}

		err := json.Unmarshal([]byte(request.Body), &entry)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
		}

		return addEntryHandler(ctx, conn, clientUUID, entry.Name, entry.Password)

	case "DELETE":
		name, ok := request.QueryStringParameters["name"]
		if ok {
			return deleteEntryHandler(ctx, conn, clientUUID, name)
		}

		return events.APIGatewayProxyResponse{Body: "missing name", StatusCode: 400}, nil
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func getEntryHandler(ctx context.Context, conn *pgx.Conn, clientUUID string, name string) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("getEntryHandler(%s)\n", name)

	entry := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{
		Name: name,
	}

	err := conn.QueryRow(ctx, `select "password" from entries where "client_uuid" = $1 and "name" = $2`, clientUUID, entry.Name).Scan(&entry.Password)
	if err != nil {
		fmt.Printf("getEntryHandler(%s) %v\n", name, err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	encoded, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("getEntryHandler(%s) %v\n", name, err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: string(encoded), StatusCode: 200}, nil
}

func listEntryNameshandler(ctx context.Context, conn *pgx.Conn, clientUUID string) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("listEntryNameshandler(%s)\n", clientUUID)

	rows, err := conn.Query(ctx, `select "name" from entries where "client_uuid" = $1`, clientUUID)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	defer rows.Close()

	names := []string{}

	for rows.Next() {
		name := ""

		err = rows.Scan(&name)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
		}

		names = append(names, name)
	}

	fmt.Printf("listEntryNameshandler(%s) %d results\n", clientUUID, len(names))

	encoded, err := json.Marshal(names)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: string(encoded), StatusCode: 200}, nil
}

func addEntryHandler(ctx context.Context, conn *pgx.Conn, clientUUID string, name string, password string) (events.APIGatewayProxyResponse, error) {
	_, err := conn.Exec(ctx, `insert into entries ("client_uuid", "name", "password") values ($1, $2, $3)`, clientUUID, name, password)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func deleteEntryHandler(ctx context.Context, conn *pgx.Conn, clientUUID string, name string) (events.APIGatewayProxyResponse, error) {
	_, err := conn.Exec(ctx, `delete from entries where "client_uuid" = $1 and "name" = $2`, clientUUID, name)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, err
}

func Fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
