module github.com/pastelnetwork/gonode/tools/pastel-api

go 1.16

require github.com/pastelnetwork/gonode/common v0.0.0

replace (
    github.com/pastelnetwork/gonode/common => ../../common
    replace github.com/pastelnetwork/gonode/proto => ../../proto
)
