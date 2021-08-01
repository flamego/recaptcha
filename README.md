# recaptcha

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/flamego/recaptcha/Go?logo=github&style=for-the-badge)](https://github.com/flamego/recaptcha/actions?query=workflow%3AGo)
[![Codecov](https://img.shields.io/codecov/c/gh/flamego/recaptcha?logo=codecov&style=for-the-badge)](https://app.codecov.io/gh/flamego/recaptcha)
[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/flamego/recaptcha?tab=doc)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/flamego/recaptcha)

Package recaptcha is a middleware that provides reCAPTCHA integration for [Flamego](https://github.com/flamego/flamego).

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
		VerifyURL: recaptcha.VerifyURLGlobal,
	}))
	f.Get("/verify", func(c flamego.Context, r recaptcha.RecaptchaV2) {
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
