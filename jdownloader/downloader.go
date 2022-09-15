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
	"errors"
	"go.uber.org/zap"
)

type DownloadState struct {
	State *string `json:"state,omitempty"`
}

type DownloadSpeedInfo struct {
	Speed *float64 `json:"speed,omitempty"`
}

type DownloadQueryLinksParams struct {
	BytesTotal       *bool    `json:"bytesTotal,omitempty"`
	Comment          *bool    `json:"comment,omitempty"`
	Status           *bool    `json:"status,omitempty"`
	Enabled          *bool    `json:"enabled,omitempty"`
	MaxResults       *int     `json:"maxResults,omitempty"`
	StartAt          *int     `json:"startAt,omitempty"`
	Host             *bool    `json:"host,omitempty"`
	BytesLoaded      *bool    `json:"bytesLoaded,omitempty"`
	Speed            *bool    `json:"speed,omitempty"`
	Eta              *bool    `json:"eta,omitempty"`
	Finished         *bool    `json:"finished,omitempty"`
	FinishedDate     *bool    `json:"finishedDate,omitempty"`
	Running          *bool    `json:"running,omitempty"`
	Skipped          *bool    `json:"skipped,omitempty"`
	ExtractionStatus *bool    `json:"extractionStatus,omitempty"`
	PackageUUIDs     *[]int64 `json:"packageUUIDs"`
	Url              *bool    `json:"url"`
	Priority         *bool    `json:"priority"`
}

type DownloadQueryLinksOptions func(params *DownloadQueryLinksParams)

func DefaultDownloadQueryLinksOptions() DownloadQueryLinksOptions {
	return func(params *DownloadQueryLinksParams) {
		params.BytesLoaded = &yes
		params.BytesTotal = &yes
		params.Comment = &yes
		params.Status = &yes
		params.Enabled = &yes
		params.Host = &yes
		params.Speed = &yes
		params.Eta = &yes
		params.Finished = &yes
		params.FinishedDate = &yes
		params.Running = &yes
		params.Skipped = &yes
		params.ExtractionStatus = &yes
		params.Url = &yes
		params.Priority = &yes
	}
}

type DownloadLink struct {
	AddedDate        *int64   `json:"addedDate,omitempty"`
	BytesTotal       *int64   `json:"bytesTotal,omitempty"`
	BytesLoaded      *int64   `json:"bytesLoaded,omitempty"`
	Comment          *string  `json:"comment,omitempty"`
	DownloadPassword *string  `json:"downloadPassword"`
	Enabled          *bool    `json:"enabled,omitempty"`
	Eta              *int64   `json:"eta,omitempty"`
	ExtractionStatus *string  `json:"extractionStatus"`
	Finished         *bool    `json:"finished,omitempty"`
	FinishedDate     *int64   `json:"finishedDate,omitempty"`
	Host             *string  `json:"host,omitempty"`
	Name             *string  `json:"name,omitempty"`
	PackageUuid      *int64   `json:"packageUUID,omitempty"`
	Priority         *string  `json:"priority"`
	Skipped          *bool    `json:"skipped,omitempty"`
	Speed            *float64 `json:"speed,omitempty"`
	Status           *string  `json:"status,omitempty"`
	StatusIconKey    *string  `json:"statusIconKey,omitempty"`
	Url              *string  `json:"url,omitempty"`
	Uuid             *int64   `json:"uuid,omitempty"`
}

type DownloadPackage struct {
	ActiveTask       *string   `json:"activeTask"`
	BytesLoaded      *int64    `json:"bytesLoaded,omitempty"`
	BytesTotal       *int64    `json:"bytesTotal,omitempty"`
	ChildCount       *int      `json:"childCount,omitempty"`
	Comment          *string   `json:"comment"`
	DownloadPassword *string   `json:"downloadPassword"`
	Enabled          *bool     `json:"enabled,omitempty"`
	Eta              *int64    `json:"eta,omitempty"`
	Finished         *bool     `json:"finished,omitempty"`
	Hosts            *[]string `json:"hosts,omitempty"`
	Name             *string   `json:"name"`
	Priority         *string   `json:"priority"`
	Running          *bool     `json:"running"`
	SaveTo           *string   `json:"saveTo,omitempty"`
	Speed            *float64  `json:"speed,omitempty"`
	Status           *string   `json:"status"`
	StatusIconKey    *string   `json:"statusIconKey"`
	Uuid             *int64    `json:"uuid,omitempty"`
}

