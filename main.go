package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/swaggest/swgui"
	swg "github.com/swaggest/swgui/v4emb"
	"log"
	"net/http"
)

// @title        Swagger RPC POC
// @version      0.1.0
// @description  Swagger RPC POC Description
func main() {
	err := CheckInTest(false)
	if err != nil {
		return
	}

	r := chi.NewRouter()

	r.Post("/v1/rpc", rpcHandler)

	r.Get("/docs/swagger.json", staticHandler)

	r.Mount("/docs", swg.NewHandlerWithConfig(swgui.Config{
		Title:       "Swagger RPC POC",
		SwaggerJSON: "/docs/swagger.json",
		BasePath:    "/docs",
		SettingsUI:  SwguiSettings(nil, "/v1/rpc"),
	}))

	// Start server.
	log.Println("http://localhost:8080/docs")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./docs/swagger.json")
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	var req RpcRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResp(err, w)
		return
	}

	reqBody, err := json.Marshal(req.Body)
	if err != nil {
		errorResp(err, w)
		return
	}

	var resp string
	switch req.Method {
	case "Method1":
		resp, err = rpcHandlerMethod1(reqBody)
		if err != nil {
			errorResp(err, w)
			return
		}
	case "Method2":
		resp, err = rpcHandlerMethod2(reqBody)
		if err != nil {
			errorResp(err, w)
			return
		}
	default:
		errorResp(errors.New("unknown rpc method"), w)
		return
	}
	goodResp := GoodResponse{resp}
	goodRespJson, err := json.Marshal(goodResp)
	if err != nil {
		errorResp(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(goodRespJson)
}

// SwguiSettings adds JSON-RPC request interceptor to Swagger UI settings.
func SwguiSettings(settingsUI map[string]string, rpcPath string) map[string]string {
	if settingsUI == nil {
		settingsUI = make(map[string]string)
	}

	settingsUI["requestInterceptor"] = `function(request) {
				if (request.loadSpec) {
					return request;
				}
				var url = window.location.protocol + '//'+ window.location.host;
				var method = request.url.substring(url.length+1);
				request.url = url + '` + rpcPath + `';
				request.body = '{"jsonrpc": "2.0", "method": "' + method + '", "body": ' + request.body + '}';
				return request;
			}`

	return settingsUI
}

type RpcRequest struct {
	JsonRpc string `example:"2.0"`
	Method  string
	Body    interface{}
} // @name RpcRequest

type Method1 struct {
	Field1 string `example:"field1"`
	Field2 int    `example:"1"`
} // @name Method1

type Method2 struct {
	Field1 bool   `example:"true"`
	Field2 string `example:"field2"`
	Field3 uint   `example:"2"`
} // @name Method2

type GoodResponse struct {
	Message string
} // @name GoodResponse

type ErrorResponse struct {
	Error string
} // @name ErrorResponse

// Method1
// @Summary  Method1 handler.
// @Tags     POC
// @Accept   json
// @Produce  json
// @Param    req  body      Method1        true  "Request Method1"
// @Success  200  {object}  GoodResponse   "Response Method1"
// @Success  500  {object}  ErrorResponse  "Error"
// @Router   /Method1 [post]
func rpcHandlerMethod1(request []byte) (string, error) {
	var req Method1
	err := json.Unmarshal(request, &req)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(req), nil
}

// Method2
// @Summary  Method2 handler.
// @Tags     POC
// @Accept   json
// @Produce  json
// @Param    req  body      Method2       true  "Request Method2"
// @Success  200  {object}  GoodResponse  "Response Method2"
// @Router   /Method2 [post]
func rpcHandlerMethod2(request []byte) (string, error) {
	var req Method2
	err := json.Unmarshal(request, &req)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(req), nil
}

func errorResp(err error, w http.ResponseWriter) {
	// Add comment which should decrease coverage
	// Add comment which should decrease coverage
	fmt.Println("error: " + err.Error())
	errResp := ErrorResponse{err.Error()}
	errJson, _ := json.Marshal(errResp)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(errJson)
}

func CheckInTest(isErr bool) error {
	// just a comment added
	// Another comment
	if isErr {
		return errors.New("some error")
	}
	return nil
}
