// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package recaptcha

import (
	"github.com/flamego/flamego"
	"github.com/imroc/req"
	"github.com/pkg/errors"
)

// VerifyURL is the url to verify user input.
type VerifyURL string

const (
	// Google is the default url to verify reCAPTCHA requests.
	VerifyURLGoogle VerifyURL = "https://www.google.com/recaptcha/api/siteverify"

	// Global is for the people who can't connect to the Google's server.
	VerifyURLGlobal VerifyURL = "https://www.recaptcha.net/recaptcha/api/siteverify"
)

type Version int

// Options contains options for the recaptcha.RecaptchaV2 middleware.
type Options struct {
	// Secret key is required. It's the shared key between your site and reCAPTCHA.
	Secret string

	VerifyURL
}

// V2 returns a middleware handler that injects recaptcha.RecaptchaV2
// into the request context, which is used for verifying reCAPTCHA V2 requests.
func V2(opts Options) flamego.Handler {
	if opts.Secret == "" {
		panic("Null secret input!")
	}

	if opts.VerifyURL == "" {
		opts.VerifyURL = Google
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
		panic("Null secret value input!")
	}

	if opts.VerifyURL == "" {
		opts.VerifyURL = Google
	}

	return flamego.ContextInvoker(func(c flamego.Context) {
		var client RecaptchaV3 = &recaptchaV3{
			secret:    opts.Secret,
			verifyURL: string(opts.VerifyURL),
		}

		c.Map(client)
	})
}

// request requests specific url and return response.
func request(url, secret, response, ip string) ([]byte, error) {
	params := req.Param{
		"secret":   secret,
		"response": response,
	}
	if ip != "" {
		params["remoteip"] = ip
	}

	resp, err := req.Get(url, params)
	if err != nil {
		return nil, errors.Wrapf(err, "request %q: ", url)
	}

	return resp.Bytes(), nil
}
