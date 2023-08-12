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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type DirectConnectionPort struct {
	Ip   *string `json:"ip"`
	Port *int    `json:"port"`
}

type DirectConnectionInfo struct {
	Mode                     *string                 `json:"mode"`
	Ports                    *[]DirectConnectionPort `json:"infos"`
	RebindProtectionDetected *bool                   `json:"rebindProtectionDetected"`
}

type Device interface {
	// LinkGrabber gets reference to LinkGrabber interface
	LinkGrabber() LinkGrabber
	// Downloader gets reference to Downloader interface
	Downloader() Downloader
	// Name gets this device's name
	Name() string
	// Id gets this device's ID
	Id() string
	// Status get this device's status
	Status() string
	// ConnectionInfo gets direct connection info
	ConnectionInfo() (*DirectConnectionInfo, error)
}

type jDevice struct {
	id     string
	name   string
	status string
	log    *zap.SugaredLogger
	impl   *jDownloaderClient
}

func (d *jDevice) LinkGrabber() LinkGrabber {
	return newLinkGrabber(d.log, d)
}

func (d *jDevice) Downloader() Downloader {
	return newDownloadController(d.log, d)
}

func (d *jDevice) Name() string {
	return d.name
}

func (d *jDevice) Status() string {
	return d.status
}

func (d *jDevice) Id() string {
	return d.id
}

func (d *jDevice) ConnectionInfo() (*DirectConnectionInfo, error) {
	data, err := d.doDevice("/device/getDirectConnectionInfos", false)
	if err != nil {
		return nil, err
	}
	info := &DirectConnectionInfo{}
	err = toObj(data, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func serializeParams(marshal bool, params ...interface{}) ([]interface{}, error) {
	if len(params) == 1 && params[0] == nil {
		return nil, nil
	}
	ps := make([]interface{}, 0)
	for _, p := range params {
		if marshal {
			s, err := json.Marshal(p)
			if err != nil {
				return nil, errors.WithMessagef(err, "unable to marshal parameter into json string: %v", p)
			}
			ps = append(ps, string(s))
		} else {
			ps = append(ps, p)
		}
	}
	return ps, nil
}

func (d *jDevice) doDevice(action string, marshal bool, params ...interface{}) (_ *DataResponse, err error) {
	err = d.impl.reconnectIfNecessary()
	if err != nil {
		return nil, err
	}
	qs := fmt.Sprintf("t_%s_%s%s", url.QueryEscape(d.impl.sessionToken), url.QueryEscape(d.id), action)
	p, err := serializeParams(marshal, params...)
	if err != nil {
		return nil, err
	}
	data := &actionRequest{
		Url:        action,
		Params:     p,
		RequestId:  d.impl.nextRid(),
		ApiVersion: 1,
	}
	plaintext, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	ciphertext, err := encrypt(plaintext, d.impl.deviceEncryptionToken)
	if err != nil {
		return nil, err
	}
	payload := base64.StdEncoding.EncodeToString(ciphertext)
	body, err := d.impl.do(fmt.Sprintf("/%s", qs), http.MethodPost, []byte(payload), d.impl.deviceEncryptionToken)
	if err != nil {
		return nil, err
	}
	result := &DataResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

var _ Device = &jDevice{}
