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
	driver "github.com/atomix/go-framework/pkg/atomix/driver/protocol/rsm"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)

	// Create an Atomix node
	node := driver.NewNode()

	// Register primitives on the Atomix node
	driver.RegisterCounterProxy(node)
	driver.RegisterElectionProxy(node)
	driver.RegisterIndexedMapProxy(node)
	driver.RegisterLockProxy(node)
	driver.RegisterLogProxy(node)
	driver.RegisterLeaderLatchProxy(node)
	driver.RegisterListProxy(node)
	driver.RegisterMapProxy(node)
	driver.RegisterSetProxy(node)
	driver.RegisterValueProxy(node)

	// Start the node
	if err := node.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Wait for an interrupt signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	// Stop the node after an interrupt
	if err := node.Stop(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
