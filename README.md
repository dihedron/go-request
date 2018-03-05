# go-request
This project implements a simple HTTP requests builder with a fluent API. The library can be imported via
``` golang
import (
	"github.com/dihedron/go-trequest"
)
```
and provides a simple factory pattern that can be used to add information and finlly retrieve an ```http.Request``` object.
Typical usage would be along the following lines:
``` golang
req, _ := New("").
	Base("https://www.example.com/").
	Path("api/v2/login?param1=value1").
	Add().
	QueryParametersFrom(testMapQP).
	QueryParametersFrom(testStructQP).
	HeadersFrom(&testMapH).
	HeadersFrom(&testStructH).
	Make()
```
The first instruction (```New("")```) creates a new factory.
