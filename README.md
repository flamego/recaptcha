# recaptcha

Package recaptcha is a middleware that provides the verifying for reCAPTCHA.

## Installation

The minimum requirement of Go is **1.16**.

	go get github.com/flamego/recaptcha

## Getting started

```go
package main

import (
	"github.com/flamego/flamego"
	"github.com/flamego/recaptcha"
)

func main() {
	f := flamego.Classic()
	f.Use(recaptcha.V2(recaptcha.Options{
		Secret: "<YOUR_SECRET_HERE>",
		VerifyURL: Global,
	}))
	f.Get("/verify", func(c flamego.Context, r *recaptcha.recaptchaV2) {
		response, err := r.Verify(input)
		if response.Success{
			//... 
		}
	})
	f.Run()
}
```

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.