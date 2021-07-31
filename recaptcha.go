// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package recaptcha

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/flamego/flamego"
)

// VerifyURL is the API URL to verify user input.
type VerifyURL string

const (
	// VerifyURLGoogle is the default API URL to verify reCAPTCHA requests.
	VerifyURLGoogle VerifyURL = "https://www.google.com/recaptcha/api/siteverify"

	// VerifyURLGlobal is API URL for the people who can't connect to the Google's server.
	VerifyURLGlobal VerifyURL = "https://www.recaptcha.net/recaptcha/api/siteverify"
)

// Options contains options for both recaptcha.RecaptchaV2 and recaptcha.RecaptchaV3 middleware.
type Options struct {
	// Secret is the shared key between your site and reCAPTCHA. This field is required.
	Secret string

	VerifyURL
}

// V2 returns a middleware handler that injects recaptcha.RecaptchaV2
// into the request context, which is used for verifying reCAPTCHA V2 requests.
func V2(opts Options) flamego.Handler {
	if opts.Secret == "" {
		panic("recaptcha: empty secret")
	}

	if opts.VerifyURL == "" {
		opts.VerifyURL = VerifyURLGoogle
	}

	return flamego.ContextInvoker(func(c flamego.Context) {
		client := &recaptchaV2{
			secret:    opts.Secret,
			verifyURL: string(opts.VerifyURL),
		}
		c.MapTo(client, (*RecaptchaV2)(nil))
	})
}

// V3 returns a middleware handler that injects recaptcha.RecaptchaV3
// into the request context, which is used for verifying reCAPTCHA V3 requests.
func V3(opts Options) flamego.Handler {
	if opts.Secret == "" {
		panic("recaptcha: empty secret")
	}

	if opts.VerifyURL == "" {
		opts.VerifyURL = VerifyURLGoogle
	}

	return flamego.ContextInvoker(func(c flamego.Context) {
		var client RecaptchaV3 = &recaptchaV3{
			secret:    opts.Secret,
			verifyURL: string(opts.VerifyURL),
		}

		c.Map(client)
	})
}

// request requests specific url and returns response.
func request(url, secret, response, remoteIP string) ([]byte, error) {
	url = fmt.Sprintf("%s?secret=%s&response=%s&remoteIP=%s", url, secret, response, remoteIP)
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "request %q", url)
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}
	return body, nil
}
