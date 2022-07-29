package xoauth2

import (
	"github.com/emersion/go-sasl"
	"github.com/sqs/go-xoauth2"
)

// The XOAUTH2 mechanism name.
const XOAuth2 = "XOAUTH2"

type xOAuth2Client struct {
	user, accessToken string
}

func (c *xOAuth2Client) Start() (mech string, ir []byte, err error) {
	mech = XOAuth2
	ir = []byte(xoauth2.OAuth2String(c.user, c.accessToken))
	return
}

func (c *xOAuth2Client) Next(challenge []byte) (response []byte, err error) {
	return nil, nil
}

// A client implementation of the XOAuth2 authentication mechanism,
// as described in https://developers.google.com/gmail/imap/xoauth2-protocol.
func New(user, accessToken string) sasl.Client {
	return &xOAuth2Client{user, accessToken}
}
