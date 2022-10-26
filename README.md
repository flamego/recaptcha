# recaptcha

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/flamego/recaptcha/Go?logo=github&style=for-the-badge)](https://github.com/flamego/recaptcha/actions?query=workflow%3AGo)
[![Codecov](https://img.shields.io/codecov/c/gh/flamego/recaptcha?logo=codecov&style=for-the-badge)](https://app.codecov.io/gh/flamego/recaptcha)
[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/flamego/recaptcha?tab=doc)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/flamego/recaptcha)

Package recaptcha is a middleware that provides reCAPTCHA integration for [Flamego](https://github.com/flamego/flamego).

## Installation

The minimum requirement of Go is **1.18**.

	go get github.com/flamego/recaptcha

## Getting started

```html
<!-- templates/home.tmpl -->
<html>
<head>
  <script src="https://www.google.com/recaptcha/api.js"></script>
</head>
<body>
<script>
  function onSubmit(token) {
    document.getElementById("demo-form").submit();
  }
</script>
<form id="demo-form" method="POST">
  <button class="g-recaptcha"
    data-sitekey="{{.SiteKey}}"
    data-callback='onSubmit'
    data-action='submit'>Submit</button>
</form>
</body>
</html>
```

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/flamego/recaptcha"
	"github.com/flamego/template"
)

func main() {
	f := flamego.Classic()
	f.Use(template.Templater())
	f.Use(recaptcha.V3(
		recaptcha.Options{
			Secret:    "<SECRET>",
			VerifyURL: recaptcha.VerifyURLGoogle,
		},
	))
	f.Get("/", func(t template.Template, data template.Data) {
		data["SiteKey"] = "<SITE KEY>"
		t.HTML(http.StatusOK, "home")
	})
	f.Post("/", func(w http.ResponseWriter, r *http.Request, re recaptcha.RecaptchaV3) {
		token := r.PostFormValue("g-recaptcha-response")
		resp, err := re.Verify(token)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		} else if !resp.Success {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("Verification failed, error codes %v", resp.ErrorCodes)))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Verified!"))
	})
	f.Run()
}
```

## Getting help

- Read [documentation and examples](https://flamego.dev/middleware/recaptcha.html).
- Please [file an issue](https://github.com/flamego/flamego/issues) or [start a discussion](https://github.com/flamego/flamego/discussions) on the [flamego/flamego](https://github.com/flamego/flamego) repository.

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.
