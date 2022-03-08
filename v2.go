// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package recaptcha

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// RecaptchaV2 is a reCAPTCHA V2 verify interface.
type RecaptchaV2 interface {
	// Verify verifies the given token. An optional remote IP of the user may be
	// passed as extra security criteria.
	Verify(token string, remoteIP ...string) (*ResponseV2, error)
}

var _ RecaptchaV2 = (*recaptchaV2)(nil)

type recaptchaV2 struct {
	client    *http.Client
	secret    string
	verifyURL string
}

// ResponseV2 is the response struct which Google send back to the client.
type ResponseV2 struct {
	// Success indicates whether the passcode valid, and does it meet security
	// criteria you specified.
	Success bool `json:"success"`
	// ChallengeTS is the timestamp of the challenge load (ISO format
	// yyyy-MM-dd'T'HH:mm:ssZZ).
	ChallengeTS time.Time `json:"challenge_ts"`
	// Hostname is the hostname of the site where the challenge was solved.
	Hostname string `json:"hostname"`
	// ErrorCodes contains the error codes when verify failed.
	ErrorCodes []string `json:"error-codes"`
}

func (r *recaptchaV2) Verify(token string, remoteIP ...string) (*ResponseV2, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}

	var ip string
	if len(remoteIP) > 0 {
		ip = remoteIP[0]
	}

	resp, err := request(r.client, r.verifyURL, r.secret, token, ip)
	if err != nil {
		return nil, errors.Wrap(err, "request reCAPTCHA server")
	}

	var response ResponseV2
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal response body")
	}
	return &response, nil
}
