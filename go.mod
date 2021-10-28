module github.com/nwehr/npass

go 1.16

replace golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad => ../crypto

require (
	filippo.io/age v1.0.0 // indirect
	github.com/ProtonMail/gopenpgp/v2 v2.1.5 // indirect
	github.com/atotto/clipboard v0.1.2
	github.com/go-piv/piv-go v1.9.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.9.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
