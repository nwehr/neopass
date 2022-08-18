package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jackc/pgx/v4"
)

func main() {
	runtime.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	influxClient := influxdb2.NewClient(os.Getenv("INFLUX_URL"), os.Getenv("INFLUX_TOKEN"))
	defer influxClient.Close()

	influx := influxClient.WriteAPI(os.Getenv("INFLUX_ORG"), "neopass")

	defer func() {
		influx.Flush()
	}()

	clientUUID, ok := request.QueryStringParameters["client_uuid"]
	if !ok {
		influx.WriteRecord("client_err value=true")
		return _err(400, "missing client_uuid parameter")
	}

	cockroach, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		influx.WriteRecord(fmt.Sprintf("err,client_uuid=%s value=true", clientUUID))
		return _err(500, "could not connect to database")
	}

	defer cockroach.Close(ctx)

	switch request.HTTPMethod {
	case "GET":
		name, ok := request.QueryStringParameters["name"]

		if ok {
			influx.WriteRecord(fmt.Sprintf("get_entry,client_uuid=%s value=true", clientUUID))

			resp, err := getEntry(ctx, cockroach, clientUUID, name)
			if err != nil {
				influx.WriteRecord("err value=true")
			}

			return resp, err

		} else {
			influx.WriteRecord(fmt.Sprintf("list_entry_names,client_uuid=%s value=true", clientUUID))

			resp, err := listEntryNames(ctx, cockroach, clientUUID)
			if err != nil {
				influx.WriteRecord("err value=true")
			}
			return resp, err
		}
	case "POST":
		influx.WriteRecord(fmt.Sprintf("set_entry,client_uuid=%s value=true", clientUUID))

		entry := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{}

		err := json.Unmarshal([]byte(request.Body), &entry)
		if err != nil {
			influx.WriteRecord(fmt.Sprintf("err,client_uuid=%s value=true", clientUUID))
			return _err(500, "could not parse json request")
		}

		resp, err := setEntry(ctx, cockroach, clientUUID, entry.Name, entry.Password)
		if err != nil {
			influx.WriteRecord(fmt.Sprintf("err,client_uuid=%s value=true", clientUUID))
		}
		return resp, err

	case "DELETE":
		influx.WriteRecord(fmt.Sprintf("delete_entry,client_uuid=%s value=true", clientUUID))

		name, ok := request.QueryStringParameters["name"]
		if ok {
			resp, err := deleteEntry(ctx, cockroach, clientUUID, name)
			if err != nil {
				influx.WriteRecord("err value=true")
			}
			return resp, err
		}

		return _err(400, "missing name")
	}

	return _err(400, "bad request")
}

func getEntry(ctx context.Context, conn *pgx.Conn, clientUUID string, name string) (events.APIGatewayProxyResponse, error) {
	entry := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{
		Name: name,
	}

	err := conn.QueryRow(ctx, `select revisions[current_revision] as "password" from entries where "client_uuid" = $1 and "name" = $2`, clientUUID, entry.Name).Scan(&entry.Password)
	if err != nil {
		return _err(500, "could not query entries")
	}

	encoded, err := json.Marshal(entry)
	if err != nil {
		return _err(500, "could not encode json response")
	}

	return _ok(string(encoded))
}

func listEntryNames(ctx context.Context, conn *pgx.Conn, clientUUID string) (events.APIGatewayProxyResponse, error) {
	rows, err := conn.Query(ctx, `select "name" from entries where "client_uuid" = $1`, clientUUID)
	if err != nil {
		return _err(500, "could not query entries")
	}

	defer rows.Close()

	names := []string{}

	for rows.Next() {
		name := ""

		err = rows.Scan(&name)
		if err != nil {
			return _err(500, "could not scan entry names")
		}

		names = append(names, name)
	}

	encoded, err := json.Marshal(names)
	if err != nil {
		return _err(500, "could not encode json response")
	}

	return _ok(string(encoded))
}

func setEntry(ctx context.Context, conn *pgx.Conn, clientUUID string, name string, password string) (events.APIGatewayProxyResponse, error) {
	q := `insert into "entries" (
		"client_uuid"
		, "name"
		, "revisions"
		, "current_revision"
	) 
	values (
		$1
		, $2
		, array[$3]
		, 1
	) 
	on conflict on constraint client_uuid_name_unique do 
		update set 
			revisions = array_append(entries.revisions, $3)
			, current_revision = entries.current_revision + 1
	`

	_, err := conn.Exec(ctx, q, clientUUID, name, password)
	if err != nil {
		return _err(500, "could not insert into entries")
	}

	return _ok("")
}

func deleteEntry(ctx context.Context, conn *pgx.Conn, clientUUID string, name string) (events.APIGatewayProxyResponse, error) {
	_, err := conn.Exec(ctx, `delete from entries where "client_uuid" = $1 and "name" = $2`, clientUUID, name)
	if err != nil {
		return _err(500, "could not delete from entries")
	}

	return _ok("")
}

func _err(code int, err string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: err, StatusCode: code}, fmt.Errorf(err)
}

func _ok(body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: body, StatusCode: 200}, nil
}
