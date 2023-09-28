package jsonrpc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/onsi/gomega"
)

// needed to retrieve requests that arrived at httpServer for further investigation
var requestChan = make(chan *RequestData, 1)

// the request datastructure that can be retrieved for test assertions
type RequestData struct {
	request *http.Request
	body    string
}

// set the response body the httpServer should return for the next request
var responseBody = ""

var httpServer *httptest.Server

// start the testhttp server and stop it when tests are finished
func TestMain(m *testing.M) {
	httpServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		// put request and body to channel for the client to investigate them
		requestChan <- &RequestData{r, string(data)}

		fmt.Fprint(w, responseBody)
	}))
	defer httpServer.Close()

	os.Exit(m.Run())
}

func TestSimpleRpcCallHeaderCorrect(t *testing.T) {
	gomega.RegisterTestingT(t)

	rpcClient := NewClient(httpServer.URL)
	rpcClient.Call("add", 1, 2)

	req := (<-requestChan).request

	gomega.Expect(req.Method).To(gomega.Equal("POST"))
	gomega.Expect(req.Header.Get("Content-Type")).To(gomega.Equal("application/json"))
	gomega.Expect(req.Header.Get("Accept")).To(gomega.Equal("application/json"))
}

// test if the structure of an rpc request is built correctly by validating the data that arrived on the test server
func TestRpcClient_Call(t *testing.T) {
	gomega.RegisterTestingT(t)
	rpcClient := NewClient(httpServer.URL)

	person := Person{
		Name:    "Alex",
		Age:     35,
		Country: "Germany",
	}

	drink := Drink{
		Name:        "Cuba Libre",
		Ingredients: []string{"rum", "cola"},
	}

	rpcClient.Call("missingParam")
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"missingParam","id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("nullParam", nil)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"nullParam","params":[null],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("nullParams", nil, nil)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"nullParams","params":[null,null],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("emptyParams", []interface{}{})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"emptyParams","params":[],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("emptyAnyParams", []string{})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"emptyAnyParams","params":[],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("emptyObject", struct{}{})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"emptyObject","params":{},"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("emptyObjectList", []struct{}{{}, {}})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"emptyObjectList","params":[{},{}],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("boolParam", true)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"boolParam","params":[true],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("boolParams", true, false, true)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"boolParams","params":[true,false,true],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("stringParam", "Alex")
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"stringParam","params":["Alex"],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("stringParams", "JSON", "RPC")
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"stringParams","params":["JSON","RPC"],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("numberParam", 123)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"numberParam","params":[123],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("numberParams", 123, 321)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"numberParams","params":[123,321],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("floatParam", 1.23)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"floatParam","params":[1.23],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("floatParams", 1.23, 3.21)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"floatParams","params":[1.23,3.21],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("manyParams", "Alex", 35, true, nil, 2.34)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"manyParams","params":["Alex",35,true,null,2.34],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("emptyMissingPublicFieldObject", struct{ name string }{name: "Alex"})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"emptyMissingPublicFieldObject","params":{},"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("singleStruct", person)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"singleStruct","params":{"name":"Alex","age":35,"country":"Germany"},"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("singlePointerToStruct", &person)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"singlePointerToStruct","params":{"name":"Alex","age":35,"country":"Germany"},"id":0,"jsonrpc":"2.0"}`))

	pp := &person
	rpcClient.Call("doublePointerStruct", &pp)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"doublePointerStruct","params":{"name":"Alex","age":35,"country":"Germany"},"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("multipleStructs", person, &drink)
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"multipleStructs","params":[{"name":"Alex","age":35,"country":"Germany"},{"name":"Cuba Libre","ingredients":["rum","cola"]}],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("singleStructInArray", []interface{}{person})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"singleStructInArray","params":[{"name":"Alex","age":35,"country":"Germany"}],"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("namedParameters", map[string]interface{}{
		"name": "Alex",
		"age":  35,
	})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"namedParameters","params":{"age":35,"name":"Alex"},"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("anonymousStructNoTags", struct {
		Name string
		Age  int
	}{"Alex", 33})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"anonymousStructNoTags","params":{"Name":"Alex","Age":33},"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("anonymousStructWithTags", struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{"Alex", 33})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"anonymousStructWithTags","params":{"name":"Alex","age":33},"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("structWithNullField", struct {
		Name    string  `json:"name"`
		Address *string `json:"address"`
	}{"Alex", nil})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"structWithNullField","params":{"name":"Alex","address":null},"id":0,"jsonrpc":"2.0"}`))

	rpcClient.Call("nestedStruct",
		Planet{
			Name: "Mars",
			Properties: Properties{
				Distance: 54600000,
				Color:    "red",
			},
		})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`{"method":"nestedStruct","params":{"name":"Mars","properties":{"distance":54600000,"color":"red"}},"id":0,"jsonrpc":"2.0"}`))
}

