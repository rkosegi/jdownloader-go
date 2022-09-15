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
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
)

type AddLinksParams struct {
	Links             string  `json:"links"`
	Autostart         *bool   `json:"autostart,omitempty"`
	PackageName       *string `json:"packageName,omitempty"`
	DestinationFolder *string `json:"destinationFolder,omitempty"`
	DownloadPassword  *string `json:"downloadPassword"`
	ExtractPassword   *string `json:"extractPassword"`
}

type AddLinksOptions func(params *AddLinksParams)

func AddLinksOptionAutostart(autostart bool) AddLinksOptions {
	return func(params *AddLinksParams) {
		params.Autostart = &autostart
	}
}

func AddLinksOptionPackage(name string) AddLinksOptions {
	return func(params *AddLinksParams) {
		params.PackageName = &name
	}
}

func AddLinksOptionDestinationDir(name string) AddLinksOptions {
	return func(params *AddLinksParams) {
		params.DestinationFolder = &name
	}
}

func AddLinksOptionDownloadPassword(pass string) AddLinksOptions {
	return func(params *AddLinksParams) {
		params.DownloadPassword = &pass
	}
}

func AddLinksOptionExtractPassword(pass string) AddLinksOptions {
	return func(params *AddLinksParams) {
		params.ExtractPassword = &pass
	}
}

type QueryPackagesParams struct {
	AvailableOfflineCount     *bool     `json:"availableOfflineCount,omitempty"`
	AvailableOnlineCount      *bool     `json:"availableOnlineCount,omitempty"`
	AvailableTempUnknownCount *bool     `json:"availableTempUnknownCount,omitempty"`
	AvailableUnknownCount     *bool     `json:"availableUnknownCount,omitempty"`
	BytesTotal                *bool     `json:"bytesTotal,omitempty"`
	ChildCount                *bool     `json:"childCount,omitempty"`
	Comment                   *bool     `json:"comment,omitempty"`
	Enabled                   *bool     `json:"enabled,omitempty"`
	Hosts                     *bool     `json:"hosts,omitempty"`
	MaxResults                *int      `json:"maxResults,omitempty"`
	PackageUUIDs              *[]string `json:"packageUUIDs,omitempty"`
	Priority                  *bool     `json:"priority,omitempty"`
	SaveTo                    *bool     `json:"saveTo,omitempty"`
	StartAt                   *int      `json:"startAt,omitempty"`
	Status                    *bool     `json:"status,omitempty"`
}

type LinkGrabberQueryPackagesOptions func(params *QueryPackagesParams)

func LinkGrabberQueryPackagesOptionPackageUUIDs(uuids []string) LinkGrabberQueryPackagesOptions {
	return func(params *QueryPackagesParams) {
		params.PackageUUIDs = &uuids
	}
}

func QueryPackagesOptionDefault() LinkGrabberQueryPackagesOptions {
	return func(params *QueryPackagesParams) {
		params.AvailableOfflineCount = &yes
		params.AvailableOnlineCount = &yes
		params.AvailableTempUnknownCount = &yes
		params.AvailableUnknownCount = &yes
		params.BytesTotal = &yes
		params.ChildCount = &yes
		params.Comment = &yes
		params.Enabled = &yes
		params.Hosts = &yes
		params.Priority = &yes
		params.SaveTo = &yes
		params.Status = &yes
	}
}

type CrawledPackage struct {
	OnlineCount      *int      `json:"onlineCount,omitempty"`
	OfflineCount     *int      `json:"offlineCount,omitempty"`
	SaveTo           *string   `json:"saveTo,omitempty"`
	UnknownCount     *int      `json:"unknownCount,omitempty"`
	TempUnknownCount *int      `json:"tempUnknownCount,omitempty"`
	Uuid             *int64    `json:"uuid,omitempty"`
	BytesTotal       *uint64   `json:"bytesTotal,omitempty"`
	ChildCount       *int      `json:"childCount,omitempty"`
	Enabled          *bool     `json:"enabled,omitempty"`
	Hosts            *[]string `json:"hosts,omitempty"`
	Name             *string   `json:"name"`
}

type LinkGrabberQueryLinksParams struct {
	BytesTotal   *bool `json:"bytesTotal,omitempty"`
	Comment      *bool `json:"comment,omitempty"`
	Status       *bool `json:"status,omitempty"`
	Enabled      *bool `json:"enabled,omitempty"`
	MaxResults   *int  `json:"maxResults,omitempty"`
	StartAt      *int  `json:"startAt,omitempty"`
	Hosts        *bool `json:"hosts,omitempty"`
	Url          *bool `json:"url,omitempty"`
	Availability *bool `json:"availability,omitempty"`
	VariantIcon  *bool `json:"variantIcon,omitempty"`
	VariantName  *bool `json:"variantName,omitempty"`
	VariantID    *bool `json:"variantID,omitempty"`
	Variants     *bool `json:"variants,omitempty"`
	Priority     *bool `json:"priority,omitempty"`
}

type LinkGrabberQueryLinksOptions func(params *LinkGrabberQueryLinksParams)

