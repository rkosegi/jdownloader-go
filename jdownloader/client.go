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
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	apiEndpoint       = "https://api.jdownloader.org"
	mediaType         = "application/aesjson-jd; charset=utf-8"
	appName           = "jdownloader-clientgo"
	paramSessionToken = "sessiontoken"
)

var (
	yes = true
)

type sessionInfo struct {
	SessionToken string `json:"sessiontoken"`
	RegainToken  string `json:"regaintoken"`
	ResponseIdentifier
}

type jDownloaderClient struct {
	connected             bool
	email                 string
	counter               int64
	appKey                string
	loginSecret           [32]byte
	deviceSecret          [32]byte
	serverEncryptionToken [32]byte
	deviceEncryptionToken [32]byte
	client                http.Client
	sessionToken          string
	regainToken           string
	lock                  sync.Mutex
	log                   *zap.SugaredLogger
	configHash            string
	lastCall              time.Time
	lastCallLock          sync.Mutex
	endpoint              string
	afterCallFn           func(error, time.Duration)
}

type ClientOption func(c *jDownloaderClient)

func ClientOptionApiEndpoint(url string) ClientOption {
	return func(c *jDownloaderClient) {
		c.endpoint = url
	}
}

func ClientOptionApiCallbacks(afterCallFn func(error, time.Duration)) ClientOption {
	return func(c *jDownloaderClient) {
		c.afterCallFn = afterCallFn
	}
}

func ClientOptionTimeout(timeout time.Duration) ClientOption {
	return func(c *jDownloaderClient) {
		c.client.Timeout = timeout
	}
}

func ClientOptionAppKey(app string) ClientOption {
	return func(c *jDownloaderClient) {
		c.appKey = app
	}
}

func (j *jDownloaderClient) ConfigHash() string {
	return j.configHash
}

func (j *jDownloaderClient) IsConnected() bool {
	return j.connected
}

func (j *jDownloaderClient) Device(name string) (_ Device, err error) {
	dl, err := j.ListDevices()
	if err != nil {
		return nil, err
	}
	var dev *DeviceInfo
	for _, d := range *dl {
		if d.Name == name {
			dev = &d
			break
		}
	}
	if dev == nil {
		return nil, errors.Errorf("no such device: %s", name)
	}
	return &jDevice{
		id:     dev.Id,
		name:   dev.Name,
		log:    j.log.Named(dev.Name),
		impl:   j,
		status: dev.Status,
	}, nil
}

func (j *jDownloaderClient) Connect() (err error) {
	j.serverEncryptionToken = j.loginSecret
	j.deviceEncryptionToken = j.deviceSecret
	j.connected = false

	data, err := j.doServer("/my/connect", http.MethodPost, []string{
		qp("email", strings.ToLower(j.email)),
		qp("appkey", strings.ToLower(j.appKey)),
	}, nil, j.loginSecret)
	if err != nil {
		return errors.Wrap(err, "API doServer failed")
	}
	session := &sessionInfo{}
	err = json.Unmarshal(data, session)
	if err != nil {
		return errors.Wrap(err, "invalid payload received from server")
	}
	if j.currentRid() != session.ResponseID {
		return errors.Errorf("mismatched RID, expected: %d, actual: %d", j.currentRid(), session.ResponseID)
	}
	newToken, err := hex.DecodeString(session.SessionToken)
	if err != nil {
		return errors.Wrap(err, "unable to decode session token from response")
	}
	j.sessionToken = session.SessionToken
	j.regainToken = session.RegainToken
	j.updateTokens(newToken)
	j.connected = true
	return nil
}

func (j *jDownloaderClient) ListDevices() (_ *[]DeviceInfo, err error) {
	err = j.reconnectIfNecessary()
	if err != nil {
		return nil, err
	}
	data, err := j.doServer("/my/listdevices", http.MethodGet, []string{
		qp(paramSessionToken, j.sessionToken),
	}, nil, j.serverEncryptionToken)
	if err != nil {
		return nil, err
	}
	dl := &DeviceList{}
	err = json.Unmarshal(data, dl)
	if err != nil {
		return nil, err
	}
	return &dl.List, nil
}

