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
	"bufio"
	"fmt"
	"log"
	"net"
)

type ConnPool struct {
	conn      map[int]*net.Conn
	nextToken int
}

func NewConnPool(connNumLimit int16) *ConnPool {
	return &ConnPool{
		conn:      make(map[int]*net.Conn, 0),
		nextToken: 0,
	}
}

func (c *ConnPool) Subscribe(conn *net.Conn) int {
	token := c.nextToken
	c.conn[c.nextToken] = conn
	c.nextToken++
	log.Printf("subscription from %d\n", token)
	return token
}

func (c *ConnPool) Unsubscribe(token int, msg []byte) error {
	if _, ok := c.conn[token]; !ok {
		return fmt.Errorf("failed to found conn %d", token)
	}

	log.Printf("unsubscribe from %d\n", token)
	cn := c.conn[token]

	(*cn).Write(msg) // fallthrough error
	(*cn).Close()
	delete(c.conn, token)

	return nil
}

func (c *ConnPool) UnsubscribeAll(msg []byte) {
	for token, _ := range c.conn {
		if err := c.Unsubscribe(token, msg); err != nil {
			continue
		}
	}
}

func (c *ConnPool) Notify(token int, msg string) error {
	if _, ok := c.conn[token]; !ok {
		return fmt.Errorf("failed to found conn %d", token)
	}

	_, err := (*c.conn[token]).Write([]byte(msg))

	if err != nil {
		return err
	}

	return nil
}

func (c *ConnPool) Start(token int) error {
	cn := c.conn[token]

	for {
		_, err := bufio.NewReader(*cn).ReadString('\n')

		if err != nil {
			if err := c.Unsubscribe(token, []byte(err.Error())); err != nil {
				return err
			}
			break
		}
	}

	if err := c.Unsubscribe(token, []byte("succeded to unsubscribe")); err != nil {
		return err
	}

	return nil
}
