// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fosite

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// WriteIntrospectionError responds with token metadata discovered by token introspection as defined in
// https://tools.ietf.org/search/rfc7662#section-2.2
//
// If the protected resource uses OAuth 2.0 client credentials to
// authenticate to the introspection endpoint and its credentials are
// invalid, the authorization server responds with an HTTP 401
// (Unauthorized) as described in Section 5.2 of OAuth 2.0 [RFC6749].
//
// If the protected resource uses an OAuth 2.0 bearer token to authorize
// its call to the introspection endpoint and the token used for
// authorization does not contain sufficient privileges or is otherwise
// invalid for this request, the authorization server responds with an
// HTTP 401 code as described in Section 3 of OAuth 2.0 Bearer Token
// Usage [RFC6750].
//
// Note that a properly formed and authorized query for an inactive or
// otherwise invalid token (or a token the protected resource is not
// allowed to know about) is not considered an error response by this
// specification.  In these cases, the authorization server MUST instead
// respond with an introspection response with the "active" field set to
// "false" as described in Section 2.2.
func (f *Fosite) WriteIntrospectionError(rw http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch errors.Cause(err) {
	case ErrInvalidRequest, ErrRequestUnauthorized:
		writeJsonError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	_ = json.NewEncoder(rw).Encode(struct {
		Active bool `json:"active"`
	}{Active: false})
}

// WriteIntrospectionResponse responds with an error if token introspection failed as defined in
// https://tools.ietf.org/search/rfc7662#section-2.3
//
// The server responds with a JSON object [RFC7159] in "application/
// json" format with the following top-level members.
//
// * active
// REQUIRED.  Boolean indicator of whether or not the presented token
// is currently active.  The specifics of a token's "active" state
// will vary depending on the implementation of the authorization
// server and the information it keeps about its tokens, but a "true"
// value return for the "active" property will generally indicate
// that a given token has been issued by this authorization server,
// has not been revoked by the resource owner, and is within its
// given time window of validity (e.g., after its issuance time and
// before its expiration time).  See Section 4 for information on
// implementation of such checks.
//
// * scope
// OPTIONAL.  A JSON string containing a space-separated list of
// scopes associated with this token, in the format described in
// Section 3.3 of OAuth 2.0 [RFC6749].
//
// * client_id
// OPTIONAL.  Client identifier for the OAuth 2.0 client that
// requested this token.
//
// * username
// OPTIONAL.  Human-readable identifier for the resource owner who
// authorized this token.
//
// * token_type
// OPTIONAL.  Type of the token as defined in Section 5.1 of OAuth
// 2.0 [RFC6749].
//
// * exp
// OPTIONAL.  Integer timestamp, measured in the number of seconds
// since January 1 1970 UTC, indicating when this token will expire,
// as defined in JWT [RFC7519].
//
// * iat
// OPTIONAL.  Integer timestamp, measured in the number of seconds
// since January 1 1970 UTC, indicating when this token was
// originally issued, as defined in JWT [RFC7519].
//
// * nbf
// OPTIONAL.  Integer timestamp, measured in the number of seconds
// since January 1 1970 UTC, indicating when this token is not to be
// used before, as defined in JWT [RFC7519].
//
// * sub
// OPTIONAL.  Subject of the token, as defined in JWT [RFC7519].
// Usually a machine-readable identifier of the resource owner who
// authorized this token.
//
// * aud
// OPTIONAL.  Service-specific string identifier or list of string
// identifiers representing the intended audience for this token, as
// defined in JWT [RFC7519].
//
// * iss
// OPTIONAL.  String representing the issuer of this token, as
// defined in JWT [RFC7519].
//
// * jti
// OPTIONAL.  String identifier for the token, as defined in JWT
// [RFC7519].
//
// Specific implementations MAY extend this structure with their own
// service-specific response names as top-level members of this JSON
// object.  Response names intended to be used across domains MUST be
// registered in the "OAuth Token Introspection Response" registry
// defined in Section 3.1.
//
// The authorization server MAY respond differently to different
// protected resources making the same request.  For instance, an
// authorization server MAY limit which scopes from a given token are
// returned for each protected resource to prevent a protected resource
// from learning more about the larger network than is necessary for its
// operation.
//
// The response MAY be cached by the protected resource to improve
// performance and reduce load on the introspection endpoint, but at the
// cost of liveness of the information used by the protected resource to
// make authorization decisions.  See Section 4 for more information
// regarding the trade off when the response is cached.
//
//
// For example, the following response contains a set of information
// about an active token:
//
// The following is a non-normative example response:
//
//	 HTTP/1.1 200 OK
//	 Content-Type: application/json
//
//	 {
//	   "active": true,
//	   "client_id": "l238j323ds-23ij4",
//	   "username": "jdoe",
//	   "scope": "read write dolphin",
//	   "sub": "Z5O3upPC88QrAjx00dis",
//	   "aud": "https://protected.example.net/resource",
//	   "iss": "https://server.example.com/",
//	   "exp": 1419356238,
//	   "iat": 1419350238,
//	   "extension_field": "twenty-seven"
//	 }
//
// If the introspection call is properly authorized but the token is not
// active, does not exist on this server, or the protected resource is
// not allowed to introspect this particular token, then the
// authorization server MUST return an introspection response with the
// "active" field set to "false".  Note that to avoid disclosing too
// much of the authorization server's state to a third party, the
// authorization server SHOULD NOT include any additional information
// about an inactive token, including why the token is inactive.
//
// The following is a non-normative example response for a token that
// has been revoked or is otherwise invalid:
//
//	 HTTP/1.1 200 OK
//	 Content-Type: application/json
//
//	 {
//	   "active": false
//	 }
func (f *Fosite) WriteIntrospectionResponse(rw http.ResponseWriter, r IntrospectionResponder) {
	if !r.IsActive() {
		_ = json.NewEncoder(rw).Encode(&struct {
			Active bool `json:"active"`
		}{Active: false})
		return
	}

	_ = json.NewEncoder(rw).Encode(struct {
		Active    bool    `json:"active"`
		ClientID  string  `json:"client_id,omitempty"`
		Scope     string  `json:"scope,omitempty"`
		ExpiresAt int64   `json:"exp,omitempty"`
		IssuedAt  int64   `json:"iat,omitempty"`
		Subject   string  `json:"sub,omitempty"`
		Username  string  `json:"username,omitempty"`
		Session   Session `json:"sess,omitempty"`
	}{
		Active:    true,
		ClientID:  r.GetAccessRequester().GetClient().GetID(),
		Scope:     strings.Join(r.GetAccessRequester().GetGrantedScopes(), " "),
		ExpiresAt: r.GetAccessRequester().GetSession().GetExpiresAt(AccessToken).Unix(),
		IssuedAt:  r.GetAccessRequester().GetRequestedAt().Unix(),
		Subject:   r.GetAccessRequester().GetSession().GetSubject(),
		Username:  r.GetAccessRequester().GetSession().GetUsername(),
		// Session:   r.GetAccessRequester().GetSession(),
	})
}
