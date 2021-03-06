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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	QueryService struct {
		ApiService
		queryURL string
	}
)

func NewQueryAPIEndpoint(url string, version string, cfg *ApiConfig) *QueryService {
	svc := &QueryService{
		ApiService: ApiService{
			logger: nil,
			Config: cfg,
		},
		queryURL: fmt.Sprint(url, "query?v=", version),
	}
	return svc
}

func DefaultQueryAPIEndpoint(cfg *ApiConfig) *QueryService {
	return NewQueryAPIEndpoint(apiAiURL, CurrentAPIVersion, cfg)
}

func (service *QueryService) TextRequest(sessionID string, text string) (*QueryResponse, error) {
	q := Query{
		Query:     []string{text},
		SessionID: sessionID,
	}
	return service.DoQuery(q)
}

func (service *QueryService) DoQuery(q Query) (*QueryResponse, error) {

	q.Lang = string(service.Config.Lang)

	jsonStr, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	service.debug("API AI request Body:", string(jsonStr))

	req, err := http.NewRequest("POST", service.queryURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+service.Config.AccessToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.ContentLength <= 0 {
		return nil, errors.New("Content length is 0")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	service.debug("API AI response Body:", string(body))

	queryResponse := &QueryResponse{}
	err = queryResponse.Decode(body)
	if err != nil {
		return nil, errors.New("Error parse body response:" + err.Error() + " Body:" + string(body))
	}

	if resp.StatusCode != http.StatusOK || queryResponse.Status.IsSuccess() == false {
		return nil, errors.New("Http Status " + resp.Status + " Body:" + string(body))
	}
	return queryResponse, nil
}
