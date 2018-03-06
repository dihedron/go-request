# go-request
This project implements a simple HTTP requests builder with a fluent API. The library can be imported via
``` golang {.line-numbers}
import (
	"github.com/dihedron/go-request"
)
```
and provides a simple factory pattern that can be used to add information and finally retrieve an ```http.Request``` object. The ```Builder``` is available under the ```request``` namespace.
Typical usage would be along the following lines:
``` golang {.line-numbers}
req, _ := request.New("").                              // create the factory
	UserAgent("myUserAgent/1.0").                   // sets the user agent
	Base("https://www.example.com/").               // base URL for requests
	Path("api/v2/login?param1=value1").             // extra request path
	Add().                                          // start adding
	QueryParameter("param1", "value1a", "value1b").	// adds a query param to URL
	QueryParameter("param2", "value2").             // adds another query param
	Header("X-Auth-Token", "1234567890abcdef").     // adds a header
	WithJSONEntity(myTaggedStruct).                 // adds the request body from a struct
	Make()
```
The first instruction (```New("")```) creates a new factory.
``` golang {.line-numbers}
	QueryParametersFrom(testMapQP).
	QueryParametersFrom(testStructQP).
	HeadersFrom(&testMapH).
	HeadersFrom(&testStructH).
```