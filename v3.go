// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package recaptcha

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// RecaptchaV3 is a reCAPTCHA V3 verify interface.
type RecaptchaV3 interface {
	// Verify verifies user's response and send result back to client.
	// It returns a score (1.0 is very likely a good interaction, 0.0 is very likely a bot).
	// Based on the score, you can take variable action in the context of your site.
	Verify(token string, remoteIP ...string) (*ResponseV3, error)
}

type recaptchaV3 struct {
	// secret is the shared key between your site and reCAPTCHA. [Required]
	secret string
	// response is the user response token provided by the reCAPTCHA client-side integration on your site. [Required]
	response string
	// remoteIP is the user's IP address. [Optional]
	remoteIP string
	// verifyURL is the reCAPTCHA backend service URL.
	verifyURL string
}

var _ RecaptchaV3 = (*recaptchaV3)(nil)

// ResponseV3 is the response struct which Google send back to the client.
type ResponseV3 struct {
	Success     bool      `json:"success"`      // whether this request was a valid reCAPTCHA token for your site
	Score       float64   `json:"score"`        // the score for this request (0.0 - 1.0)
	Action      string    `json:"action"`       // the action name for this request (important to verify)
	ChallengeTS time.Time `json:"challenge_ts"` // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	Hostname    string    `json:"hostname"`     // the hostname of the site where the reCAPTCHA was solved
	ErrorCodes  []string  `json:"error-codes"`  // optional
}

func (r *recaptchaV3) Verify(token string, remoteIP ...string) (*ResponseV3, error) {
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

	var response ResponseV3
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal reCAPTCHA response JSON: %v", string(resp))
	}

	return &response, nil
}
