package restclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type (
	RestBuilder struct {
		baseUrl         string
		paths           string
		params          string
		client          *http.Client
		errorResponse   ResponseHandler
		responseHandler ResponseHandler
		requestHandler  ResponseHandler
		middleware      []Middleware
	}

	// ResponseHandler is used to validate or handle the response to a request.
	ResponseHandler = func(*http.Response) error

	// MiddlewareFunc defines a function to process middleware.
	Middleware func(next HandlerFunc) HandlerFunc

	// HandlerFunc defines a function to serve HTTP requests.
	HandlerFunc func(c context.Context) error

	// HTTPErrorHandler is a centralized HTTP error handler.
	HTTPErrorHandler func(err error, c context.Context)
)

func NewRestBuilder() *RestBuilder {
	return &RestBuilder{}
}

func (rb *RestBuilder) WithBaseUrl(url string) *RestBuilder {
	rb.baseUrl = url
	return rb
}

func (rb *RestBuilder) WithPathParams(format string, values ...any) *RestBuilder {
	rb.paths = fmt.Sprintf(format, values...)
	return rb
}

func (rb *RestBuilder) WithQueryParams(params map[string]any) *RestBuilder {
	var query strings.Builder
	for k, v := range params {

		switch j := reflect.ValueOf(v); j.Kind() {
		case reflect.String:
			fmt.Fprintf(&query, "%s=%s&", k, v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Fprintf(&query, "%s=%d&", k, v)
		default:
			log.Fatal("URL: Unkown format of query params. It has to be either ~Int or String")
		}
	}
	rb.params = strings.TrimSuffix(query.String(), "&")

	return rb
}

func (rb *RestBuilder) ResultJSON(v any) *RestBuilder {
	formatJson := func(res *http.Response) error {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, v); err != nil {
			return err
		}
		return nil
	}

	rb.responseHandler = formatJson
	return rb
}

func (rb *RestBuilder) ResultError(v any) *RestBuilder {
	formatError := func(res *http.Response) error {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			v = fmt.Errorf("%w", err)
		}
		res.Body.Close()
		fmt.Printf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
		v = fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
		return fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	rb.errorResponse = formatError
	return rb
}

func (rb *RestBuilder) Do(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rb.baseUrl+rb.paths+"?"+rb.params, nil)
	if err != nil {
		return err
	}

	if rb.client == nil {
		rb.client = http.DefaultClient
	}

	res, err := rb.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return rb.errorResponse(res)
	}

	if res.StatusCode == http.StatusOK {
		return rb.responseHandler(res)
	}

	return nil
}
