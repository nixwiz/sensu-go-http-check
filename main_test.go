package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
}

func TestExecuteCheck(t *testing.T) {

	testcases := []struct{
		status	int
		search	string
	}{
		{sensu.CheckStateOK, "SUCCESS"},
		{sensu.CheckStateCritical, "FAILURE"},
	}

	for _, tc := range testcases {
		event := corev2.FixtureEvent("entity1", "check")
		assert := assert.New(t)

		var test = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedMethod := "GET"
			expectedURI := "/"
			assert.Equal(expectedMethod, r.Method)
			assert.Equal(expectedURI, r.RequestURI)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("SUCCESS"))
		}))
		_, err := url.ParseRequestURI(test.URL)
		require.NoError(t, err)
		plugin.URL = test.URL
		plugin.SearchString = tc.search
		status, err := executeCheck(event)
		assert.NoError(err)
		assert.Equal(status, tc.status)
	}
}
