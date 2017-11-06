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

package oauth2

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ory/fosite"
	"github.com/ory/fosite/internal"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestClientCredentials_HandleTokenEndpointRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockClientCredentialsGrantStorage(ctrl)
	chgen := internal.NewMockAccessTokenStrategy(ctrl)
	areq := internal.NewMockAccessRequester(ctrl)
	defer ctrl.Finish()

	h := ClientCredentialsGrantHandler{
		HandleHelper: &HandleHelper{
			AccessTokenStorage:  store,
			AccessTokenStrategy: chgen,
			AccessTokenLifespan: time.Hour,
		},
		ScopeStrategy: fosite.HierarchicScopeStrategy,
	}
	for k, c := range []struct {
		description string
		mock        func()
		req         *http.Request
		expectErr   error
	}{
		{
			description: "should fail because not responsible",
			expectErr:   fosite.ErrUnknownRequest,
			mock: func() {
				areq.EXPECT().GetGrantTypes().Return(fosite.Arguments{""})
			},
		},
		{
			description: "should pass",
			mock: func() {
				areq.EXPECT().GetSession().Return(new(fosite.DefaultSession))
				areq.EXPECT().GetGrantTypes().Return(fosite.Arguments{"client_credentials"})
				areq.EXPECT().GetRequestedScopes().Return([]string{"foo", "bar", "baz.bar"})
				areq.EXPECT().GetClient().Return(&fosite.DefaultClient{
					GrantTypes: fosite.Arguments{"client_credentials"},
					Scopes:     []string{"foo", "bar", "baz"},
				})
			},
		},
	} {
		c.mock()
		err := h.HandleTokenEndpointRequest(nil, areq)
		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		t.Logf("Passed test case %d", k)
	}
}

func TestClientCredentials_PopulateTokenEndpointResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := internal.NewMockClientCredentialsGrantStorage(ctrl)
	chgen := internal.NewMockAccessTokenStrategy(ctrl)
	areq := fosite.NewAccessRequest(new(fosite.DefaultSession))
	aresp := fosite.NewAccessResponse()
	defer ctrl.Finish()

	h := ClientCredentialsGrantHandler{
		HandleHelper: &HandleHelper{
			AccessTokenStorage:  store,
			AccessTokenStrategy: chgen,
			AccessTokenLifespan: time.Hour,
		},
		ScopeStrategy: fosite.HierarchicScopeStrategy,
	}
	for k, c := range []struct {
		description string
		mock        func()
		req         *http.Request
		expectErr   error
	}{
		{
			description: "should fail because not responsible",
			expectErr:   fosite.ErrUnknownRequest,
			mock: func() {
				areq.GrantTypes = fosite.Arguments{""}
			},
		},
		{
			description: "should fail because client not allowed",
			expectErr:   fosite.ErrInvalidGrant,
			mock: func() {
				areq.GrantTypes = fosite.Arguments{"client_credentials"}
				areq.Client = &fosite.DefaultClient{GrantTypes: fosite.Arguments{"foo"}}
			},
		},
		{
			description: "should pass",
			mock: func() {
				areq.GrantTypes = fosite.Arguments{"client_credentials"}
				areq.Session = &fosite.DefaultSession{}
				areq.Client = &fosite.DefaultClient{GrantTypes: fosite.Arguments{"client_credentials"}}
				chgen.EXPECT().GenerateAccessToken(nil, areq).Return("tokenfoo.bar", "bar", nil)
				store.EXPECT().CreateAccessTokenSession(nil, "bar", areq).Return(nil)
			},
		},
	} {
		c.mock()
		err := h.PopulateTokenEndpointResponse(nil, areq, aresp)
		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)
		t.Logf("Passed test case %d", k)
	}
}
