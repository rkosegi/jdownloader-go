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

import "github.com/pkg/errors"

type MockClient struct {
	devs      *[]DeviceInfo
	connected bool
	id        string
}

func (m *MockClient) SetDevices(devs *[]DeviceInfo) {
	m.devs = devs
}

func (m *MockClient) Connect() error {
	m.connected = true
	return nil
}

func (m *MockClient) IsConnected() bool {
	return m.connected
}

func (m *MockClient) Reconnect() error {
	return nil
}

func (m *MockClient) Disconnect() error {
	m.connected = false
	return nil
}

func (m *MockClient) ListDevices() (*[]DeviceInfo, error) {
	return m.devs, nil
}

func (m *MockClient) Device(name string) (Device, error) {
	for _, d := range *m.devs {
		if d.Name == name {
			return &MockDevice{
				id: name,
			}, nil
		}
	}
	return nil, errors.Errorf("no such device: %s", name)
}

func (m *MockClient) ConfigHash() string {
	return "123"
}

func NewMockClient() *MockClient {
	return &MockClient{
		devs: &[]DeviceInfo{},
	}
}

type MockDevice struct {
	Device
	links []DownloadLink
	id    string
}

func (d *MockDevice) LinkGrabber() LinkGrabber {
	//TODO implement me
	panic("implement me")
}

func (d *MockDevice) Downloader() Downloader {
	return &MockDownloader{
		dev: *d,
	}
}

func (d *MockDevice) Name() string {
	//TODO implement me
	panic("implement me")
}

func (d *MockDevice) Id() string {
	return d.id
}

func (d *MockDevice) Status() string {
	return "UNKNOWN"
}

func (d *MockDevice) ConnectionInfo() (*DirectConnectionInfo, error) {
	return &DirectConnectionInfo{}, nil
}

func (d *MockDevice) Packages(...LinkGrabberQueryPackagesOptions) (*[]DownloadPackage, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDevice) Links(...DownloadQueryLinksOptions) (*[]DownloadLink, error) {
	//TODO implement me
	panic("implement me")
}

func (d *MockDevice) SetLinks(links *[]DownloadLink) {
	d.links = *links
}

type MockDownloader struct {
	dev MockDevice
}

func (dw *MockDownloader) Remove([]int64, []int64) error {
	return nil
}

func (dw *MockDownloader) Packages(...LinkGrabberQueryPackagesOptions) (*[]DownloadPackage, error) {
	//TODO implement me
	panic("implement me")
}

func (dw *MockDownloader) Links(...DownloadQueryLinksOptions) (*[]DownloadLink, error) {
	return &dw.dev.links, nil
}

func (dw *MockDownloader) Start() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (dw *MockDownloader) Stop() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (dw *MockDownloader) Pause() (bool, error) {
	return false, nil
}

func (dw *MockDownloader) Speed() (*DownloadSpeedInfo, error) {
	return nil, nil
}

func (dw *MockDownloader) Force([]int64, []int64) error {
	return nil
}

func (dw *MockDownloader) State() (*DownloadState, error) {
	return nil, nil
}
