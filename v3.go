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

// RecaptchaV3 is a reCAPTCHA V3 verify interface.
type RecaptchaV3 interface {
	// Verify verifies the given token. An optional remote IP of the user may be
	// passed as extra security criteria.
	Verify(token string, remoteIP ...string) (*ResponseV3, error)
}

var _ RecaptchaV3 = (*recaptchaV3)(nil)

type recaptchaV3 struct {
	client    *http.Client
	secret    string
	verifyURL string
}

// ResponseV3 is the response struct which Google send back to the client.
type ResponseV3 struct {
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
	// Score indicates the score of the request (0.0 - 1.0).
	Score float64 `json:"score"`
	// Action is the action name of the request.
	Action string `json:"action"`
}

func (r *recaptchaV3) Verify(token string, remoteIP ...string) (*ResponseV3, error) {
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

	var response ResponseV3
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal response body")
	}
	return &response, nil
}