func TestRpcClient_CallBatch(t *testing.T) {
	gomega.RegisterTestingT(t)
	rpcClient := NewClient(httpServer.URL)

	person := Person{
		Name:    "Alex",
		Age:     35,
		Country: "Germany",
	}

	drink := Drink{
		Name:        "Cuba Libre",
		Ingredients: []string{"rum", "cola"},
	}

	// invalid parameters are possible by manually defining *RPCRequest
	rpcClient.CallBatch(RPCRequests{
		{
			Method: "singleRequest",
			Params: 3, // invalid, should be []int{3}
		},
	})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`[{"method":"singleRequest","params":3,"id":0,"jsonrpc":"2.0"}]`))

	// better use Params() unless you know what you are doing
	rpcClient.CallBatch(RPCRequests{
		{
			Method: "singleRequest",
			Params: Params(3), // always valid json rpc
		},
	})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`[{"method":"singleRequest","params":[3],"id":0,"jsonrpc":"2.0"}]`))

	// even better, use NewRequest()
	rpcClient.CallBatch(RPCRequests{
		NewRequest("multipleRequests1", 1),
		NewRequest("multipleRequests2", 2),
		NewRequest("multipleRequests3", 3),
	})
	gomega.Expect((<-requestChan).body).To(gomega.Equal(`[{"method":"multipleRequests1","params":[1],"id":0,"jsonrpc":"2.0"},{"method":"multipleRequests2","params":[2],"id":1,"jsonrpc":"2.0"},{"method":"multipleRequests3","params":[3],"id":2,"jsonrpc":"2.0"}]`))

	// test a huge batch request
	requests := RPCRequests{
		NewRequest("nullParam", nil),
		NewRequest("nullParams", nil, nil),
		NewRequest("emptyParams", []interface{}{}),
		NewRequest("emptyAnyParams", []string{}),
		NewRequest("emptyObject", struct{}{}),
		NewRequest("emptyObjectList", []struct{}{{}, {}}),
		NewRequest("boolParam", true),
		NewRequest("boolParams", true, false, true),
		NewRequest("stringParam", "Alex"),
		NewRequest("stringParams", "JSON", "RPC"),
		NewRequest("numberParam", 123),
		NewRequest("numberParams", 123, 321),
		NewRequest("floatParam", 1.23),
		NewRequest("floatParams", 1.23, 3.21),
		NewRequest("manyParams", "Alex", 35, true, nil, 2.34),
		NewRequest("emptyMissingPublicFieldObject", struct{ name string }{name: "Alex"}),
		NewRequest("singleStruct", person),
		NewRequest("singlePointerToStruct", &person),
		NewRequest("multipleStructs", person, &drink),
		NewRequest("singleStructInArray", []interface{}{person}),
		NewRequest("namedParameters", map[string]interface{}{
			"name": "Alex",
			"age":  35,
		}),
		NewRequest("anonymousStructNoTags", struct {
			Name string
			Age  int
		}{"Alex", 33}),
		NewRequest("anonymousStructWithTags", struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{"Alex", 33}),
		NewRequest("structWithNullField", struct {
			Name    string  `json:"name"`
			Address *string `json:"address"`
		}{"Alex", nil}),
	}
	rpcClient.CallBatch(requests)

	gomega.Expect((<-requestChan).body).To(gomega.Equal(`[{"method":"nullParam","params":[null],"id":0,"jsonrpc":"2.0"},` +
		`{"method":"nullParams","params":[null,null],"id":1,"jsonrpc":"2.0"},` +
		`{"method":"emptyParams","params":[],"id":2,"jsonrpc":"2.0"},` +
		`{"method":"emptyAnyParams","params":[],"id":3,"jsonrpc":"2.0"},` +
		`{"method":"emptyObject","params":{},"id":4,"jsonrpc":"2.0"},` +
		`{"method":"emptyObjectList","params":[{},{}],"id":5,"jsonrpc":"2.0"},` +
		`{"method":"boolParam","params":[true],"id":6,"jsonrpc":"2.0"},` +
		`{"method":"boolParams","params":[true,false,true],"id":7,"jsonrpc":"2.0"},` +
		`{"method":"stringParam","params":["Alex"],"id":8,"jsonrpc":"2.0"},` +
		`{"method":"stringParams","params":["JSON","RPC"],"id":9,"jsonrpc":"2.0"},` +
		`{"method":"numberParam","params":[123],"id":10,"jsonrpc":"2.0"},` +
		`{"method":"numberParams","params":[123,321],"id":11,"jsonrpc":"2.0"},` +
		`{"method":"floatParam","params":[1.23],"id":12,"jsonrpc":"2.0"},` +
		`{"method":"floatParams","params":[1.23,3.21],"id":13,"jsonrpc":"2.0"},` +
		`{"method":"manyParams","params":["Alex",35,true,null,2.34],"id":14,"jsonrpc":"2.0"},` +
		`{"method":"emptyMissingPublicFieldObject","params":{},"id":15,"jsonrpc":"2.0"},` +
		`{"method":"singleStruct","params":{"name":"Alex","age":35,"country":"Germany"},"id":16,"jsonrpc":"2.0"},` +
		`{"method":"singlePointerToStruct","params":{"name":"Alex","age":35,"country":"Germany"},"id":17,"jsonrpc":"2.0"},` +
		`{"method":"multipleStructs","params":[{"name":"Alex","age":35,"country":"Germany"},{"name":"Cuba Libre","ingredients":["rum","cola"]}],"id":18,"jsonrpc":"2.0"},` +
		`{"method":"singleStructInArray","params":[{"name":"Alex","age":35,"country":"Germany"}],"id":19,"jsonrpc":"2.0"},` +
		`{"method":"namedParameters","params":{"age":35,"name":"Alex"},"id":20,"jsonrpc":"2.0"},` +
		`{"method":"anonymousStructNoTags","params":{"Name":"Alex","Age":33},"id":21,"jsonrpc":"2.0"},` +
		`{"method":"anonymousStructWithTags","params":{"name":"Alex","age":33},"id":22,"jsonrpc":"2.0"},` +
		`{"method":"structWithNullField","params":{"name":"Alex","address":null},"id":23,"jsonrpc":"2.0"}]`))

	// create batch manually
	requests = []*RPCRequest{
		{
			Method:  "myMethod1",
			Params:  []int{1},
			ID:      123,   // will be forced to requests[i].ID == i unless you use CallBatchRaw
			JSONRPC: "7.0", // will be forced to "2.0"  unless you use CallBatchRaw
		},
		{
			Method:  "myMethod2",
			Params:  &person,
			ID:      321,     // will be forced to requests[i].ID == i unless you use CallBatchRaw
			JSONRPC: "wrong", // will be forced to "2.0" unless you use CallBatchRaw
		},
	}
	rpcClient.CallBatch(requests)

	gomega.Expect((<-requestChan).body).To(gomega.Equal(`[{"method":"myMethod1","params":[1],"id":0,"jsonrpc":"2.0"},` +
		`{"method":"myMethod2","params":{"name":"Alex","age":35,"country":"Germany"},"id":1,"jsonrpc":"2.0"}]`))

	// use raw batch
	requests = []*RPCRequest{
		{
			Method:  "myMethod1",
			Params:  []int{1},
			ID:      123,
			JSONRPC: "7.0",
		},
		{
			Method:  "myMethod2",
			Params:  &person,
			ID:      321,
			JSONRPC: "wrong",
		},
	}
	rpcClient.CallBatchRaw(requests)

	gomega.Expect((<-requestChan).body).To(gomega.Equal(`[{"method":"myMethod1","params":[1],"id":123,"jsonrpc":"7.0"},` +
		`{"method":"myMethod2","params":{"name":"Alex","age":35,"country":"Germany"},"id":321,"jsonrpc":"wrong"}]`))
}

