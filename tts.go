package gapiai

/***********************************************************************************************************************
 *
 * Go client-side library for API.AI
 * =================================================
 *
 * Copyright (C) 2017 by Slava Vasylyev
 *
 *
 * *********************************************************************************************************************
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 ***********************************************************************************************************************/

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type (
	TtsService struct {
		ApiService
		url string
	}

	SpeechHandler func(io.Reader, int64)
)

func NewTtsAPIEndpoint(url string, version string, cfg *ApiConfig) *TtsService {
	svc := &TtsService{
		ApiService: ApiService{
			logger: nil,
			Config: cfg,
		},
		url: fmt.Sprint(url, "tts?v=", version),
	}
	return svc
}

func DefaultTtsAPIEndpoint(cfg *ApiConfig) *TtsService {
	return NewTtsAPIEndpoint(apiAiURL, CurrentAPIVersion, cfg)
}

func (service *TtsService) DoTts(text string, handler SpeechHandler) error {

	req, err := http.NewRequest("GET", service.url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+service.Config.AccessToken)
	req.Header.Set("Accept-Language", string(service.Config.Lang))

	// get params
	query := req.URL.Query()
	query.Add("text", text)

	req.URL.RawQuery = query.Encode()
	service.debug("Raw query", req.URL.RawQuery)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.ContentLength <= 0 {
		return errors.New("Content length is 0")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Http Status " + resp.Status)
	}

	handler(resp.Body, resp.ContentLength)
	return nil
}