func DefaultLinkGrabberQueryLinksOptions() LinkGrabberQueryLinksOptions {
	return func(params *LinkGrabberQueryLinksParams) {
		params.BytesTotal = &yes
		params.Comment = &yes
		params.Status = &yes
		params.Enabled = &yes
		params.Hosts = &yes
		params.Url = &yes
		params.Availability = &yes
		params.VariantIcon = &yes
		params.VariantID = &yes
		params.Variants = &yes
		params.Priority = &yes
	}
}

type CrawledLink struct {
	Availability     *string `json:"availability,omitempty	"`
	BytesTotal       *uint64 `json:"bytesTotal,omitempty"`
	Comment          *string `json:"comment,omitempty"`
	DownloadPassword *string `json:"downloadPassword,omitempty"`
	Enabled          *bool   `json:"enabled,omitempty"`
	Host             *string `json:"host,omitempty"`
	Name             *string `json:"name,omitempty"`
	PackageUuid      *int64  `json:"packageUUID,omitempty"`
	Priority         *string `json:"priority,omitempty"`
	Url              *string `json:"url,omitempty"`
	Uuid             *int64  `json:"uuid,omitempty"`
	Status           *string `json:"status,omitempty"`
	Variants         *bool   `json:"variants,omitempty"`
}

type LinkGrabber interface {
	// Clear clears list of links
	Clear() error
	// Packages gets list of packages
	Packages(options ...LinkGrabberQueryPackagesOptions) (*[]CrawledPackage, error)
	// Links queries links currently being present
	Links(...LinkGrabberQueryLinksOptions) (*[]CrawledLink, error)
	// Add adds one or more links into download queue
	Add([]string, ...AddLinksOptions) (*DataResponse, error)
	// IsCollecting checks if link grabber is collecting links
	IsCollecting() (bool, error)
	// Remove removes given linksIds and/or packageIds
	Remove([]int64, []int64) error
	// RenameLink renames link
	RenameLink(int64, string) error
}

type linkGrabber struct {
	log *zap.SugaredLogger
	d   *jDevice
}

func newLinkGrabber(log *zap.SugaredLogger, device *jDevice) LinkGrabber {
	return &linkGrabber{
		log: log.Named("links"),
		d:   device,
	}
}

func (l *linkGrabber) Links(options ...LinkGrabberQueryLinksOptions) (*[]CrawledLink, error) {
	params := &LinkGrabberQueryLinksParams{}
	if len(options) == 0 {
		defaults := DefaultLinkGrabberQueryLinksOptions()
		options = append(options, defaults)
	}
	for _, opt := range options {
		opt(params)
	}
	data, err := l.d.doDevice("/linkgrabberv2/queryLinks", true, params)
	if err != nil {
		return nil, err
	}
	items := make([]CrawledLink, 0)
	err = toObj(data, &items)
	if err != nil {
		return nil, err
	}
	return &items, nil
}

func (l *linkGrabber) Packages(options ...LinkGrabberQueryPackagesOptions) (*[]CrawledPackage, error) {
	return queryPackages("linkgrabberv2", l.d, options...)
}

func (l *linkGrabber) Add(links []string, options ...AddLinksOptions) (*DataResponse, error) {
	params := &AddLinksParams{
		Links: strings.Join(links, ","),
	}
	for _, opt := range options {
		opt(params)
	}
	data, err := l.d.doDevice("/linkgrabberv2/addLinks", true, params)
	if err != nil {
		return nil, err
	}
	resp := &DataResponse{}
	err = toObj(data, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *linkGrabber) IsCollecting() (bool, error) {
	data, err := l.d.doDevice("/linkgrabberv2/isCollecting", false, nil)
	if err != nil {
		return false, err
	}
	var res bool
	err = toObj(data, &res)
	return res, err
}

func (l *linkGrabber) Remove(linkIds []int64, packageIds []int64) error {
	if len(linkIds) == 0 && len(packageIds) == 0 {
		return errors.New("One of linkIds or packageIds must not be empty")
	}
	_, err := l.d.doDevice("/linkgrabberv2/removeLinks", false, linkIds, packageIds)
	return err
}

func (l *linkGrabber) Clear() error {
	_, err := l.d.doDevice("/linkgrabberv2/clearList", false)
	return err
}

func (l *linkGrabber) RenameLink(id int64, name string) error {
	_, err := l.d.doDevice("/linkgrabberv2/renameLink", false, id, name)
	return err
}

func queryPackages(prefix string, d *jDevice, options ...LinkGrabberQueryPackagesOptions) (*[]CrawledPackage, error) {
	params := &QueryPackagesParams{}
	if len(options) == 0 {
		options = append(options, QueryPackagesOptionDefault())
	}
	for _, opt := range options {
		opt(params)
	}
	data, err := d.doDevice(fmt.Sprintf("/%s/queryPackages", prefix), true, params)
	if err != nil {
		return nil, err
	}
	items := make([]CrawledPackage, 0)
	err = toObj(data, &items)
	if err != nil {
		return nil, err
	}
	return &items, nil
}