func (j *jDownloaderClient) Reconnect() (err error) {
	data, err := j.doServer("/my/reconnect", http.MethodGet, []string{
		qp(paramSessionToken, j.sessionToken),
		qp("regaintoken", j.regainToken),
	}, nil, j.serverEncryptionToken)
	if err != nil {
		return err
	}
	session := &sessionInfo{}
	err = json.Unmarshal(data, session)
	if err != nil {
		return errors.Wrap(err, "invalid payload received from server")
	}
	if j.currentRid() != session.ResponseID {
		return errors.Errorf("mismatched RID, expected: %d, actual: %d", j.currentRid(), session.ResponseID)
	}
	newToken, err := hex.DecodeString(session.SessionToken)
	if err != nil {
		return errors.Wrap(err, "unable to decode session token from response")
	}
	j.sessionToken = session.SessionToken
	j.regainToken = session.RegainToken
	j.updateTokens(newToken)
	j.connected = true
	return nil
}

func (j *jDownloaderClient) Disconnect() error {
	_, err := j.doServer("/my/disconnect", http.MethodPost, []string{
		qp(paramSessionToken, j.sessionToken),
	}, nil, j.serverEncryptionToken)
	return err
}

func (j *jDownloaderClient) updateTokens(newToken []byte) {
	j.serverEncryptionToken = updateToken(newToken, j.serverEncryptionToken)
	j.deviceEncryptionToken = updateToken(newToken, j.deviceEncryptionToken)
}

func (j *jDownloaderClient) reconnectIfNecessary() error {
	j.lastCallLock.Lock()
	defer j.lastCallLock.Unlock()
	if !j.connected || time.Now().Sub(j.lastCall).Seconds() > 30 {
		j.lastCall = time.Now()
		return j.Connect()
	}
	return nil
}

func (j *jDownloaderClient) do(path string, method string, data []byte, key [32]byte) (_ []byte, err error) {
	defer j.onApiDone(err, time.Now())
	j.lock.Lock()
	defer j.lock.Unlock()
	uri := fmt.Sprintf("%s%s", j.endpoint, path)
	var resp *http.Response
	j.log.Debugf("%s %s @ %d", method, uri, j.currentRid())
	if http.MethodGet == method {
		resp, err = j.client.Get(uri)
	} else {
		if data != nil {
			resp, err = j.client.Post(uri, mediaType, bytes.NewBuffer(data))
		} else {
			resp, err = j.client.Post(uri, mediaType, nil)
		}
	}
	if err != nil {
		return nil, err
	}
	j.log.Debugf("HTTP%d @ %d", resp.StatusCode, j.currentRid())
	defer bodycloser(resp.Body, j.log)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fully consume response body")
	}
	if resp.StatusCode != http.StatusOK {
		e := parseError(body, key, j.log)
		return nil, errors.Errorf("API doServer failed: %v", e)
	} else {
		return decode(body, key)
	}
}

func (j *jDownloaderClient) doServer(path string, method string, args []string, data []byte, key [32]byte) (_ []byte, err error) {
	if args == nil {
		args = make([]string, 0)
	}
	args = append(args, qp("rid", strconv.FormatInt(j.nextRid(), 10)))
	uri := strings.Join(args, "&")
	uri = fmt.Sprintf("%s?%s", path, uri)
	uri = fmt.Sprintf("%s&signature=%s", uri, sign(uri, key[:]))
	return j.do(uri, method, data, key)
}

func (j *jDownloaderClient) nextRid() int64 {
	return atomic.AddInt64(&j.counter, 1)
}

func (j *jDownloaderClient) currentRid() int64 {
	return j.counter
}

func (j *jDownloaderClient) onApiDone(err error, start time.Time) {
	if j.afterCallFn != nil {
		j.afterCallFn(err, time.Since(start))
	}
}
