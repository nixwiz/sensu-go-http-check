package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	URL          string
	SearchString string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-example-check",
			Short:    "Sensu Simple HTTP Check",
			Keyspace: "sensu.io/plugins/sensu-example-check/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "url",
			Env:       "CHECK_URL",
			Argument:  "url",
			Shorthand: "u",
			Default:   "http://localhost:80/",
			Usage:     "URL to test",
			Value:     &plugin.URL,
		},
		&sensu.PluginConfigOption{
			Path:      "search-string",
			Env:       "CHECK_SEARCH_STRING",
			Argument:  "search-string",
			Shorthand: "s",
			Default:   "",
			Usage:     "String to search for",
			Value:     &plugin.SearchString,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if len(plugin.URL) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--url or CHECK_URL environment variable is required")
	}
	if len(plugin.SearchString) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--search-string or CHECK_SEARCH_STRING environment variable is required")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {

	resp, err := http.Get(plugin.URL)
	if err != nil {
		return sensu.CheckStateCritical, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return sensu.CheckStateCritical, err
	}

	if strings.Contains(string(body), plugin.SearchString) {
		fmt.Printf("%s OK: %v, found %s at %s\n", plugin.PluginConfig.Name, resp.StatusCode, plugin.SearchString, plugin.URL)
		return sensu.CheckStateOK, nil
	} else {
		fmt.Printf("%s CRITICAL: %v, %s not found at %s\n", plugin.PluginConfig.Name, resp.StatusCode, plugin.SearchString, plugin.URL)
		return sensu.CheckStateCritical, err
	}

}
