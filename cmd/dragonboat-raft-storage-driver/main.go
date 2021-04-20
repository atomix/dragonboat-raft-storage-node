// Copyright 2019-present Open Networking Foundation.
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

package main

import (
	"fmt"
	"github.com/atomix/go-framework/pkg/atomix/cluster"
	"github.com/atomix/go-framework/pkg/atomix/driver"
	"github.com/atomix/go-framework/pkg/atomix/driver/proxy"
	"github.com/atomix/go-framework/pkg/atomix/driver/proxy/rsm"
	"github.com/atomix/go-framework/pkg/atomix/driver/proxy/rsm/counter"
	"github.com/atomix/go-framework/pkg/atomix/driver/proxy/rsm/election"
	"github.com/atomix/go-framework/pkg/atomix/driver/proxy/rsm/indexedmap"
	"github.com/atomix/go-framework/pkg/atomix/driver/proxy/rsm/leader"
	"github.com/atomix/go-framework/pkg/atomix/driver/proxy/rsm/list"
	"github.com/atomix/go-framework/pkg/atomix/driver/proxy/rsm/lock"
	"github.com/atomix/go-framework/pkg/atomix/driver/proxy/rsm/log"
	_map "github.com/atomix/go-framework/pkg/atomix/driver/proxy/rsm/map"
	"github.com/atomix/go-framework/pkg/atomix/driver/proxy/rsm/set"
	"github.com/atomix/go-framework/pkg/atomix/driver/proxy/rsm/value"
	"github.com/atomix/go-framework/pkg/atomix/logging"
	"os"
	"os/signal"
	"strconv"
)

const (
	driverNodeEnv      = "ATOMIX_DRIVER_NODE"
	driverNamespaceEnv = "ATOMIX_DRIVER_NAMESPACE"
	driverNameEnv      = "ATOMIX_DRIVER_NAME"
	driverPortEnv      = "ATOMIX_DRIVER_PORT"
)

func main() {
	logging.SetLevel(logging.DebugLevel)

	memberID := fmt.Sprintf("%s.%s", os.Getenv(driverNameEnv), os.Getenv(driverNamespaceEnv))
	nodeID := os.Getenv(driverNodeEnv)
	port, err := strconv.Atoi(os.Getenv(driverPortEnv))
	if err != nil {
		panic(err)
	}

	provider := func(c cluster.Cluster) proxy.Protocol {
		p := rsm.NewProtocol(c)
		counter.Register(p)
		election.Register(p)
		indexedmap.Register(p)
		lock.Register(p)
		log.Register(p)
		leader.Register(p)
		list.Register(p)
		_map.Register(p)
		set.Register(p)
		value.Register(p)
		return p
	}

	// Create a Raft driver
	d := driver.NewDriver(
		provider,
		driver.WithDriverID(memberID),
		driver.WithNodeID(nodeID),
		driver.WithPort(port))

	// Start the node
	if err := d.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Wait for an interrupt signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	// Stop the node after an interrupt
	if err := d.Stop(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
