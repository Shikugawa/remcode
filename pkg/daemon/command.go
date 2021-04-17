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
	"encoding/json"
	"net/http"
	"strconv"
)

type Command struct {
	Token   string `json:"token"`
	Command string `json:"command"`
}

type CommandReceiverHandler struct {
	Pool *ConnPool
}

func (c *CommandReceiverHandler) notifyHandler(w http.ResponseWriter, r *http.Request) {
	var cmd Command

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := strconv.Atoi(cmd.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := c.Pool.Notify(token, cmd.Command); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