// test if the result of an an rpc request is parsed correctly and if errors are thrown correctly
func TestRpcJsonResponseStruct(t *testing.T) {
	gomega.RegisterTestingT(t)
	rpcClient := NewClient(httpServer.URL)

	// empty return body is an error
	responseBody = ``
	res, err := rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(res).To(gomega.BeNil())

	// not a json body is an error
	responseBody = `{ "not": "a", "json": "object"`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(res).To(gomega.BeNil())

	// field "anotherField" not allowed in rpc response is an error
	responseBody = `{ "anotherField": "norpc"}`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(res).To(gomega.BeNil())

	// result null is ok
	responseBody = `{"result": null}`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Result).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())

	// error null is ok
	responseBody = `{"error": null}`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Result).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())

	// result and error null is ok
	responseBody = `{"result": null, "error": null}`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Result).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())

	// result string is ok
	responseBody = `{"result": "ok"}`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Result).To(gomega.Equal("ok"))

	// result with error null is ok
	responseBody = `{"result": "ok", "error": null}`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Result).To(gomega.Equal("ok"))

	// error with result null is ok
	responseBody = `{"error": {"code": 123, "message": "something wrong"}, "result": null}`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Result).To(gomega.BeNil())
	gomega.Expect(res.Error.Code).To(gomega.Equal(123))
	gomega.Expect(res.Error.Message).To(gomega.Equal("something wrong"))

	// error with code and message is ok
	responseBody = `{ "error": {"code": 123, "message": "something wrong"}}`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Result).To(gomega.BeNil())
	gomega.Expect(res.Error.Code).To(gomega.Equal(123))
	gomega.Expect(res.Error.Message).To(gomega.Equal("something wrong"))

	// check results

	// should return int correctly
	responseBody = `{ "result": 1 }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	i, err := res.GetInt()
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(i).To(gomega.Equal(int64(1)))

	// error on wrong type
	i = 3
	responseBody = `{ "result": "notAnInt" }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	i, err = res.GetInt()
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(i).To(gomega.Equal(int64(0)))

	// error on result null
	i = 3
	responseBody = `{ "result": null }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	i, err = res.GetInt()
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(i).To(gomega.Equal(int64(0)))

	b := false
	responseBody = `{ "result": true }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	b, err = res.GetBool()
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(b).To(gomega.Equal(true))

	b = true
	responseBody = `{ "result": 123 }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	b, err = res.GetBool()
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(b).To(gomega.Equal(false))

	var p *Person
	responseBody = `{ "result": {"name": "Alex", "age": 35, "anotherField": "something"} }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	err = res.GetObject(&p)
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(p.Name).To(gomega.Equal("Alex"))
	gomega.Expect(p.Age).To(gomega.Equal(35))
	gomega.Expect(p.Country).To(gomega.Equal(""))

	// TODO: How to check if result could be parsed or if it is default?
	p = nil
	responseBody = `{ "result": {"anotherField": "something"} }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	err = res.GetObject(&p)
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(p).NotTo(gomega.BeNil())

	// TODO: HERE######
	var pp *PointerFieldPerson
	responseBody = `{ "result": {"anotherField": "something", "country": "Germany"} }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	err = res.GetObject(&pp)
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(pp.Name).To(gomega.BeNil())
	gomega.Expect(pp.Age).To(gomega.BeNil())
	gomega.Expect(*pp.Country).To(gomega.Equal("Germany"))

	p = nil
	responseBody = `{ "result": null }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	err = res.GetObject(&p)
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(p).To(gomega.BeNil())

	// passing nil is an error
	p = nil
	responseBody = `{ "result": null }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	err = res.GetObject(p)
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(p).To(gomega.BeNil())

	p2 := &Person{
		Name: "Alex",
	}
	responseBody = `{ "result": null }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	err = res.GetObject(&p2)
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(p2).To(gomega.BeNil())

	p2 = &Person{
		Name: "Alex",
	}
	responseBody = `{ "result": {"age": 35} }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	err = res.GetObject(p2)
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(p2.Name).To(gomega.Equal("Alex"))
	gomega.Expect(p2.Age).To(gomega.Equal(35))

	// prefilled struct is kept on no result
	p3 := Person{
		Name: "Alex",
	}
	responseBody = `{ "result": null }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	err = res.GetObject(&p3)
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(p3.Name).To(gomega.Equal("Alex"))

	// prefilled struct is extended / overwritten
	p3 = Person{
		Name: "Alex",
		Age:  123,
	}
	responseBody = `{ "result": {"age": 35, "country": "Germany"} }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	err = res.GetObject(&p3)
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(p3.Name).To(gomega.Equal("Alex"))
	gomega.Expect(p3.Age).To(gomega.Equal(35))
	gomega.Expect(p3.Country).To(gomega.Equal("Germany"))

	// nil is an error
	responseBody = `{ "result": {"age": 35} }`
	res, err = rpcClient.Call("something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.Error).To(gomega.BeNil())
	err = res.GetObject(nil)
	gomega.Expect(err).NotTo(gomega.BeNil())
}

