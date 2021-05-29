// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package recaptcha

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// RecaptchaV2 is a reCAPTCHA V2 verify interface.
type RecaptchaV2 interface {
	// Verify verifies user's response and send result back to client.
	Verify(token string, remoteIP ...string) (*ResponseV2, error)
}

var _ RecaptchaV2 = (*recaptchaV2)(nil)

type recaptchaV2 struct {
	// secret is the shared key between your site and reCAPTCHA. [Required]
	secret string
	// response is the user response token provided by the reCAPTCHA client-side integration on your site. [Required]
	response string
	// remoteIP is the user's IP address. [Optional]
	remoteIP string
	// verifyURL is the reCAPTCHA backend service URL.
	verifyURL string
}

// ResponseV2 is the response struct which Google send back to the client.
type ResponseV2 struct {
	Success bool `json:"success"`
	// ChallengeTS is the timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ).
	ChallengeTS time.Time `json:"challenge_ts"`
	// Hostname is the hostname of the site where the reCAPTCHA was solved.
	Hostname string `json:"hostname"`
	// ErrorCodes returns the error codes when verify failed.
	ErrorCodes []string `json:"error-codes"`
}

func (r *recaptchaV2) Verify(token string, remoteIP ...string) (*ResponseV2, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}
	if len(remoteIP) > 0 {
		r.remoteIP = remoteIP[0]
	}

	resp, err := request(r.verifyURL, r.secret, r.response, r.remoteIP)
	if err != nil {
		return nil, errors.Wrap(err, "request reCAPTCHA server")
	}

	var response ResponseV2
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal reCAPTCHA response JSON: %v", string(resp))
	}

	return &response, nil
}
