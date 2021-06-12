// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package recaptcha

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/pkg/errors"
)

// VerifyURL is the url to verify user input.
type VerifyURL string

const (
	// VerifyURLGoogle is the default url to verify reCAPTCHA requests.
	VerifyURLGoogle VerifyURL = "https://www.google.com/recaptcha/api/siteverify"

	// VerifyURLGlobal is for the people who can't connect to the Google's server.
	VerifyURLGlobal VerifyURL = "https://www.recaptcha.net/recaptcha/api/siteverify"
)

type Version int

// Options contains options for the recaptcha.RecaptchaV2 and recaptcha.RecaptchaV3 middleware.
type Options struct {
	// Secret key is required. It's the shared key between your site and reCAPTCHA.
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
		var client RecaptchaV2 = &recaptchaV2{
			secret:    opts.Secret,
			verifyURL: string(opts.VerifyURL),
		}
		c.Map(client)
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
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "request %q", url)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}
	return body, nil
}
