module github.com/aptible/cloud-cli

go 1.17

require github.com/aptible/cloud-api-clients v0.0.0

require (
	github.com/golang/protobuf v1.4.2 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/oauth2 v0.0.0-20210323180902-22b0adad7558 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

replace github.com/aptible/cloud-api-clients v0.0.0 => ../cloud-api-clients/clients/go
