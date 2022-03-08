// Copyright 2022 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/flamego/flamego"

	"github.com/flamego/recaptcha"
)

func main() {
	siteKey := flag.String("site-key", "", "The reCAPTCHA site key")
	secretKey := flag.String("secret-key", "", "The reCAPTCHA secret key")
	flag.Parse()

	f := flamego.Classic()
	f.Use(recaptcha.V3(
		recaptcha.Options{
			Secret:    *secretKey,
			VerifyURL: recaptcha.VerifyURLGoogle,
		},
	))

	f.Get("/", func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		_, _ = w.Write([]byte(fmt.Sprintf(`
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
    		data-sitekey="%s"
    		data-callback='onSubmit'
    		data-action='submit'>Submit</button>
	</form>
</body>
</html>
`, *siteKey)))
	})

	f.Post("/", func(w http.ResponseWriter, r *http.Request, re recaptcha.RecaptchaV3) {
		token := r.PostFormValue("g-recaptcha-response")
		fmt.Println("token", token)
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
