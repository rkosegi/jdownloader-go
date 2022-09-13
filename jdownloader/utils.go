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
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/url"
)

func toObj(response *DataResponse, dst interface{}) error {
	if response == nil || response.Data == nil || dst == nil {
		return nil
	}
	data, err := json.Marshal(response.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func parseError(data []byte, key [32]byte, log *zap.SugaredLogger) map[string]interface{} {
	//1, try simple json unmarshal
	var v map[string]interface{}
	err := json.Unmarshal(data, &v)
	//2, if that doesn't work, it must be encrypted
	if err != nil {
		decoded, err := decode(data, key)
		if err != nil {
			log.Warnf("unable to decrypt error response: %v", err)
		} else {
			//3, now try to unmarshal error
			err = json.Unmarshal(decoded, &v)
			if err != nil {
				log.Warnf("error response is not a json: %v", err)
				return nil
			} else {
				return v
			}
		}
	} else {
		return v
	}
	return nil
}

//qp Creates escaped query parameter
func qp(key string, value string) string {
	return fmt.Sprintf("%s=%s", key, url.QueryEscape(value))
}

func bodycloser(b io.ReadCloser, log *zap.SugaredLogger) {
	err := b.Close()
	if err != nil {
		log.Warnf("error while closing reader: %v", err)
	}
}