type Downloader interface {
	//Packages queries information about existing packages
	Packages(...LinkGrabberQueryPackagesOptions) (*[]DownloadPackage, error)
	//Links queries information about existing links
	Links(...DownloadQueryLinksOptions) (*[]DownloadLink, error)
	//Remove removes given links and/or packages
	Remove([]int64, []int64) error
	//Start starts download process
	Start() (bool, error)
	//Stop stops download process
	Stop() (bool, error)
	//Pause pauses download process
	Pause() (bool, error)
	//Speed get current download speed
	Speed() (*DownloadSpeedInfo, error)
	//Force forces download of given links/packages
	Force([]int64, []int64) error
	//State gets currect state of download process
	State() (*DownloadState, error)
}

type downloadController struct {
	l *zap.SugaredLogger
	d *jDevice
}

func newDownloadController(log *zap.SugaredLogger, d *jDevice) Downloader {
	return &downloadController{
		d: d,
		l: log.Named("download"),
	}
}

func (dc *downloadController) Links(options ...DownloadQueryLinksOptions) (*[]DownloadLink, error) {
	params := &DownloadQueryLinksParams{}
	if len(options) == 0 {
		defaults := DefaultDownloadQueryLinksOptions()
		options = append(options, defaults)
	}
	for _, opt := range options {
		opt(params)
	}
	data, err := dc.d.doDevice("/downloadsV2/queryLinks", true, params)
	if err != nil {
		return nil, err
	}
	items := make([]DownloadLink, 0)
	err = toObj(data, &items)
	if err != nil {
		return nil, err
	}
	return &items, nil
}

func (dc *downloadController) Packages(options ...LinkGrabberQueryPackagesOptions) (*[]DownloadPackage, error) {
	params := &QueryPackagesParams{}
	if len(options) == 0 {
		options = append(options, QueryPackagesOptionDefault())
	}
	for _, opt := range options {
		opt(params)
	}
	data, err := dc.d.doDevice("/downloadsV2/queryPackages", true, params)
	if err != nil {
		return nil, err
	}
	items := make([]DownloadPackage, 0)
	err = toObj(data, &items)
	if err != nil {
		return nil, err
	}
	return &items, nil
}

func (dc *downloadController) Start() (bool, error) {
	data, err := dc.d.doDevice("/downloadcontroller/start", false)
	if err != nil {
		return false, err
	}
	return data.Data.(bool), err
}

func (dc *downloadController) Stop() (bool, error) {
	data, err := dc.d.doDevice("/downloadcontroller/stop", false)
	if err != nil {
		return false, err
	}
	return data.Data.(bool), err
}

func (dc *downloadController) Pause() (bool, error) {
	data, err := dc.d.doDevice("/downloadcontroller/pause", false)
	if err != nil {
		return false, err
	}
	return data.Data.(bool), err
}

func (dc *downloadController) Speed() (*DownloadSpeedInfo, error) {
	data, err := dc.d.doDevice("/downloadcontroller/getSpeedInBps", false)
	if err != nil {
		return nil, err
	}
	speed := data.Data.(float64)
	return &DownloadSpeedInfo{Speed: &speed}, nil
}

func (dc *downloadController) Force(linkIds []int64, packageIds []int64) error {
	if len(linkIds) == 0 && len(packageIds) == 0 {
		return errors.New("one of linkIds or packageIds must not be empty")
	}
	_, err := dc.d.doDevice("/downloadcontroller/forceDownload", false, linkIds, packageIds)
	return err
}

func (dc *downloadController) Remove(linkIds []int64, packageIds []int64) error {
	if len(linkIds) == 0 && len(packageIds) == 0 {
		return errors.New("one of linkIds or packageIds must not be empty")
	}
	_, err := dc.d.doDevice("/downloadsV2/removeLinks", false, linkIds, packageIds)
	return err
}

func (dc *downloadController) State() (*DownloadState, error) {
	data, err := dc.d.doDevice("/downloadcontroller/getCurrentState", false)
	if err != nil {
		return nil, err
	}
	state := data.Data.(string)
	return &DownloadState{State: &state}, nil
}
