/*
Copyright 2022 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package jdownloader

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type DeviceInfo struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ResponseIdentifier struct {
	ResponseID int64 `json:"rid,omitempty"`
}

type DeviceList struct {
	ResponseIdentifier
	List []DeviceInfo `json:"list"`
}

type actionRequest struct {
	Url        string        `json:"url"`
	Params     []interface{} `json:"params,omitempty"`
	RequestId  int64         `json:"rid"`
	ApiVersion int           `json:"apiVer"`
}

type DataResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Source string      `json:"src,omitempty"`
	Type   string      `json:"type,omitempty"`
	ResponseIdentifier
}

type JdClient interface {
	// Connect connects to device and obtains session key
	Connect() error
	// IsConnected returns true if client is connected to API server
	IsConnected() bool
	// Reconnect reconnects client to API server
	Reconnect() error
	// Disconnect disconnects client from API server
	Disconnect() error
	// ListDevices lists all devices associated with account used to connect to API server
	ListDevices() (*[]DeviceInfo, error)
	// Device gets specific device instance based on device name
	Device(string) (Device, error)
	// ConfigHash return hash code of configuration.
	// This method can be used to determine if there was configuration change
	ConfigHash() string
}

func NewClient(email string, password string, logger *zap.SugaredLogger, opts ...ClientOption) JdClient {
	c := &jDownloaderClient{
		connected: false,
		email:     email,
		appKey:    appName,
		client: http.Client{
			Timeout: 15 * time.Second,
		},
		loginSecret:  createSecret(email, password, "server"),
		deviceSecret: createSecret(email, password, "device"),
		log:          logger,
		configHash:   hashConfigKeys(email, password),
		lastCall:     time.Now(),
		endpoint:     apiEndpoint,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
