module github.com/janoszen/openshiftci-inspector

go 1.14


require (
	github.com/aws/aws-sdk-go v1.36.13
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/mux v1.7.3
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/prometheus v1.0.0
	k8s.io/api v0.20.0
)

// Note on Prometheus dependency: Prometheus is not built with the Go module spec in mind.
// We are therefore faking version 1.0.0 and replacing it with the specific version we want.
// 2.21.0 -> v0.0.0-20200911110723-e83ef207b6c2
replace github.com/prometheus/prometheus v1.0.0 => github.com/prometheus/prometheus v0.0.0-20200911110723-e83ef207b6c2