func TestRpcBatchJsonResponseStruct(t *testing.T) {
	gomega.RegisterTestingT(t)
	rpcClient := NewClient(httpServer.URL)

	// empty return body is an error
	responseBody = ``
	res, err := rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(res).To(gomega.BeNil())

	// not a json body is an error
	responseBody = `{ "not": "a", "json": "object"`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(res).To(gomega.BeNil())

	// field "anotherField" not allowed in rpc response is an error
	responseBody = `{ "anotherField": "norpc"}`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(res).To(gomega.BeNil())

	// result must be wrapped in array on batch request
	responseBody = `{"result": null}`
	_, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err.Error()).NotTo(gomega.BeNil())

	// result ok since in array
	responseBody = `[{"result": null}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(len(res)).To(gomega.Equal(1))
	gomega.Expect(res[0].Result).To(gomega.BeNil())

	// error null is ok
	responseBody = `[{"error": null}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res[0].Result).To(gomega.BeNil())
	gomega.Expect(res[0].Error).To(gomega.BeNil())

	// result and error null is ok
	responseBody = `[{"result": null, "error": null}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res[0].Result).To(gomega.BeNil())
	gomega.Expect(res[0].Error).To(gomega.BeNil())

	// result string is ok
	responseBody = `[{"result": "ok","id":0}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res[0].Result).To(gomega.Equal("ok"))
	gomega.Expect(res[0].ID).To(gomega.Equal(0))

	// result with error null is ok
	responseBody = `[{"result": "ok", "error": null}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res[0].Result).To(gomega.Equal("ok"))

	// error with result null is ok
	responseBody = `[{"error": {"code": 123, "message": "something wrong"}, "result": null}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res[0].Result).To(gomega.BeNil())
	gomega.Expect(res[0].Error.Code).To(gomega.Equal(123))
	gomega.Expect(res[0].Error.Message).To(gomega.Equal("something wrong"))

	// error with code and message is ok
	responseBody = `[{ "error": {"code": 123, "message": "something wrong"}}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res[0].Result).To(gomega.BeNil())
	gomega.Expect(res[0].Error.Code).To(gomega.Equal(123))
	gomega.Expect(res[0].Error.Message).To(gomega.Equal("something wrong"))

	// check results

	// should return int correctly
	responseBody = `[{ "result": 1 }]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res[0].Error).To(gomega.BeNil())
	i, err := res[0].GetInt()
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(i).To(gomega.Equal(int64(1)))

	// error on wrong type
	i = 3
	responseBody = `[{ "result": "notAnInt" }]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res[0].Error).To(gomega.BeNil())
	i, err = res[0].GetInt()
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(i).To(gomega.Equal(int64(0)))

	var p *Person
	responseBody = `[{"id":0, "result": {"name": "Alex", "age": 35}}, {"id":2, "result": {"name": "Lena", "age": 2}}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})

	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())

	gomega.Expect(res[0].Error).To(gomega.BeNil())
	gomega.Expect(res[0].ID).To(gomega.Equal(0))

	gomega.Expect(res[1].Error).To(gomega.BeNil())
	gomega.Expect(res[1].ID).To(gomega.Equal(2))

	_ = res[0].GetObject(&p)
	gomega.Expect(p.Name).To(gomega.Equal("Alex"))
	gomega.Expect(p.Age).To(gomega.Equal(35))

	_ = res[1].GetObject(&p)
	gomega.Expect(p.Name).To(gomega.Equal("Lena"))
	gomega.Expect(p.Age).To(gomega.Equal(2))

	// check if error occurred
	responseBody = `[{ "result": "someresult", "error": null}, { "result": null, "error": {"code": 123, "message": "something wrong"}}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.HasError()).To(gomega.BeTrue())

	// check if error occurred
	responseBody = `[{ "result": null, "error": {"code": 123, "message": "something wrong"}}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.HasError()).To(gomega.BeTrue())
	// check if error occurred
	responseBody = `[{ "result": null, "error": {"code": 123, "message": "something wrong"}}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.HasError()).To(gomega.BeTrue())

	// check if response mapping works
	responseBody = `[{ "id":123,"result": 123},{ "id":1,"result": 1}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.HasError()).To(gomega.BeFalse())
	resMap := res.AsMap()

	int1, _ := resMap[1].GetInt()
	int123, _ := resMap[123].GetInt()
	gomega.Expect(int1).To(gomega.Equal(int64(1)))
	gomega.Expect(int123).To(gomega.Equal(int64(123)))

	// check if getByID works
	int123, _ = res.GetByID(123).GetInt()
	gomega.Expect(int123).To(gomega.Equal(int64(123)))

	// check if error occurred
	responseBody = `[{ "result": null, "error": {"code": 123, "message": "something wrong"}}]`
	res, err = rpcClient.CallBatch(RPCRequests{
		NewRequest("something", 1, 2, 3),
	})
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(res.HasError()).To(gomega.BeTrue())
}

func TestRpcClient_CallFor(t *testing.T) {
	gomega.RegisterTestingT(t)
	rpcClient := NewClient(httpServer.URL)

	i := 0
	responseBody = `{"result":3,"id":0,"jsonrpc":"2.0"}`
	err := rpcClient.CallFor(&i, "something", 1, 2, 3)
	<-requestChan
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(i).To(gomega.Equal(3))
}

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Country string `json:"country"`
}

type PointerFieldPerson struct {
	Name    *string `json:"name"`
	Age     *int    `json:"age"`
	Country *string `json:"country"`
}

type Drink struct {
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
}

type Planet struct {
	Name       string     `json:"name"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	Distance int    `json:"distance"`
	Color    string `json:"color"`
}
