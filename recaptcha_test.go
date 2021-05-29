// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package recaptcha

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flamego/flamego"
	"github.com/stretchr/testify/assert"
)

func TestV2(t *testing.T) {
	f := flamego.NewWithLogger(&bytes.Buffer{})
	f.Use(V2(Options{
		Secret:    "6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe",
		VerifyURL: VerifyURLGlobal,
	}))
	f.Get("/", func(c flamego.Context, r *recaptchaV2) bool {
		response, err := c.Request().Body().String()
		assert.Nil(t, err)

		responseV2, err := r.Verify(response)
		assert.Nil(t, err)
		return responseV2.Success
	})

	resp := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)

	f.ServeHTTP(resp, req)
}
