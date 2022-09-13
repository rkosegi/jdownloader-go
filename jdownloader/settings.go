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

const (
	GeneralSettingsFile       = "org.jdownloader.settings.GeneralSettings.json"
	MyJdownloaderSettingsfile = "org.jdownloader.api.myjdownloader.MyJDownloaderSettings.json"
	LinkGrabberSettingsFile   = "org.jdownloader.gui.views.linkgrabber.addlinksdialog.LinkgrabberSettings.json"
	CustomProxyListFile       = "org.jdownloader.settings.InternetConnectionSettings.customproxylist.json"
)

type GeneralSettings struct {
	MaxSimultaneDownloadsPerHost   int    `json:"maxsimultanedownloadsperhost"`
	DefaultDownloadFolder          string `json:"defaultdownloadfolder"`
	OnSkipDueToAlreadyExistsAction string `json:"onskipduetoalreadyexistsaction"`
}

type LinkGrabberSettings struct {
	AutoExtractionEnabled           bool   `json:"autoextractionenabled"`
	LinkGrabberAutoStartEnabled     bool   `json:"linkgrabberautostartenabled"`
	VariousPackageEnabled           bool   `json:"variouspackageenabled"`
	LatestDownloadDestinationFolder string `json:"latestdownloaddestinationfolder"`
	LinkGrabberAddAtTop             bool   `json:"linkgrabberaddattop"`
	LinkGrabberAutoConfirmEnabled   bool   `json:"linkgrabberautoconfirmenabled"`
}

type MyJdownloaderSettings struct {
	UniqueDeviceIdSaltV2 string  `json:"uniquedeviceidsaltv2"`
	AutoConnectEnabledV2 bool    `json:"autoconnectenabledv2"`
	Password             string  `json:"password"`
	DebugEnabled         bool    `json:"debugenabled"`
	UniqueDeviceId       *string `json:"uniquedeviceid"`
	ServerHost           string  `json:"serverhost"`
	DirectConnectMode    string  `json:"directconnectmode"`
	UniqueDeviceIdV2     string  `json:"uniquedeviceidv2"`
	Email                string  `json:"email"`
	LastError            string  `json:"lasterror"`
	DeviceName           string  `json:"devicename"`
	LastLocalPort        int     `json:"lastlocalport"`
}

func DefaultGeneralSettings() *GeneralSettings {
	return &GeneralSettings{
		MaxSimultaneDownloadsPerHost:   1,
		DefaultDownloadFolder:          "download",
		OnSkipDueToAlreadyExistsAction: "SKIP_FILE",
	}
}

func DefaultMyJdownloaderSettings() *MyJdownloaderSettings {
	return &MyJdownloaderSettings{
		AutoConnectEnabledV2: true,
		DebugEnabled:         false,
		LastError:            "NONE",
		ServerHost:           "api.jdownloader.org",
		DirectConnectMode:    "LAN",
	}
}

func DefaultLinkGrabberSettings() *LinkGrabberSettings {
	return &LinkGrabberSettings{
		AutoExtractionEnabled:           true,
		LinkGrabberAutoStartEnabled:     true,
		VariousPackageEnabled:           true,
		LatestDownloadDestinationFolder: "download",
		LinkGrabberAddAtTop:             false,
		LinkGrabberAutoConfirmEnabled:   false,
	}
}

type ProxyServer struct {
	Type                       string  `json:"type"`
	Address                    *string `json:"address"`
	Port                       int     `json:"port"`
	PreferNativeImplementation bool    `json:"preferNativeImplementation"`
	ResolveHostName            bool    `json:"resolveHostName"`
	Username                   *string `json:"username"`
	ConnectMethodPrefered      bool    `json:"connectMethodPrefered"`
	Password                   *string `json:"password"`
}

type ProxyFilter struct {
	Type    string    `json:"type"`
	Entries *[]string `json:"entries"`
}

type ProxyServerEntry struct {
	Filter                 *ProxyFilter `json:"filter"`
	Enabled                bool         `json:"enabled"`
	Pac                    bool         `json:"pac"`
	ReconnectSupported     bool         `json:"reconnectSupported"`
	RangeRequestsSupported bool         `json:"rangeRequestsSupported"`
	Proxy                  *ProxyServer `json:"proxy"`
}

func DefaultProxyServerEntry() *ProxyServerEntry {
	return &ProxyServerEntry{
		Enabled:                false,
		Pac:                    false,
		ReconnectSupported:     false,
		RangeRequestsSupported: false,
		Proxy: &ProxyServer{
			Type:                       "NONE",
			Address:                    nil,
			Port:                       80,
			PreferNativeImplementation: false,
			ResolveHostName:            false,
			Username:                   nil,
			ConnectMethodPrefered:      false,
			Password:                   nil,
		},
	}
}

func DefaultProxyList() []*ProxyServerEntry {
	return []*ProxyServerEntry{
		DefaultProxyServerEntry(),
	}
}
