// This source file is part of the Gel open source project.
//
// Copyright Gel Data Inc. and the Gel authors.
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

package gel

import (
	"fmt"
	"maps"
	"strings"

	"github.com/geldata/gel-go/gelcfg"
	gel "github.com/geldata/gel-go/internal/client"
	gelerrint "github.com/geldata/gel-go/internal/gelerr"
)

// WithTxOptions returns a copy of c with the [gelcfg.TxOptions] set to opts.
func (c Client) WithTxOptions(opts gelcfg.TxOptions) *Client { //nolint:gocritic,lll
	if !opts.IsValid() {
		panic("TxOptions not created with NewTxOptions() are not valid")
	}

	c.copyPool()
	c.pool.QueryConfig.TxOptions = opts
	return &c
}

// WithRetryOptions returns a copy of c
// with the RetryOptions set to opts.
func (c Client) WithRetryOptions(opts gelcfg.RetryOptions) *Client { //nolint:gocritic,lll
	if !opts.IsValid() {
		panic("RetryOptions not created with NewRetryOptions() are not valid")
	}

	c.copyPool()
	c.pool.QueryConfig.RetryOptions = opts
	return &c
}

func (c *Client) copyPool() {
	annotations := make(map[string]string, len(c.pool.QueryConfig.Annotations))
	maps.Copy(annotations, c.pool.QueryConfig.Annotations)
	pool := *c.pool
	pool.QueryConfig.Annotations = annotations
	c.pool = &pool
}

// WithConfig returns a copy of c
// with configuration values set to cfg.
// This is equivalent to using the edgeql configure session command.
// For available configuration parameters refer to the [config documentation].
//
// [config documentation]: https://docs.geldata.com/reference/stdlib/cfg#ref-std-cfg
func (c Client) WithConfig(cfg map[string]interface{}) *Client { //nolint:gocritic,lll
	state := gel.CopyState(c.pool.State)

	var config map[string]interface{}
	if c, ok := state["config"]; ok {
		config = c.(map[string]interface{})
	} else {
		config = make(map[string]interface{}, len(cfg))
	}

	for k, v := range cfg {
		config[k] = v
	}

	state["config"] = config
	c.copyPool()
	c.pool.State = state
	return &c
}

// WithoutConfig returns a copy of c with keys unset from the configuration.
func (c Client) WithoutConfig(key ...string) *Client { // nolint:gocritic
	state := gel.CopyState(c.pool.State)

	if c, ok := state["config"]; ok {
		config := c.(map[string]interface{})
		for _, k := range key {
			delete(config, k)
		}
	}

	c.copyPool()
	c.pool.State = state
	return &c
}

// WithModuleAliases returns a copy of c with module name aliases set to
// aliases.
func (c Client) WithModuleAliases(aliases ...gelcfg.ModuleAlias) *Client { //nolint:gocritic,lll
	state := gel.CopyState(c.pool.State)

	var a []interface{}
	if b, ok := state["aliases"]; ok {
		a = b.([]interface{})
	}

	for i := 0; i < len(aliases); i++ {
		a = append(a, []interface{}{aliases[i].Alias, aliases[i].Module})
	}

	state["aliases"] = a
	c.copyPool()
	c.pool.State = state
	return &c
}

// WithoutModuleAliases returns a copy of c with aliases unset.
func (c Client) WithoutModuleAliases(aliases ...string) *Client { //nolint:gocritic,lll
	state := gel.CopyState(c.pool.State)

	if a, ok := state["aliases"]; ok {
		blacklist := make(map[string]struct{}, len(aliases))
		for _, name := range aliases {
			blacklist[name] = struct{}{}
		}

		var without []interface{}
		for _, p := range a.([]interface{}) {
			pair := p.([]interface{})
			key := pair[0].(string)
			if _, ok := blacklist[key]; !ok {
				without = append(without, []interface{}{key, pair[1]})
			}
		}

		state["aliases"] = without
	}

	c.copyPool()
	c.pool.State = state
	return &c
}

// WithGlobals returns a copy of c with its global variables updated from
// globals.
//
// WithGlobals does not remove variables that are not mentioned in globals.
// Instead use [Client.WithoutGlobals].
func (c Client) WithGlobals(globals map[string]interface{}) *Client { //nolint:gocritic,lll
	state := gel.CopyState(c.pool.State)

	var g map[string]interface{}
	if x, ok := state["globals"]; ok {
		g = x.(map[string]interface{})
	} else {
		g = make(map[string]interface{}, len(globals))
	}

	for k, v := range globals {
		g[k] = v
	}

	state["globals"] = g
	c.copyPool()
	c.pool.State = state
	return &c
}

// WithoutGlobals returns a copy of c with the specified global names unset.
func (c Client) WithoutGlobals(globals ...string) *Client { //nolint:gocritic
	state := gel.CopyState(c.pool.State)

	if c, ok := state["globals"]; ok {
		config := c.(map[string]interface{})
		for _, k := range globals {
			delete(config, k)
		}
	}

	c.copyPool()
	c.pool.State = state
	return &c
}

// WithWarningHandler returns a copy of c with its [gelcfg.WarningHandler] set
// to handler. If handler is nil, [gelcfg.LogWarnings] is used.
func (c Client) WithWarningHandler(handler gelcfg.WarningHandler) *Client { //nolint:gocritic,lll
	if handler == nil {
		handler = gelcfg.LogWarnings
	}

	c.copyPool()
	c.pool.QueryConfig.WarningHandler = handler
	return &c
}

// WithQueryOptions returns a copy of c with its gelcfg.Queryoptions set to
// opts.
func (c Client) WithQueryOptions(opts gelcfg.QueryOptions) *Client { //nolint:gocritic,lll
	c.copyPool()
	c.pool.QueryConfig.QueryOptions = opts
	return &c
}

// WithQueryTag returns a copy of c with the [sys::QueryStats] tag set.
//
// sys::QueryStats only records the tag from the first time a query is run.
// Running the query again with a different tag will not change the tag in the
// sys::QueryStats entry.
//
// [sys::QueryStats]: https://docs.geldata.com/reference/stdlib/sys#type::sys::QueryStats
func (c Client) WithQueryTag(tag string) (*Client, error) {
	for _, prefix := range []string{"gel/", "edgedb/"} {
		if strings.HasPrefix(tag, prefix) {
			return nil, gelerrint.NewInvalidArgumentError(
				fmt.Sprintf("reserved tag: %s*", prefix),
				nil,
			)
		}
	}

	if len(tag) > 128 {
		return nil, gelerrint.NewInvalidArgumentError(
			"tag too long (> 128 characters)",
			nil,
		)
	}

	c.copyPool()
	c.pool.QueryConfig.Annotations["tag"] = tag
	return &c, nil
}

// WithoutQueryTag returns a copy of c with the [sys::QueryStats] tag removed.
//
// [sys::QueryStats]: https://docs.geldata.com/reference/stdlib/sys#type::sys::QueryStats
func (c Client) WithoutQueryTag() *Client {
	c.copyPool()
	delete(c.pool.QueryConfig.Annotations, "tag")
	return &c
}
