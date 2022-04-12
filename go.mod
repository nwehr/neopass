module github.com/nwehr/neopass

go 1.18

replace golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad => ../crypto

require (
	filippo.io/age v1.0.0
	github.com/atotto/clipboard v0.1.2
	github.com/aws/aws-lambda-go v1.29.0
	github.com/go-piv/piv-go v1.9.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgx/v4 v4.15.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.11.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.10.0 // indirect
	golang.org/x/sys v0.0.0-20210903071746-97244b99971b // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	golang.org/x/text v0.3.6 // indirect
)
