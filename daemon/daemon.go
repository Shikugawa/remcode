// Copyright 2021 Rei Shimizu

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package daemon

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func CommandReceiver(port int16, pool *ConnPool) {
	crh := &CommandReceiverHandler{
		Pool: pool,
	}
	http.HandleFunc("/notify", crh.notifyHandler)

	log.Println("command receiver started")
	if err := http.ListenAndServe(":"+fmt.Sprint(port), nil); err != nil {
		log.Fatalln(err)
	}
}

func SubscriptionReceiver(port int16, pool *ConnPool) {
	l, err := net.Listen("tcp", ":"+fmt.Sprint(port))

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error()+"\n")
		return
	}

	defer l.Close()

	log.Println("subscription receiver started")

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("failed to accept")
			continue
		}

		token := pool.Subscribe(&c)
		go pool.Start(token)
	}
}

func Run() {
	pool := NewConnPool(1024)

	go CommandReceiver(4000, pool)
	go SubscriptionReceiver(3000, pool)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	pool.UnsubscribeAll([]byte("server finished"))
}
