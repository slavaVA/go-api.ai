package gapiai_test

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
	. "github.com/slavaVA/go-api.ai"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"io"
	"net/http"
	"os"
)

var _ = Describe("Service", func() {
	var server *ghttp.Server

	BeforeEach(func() {
		server = ghttp.NewServer()
	})

	AfterEach(func() {
		//shut down the server between tests
		server.Close()
	})
	Describe("Text Query", func() {
		testToken := "123456789"
		var apiService *QueryService
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v1/query", "v=20150910"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"Bearer " + testToken},
						"Content-Type":  []string{"application/json; charset=utf-8"},
					}),
					ghttp.RespondWith(http.StatusOK, `{
						  "id": "5bb49696-549d-4655-bfb1-21e1dc806379",
						  "timestamp": "2016-12-30T14:29:02.746Z",
						  "result": {
						    "source": "agent",
						    "resolvedQuery": "Some query",
						    "action": "ActionName",
						    "actionIncomplete": false,
						    "parameters": {
						      "param1": "value1",
						      "param2": "value2"
						    },
						    "contexts": [],
						    "metadata": {
						      "intentId": "6500dd00-5f37-4fa0-a050-a8cf2428867b",
						      "webhookUsed": "false",
						      "webhookForSlotFillingUsed": "false",
						      "intentName": "WhatExchangeRates"
						    },
						    "fulfillment": {
						      "speech": "Some speech text",
						      "messages": [
							{
							  "type": 0,
							  "speech": "Message speech text"
							}
						      ]
						    },
						    "score": 0.69
						  },
						  "status": {
						    "code": 200,
						    "errorType": "success"
						  },
						  "sessionId": "111"
						}`),
				),
			)
			apiConfig := &ApiConfig{
				AccessToken: testToken,
				Lang:        English,
			}
			apiService = NewQueryAPIEndpoint(server.URL()+"/v1/", CurrentAPIVersion, apiConfig)
			apiService.EnableLogger(os.Stdout)
		})

		It("Should do simple Text Request", func() {
			sessionId := "111"
			requestText := "RequestText"
			response, err := apiService.TextRequest(sessionId, requestText)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(response).ShouldNot(BeNil())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.SessionID).Should(Equal(sessionId))
			Ω(response.Status.Code).Should(Equal(200))
		})
	})

	Describe("TTS", func() {
		testToken := "123456789"
		var apiService *TtsService
		speech := make([]byte, 100)
		BeforeEach(func() {
			for i, _ := range speech {
				speech[i] = byte(i)
			}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/tts", "text=Hello&v=20150910"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"Bearer " + testToken},
						"Accept-language": []string{"en"},
					}),
					ghttp.RespondWith(http.StatusOK, speech),
				),
			)
			apiConfig := &ApiConfig{
				AccessToken: testToken,
				Lang:        English,
			}
			apiService = NewTtsAPIEndpoint(server.URL()+"/v1/", CurrentAPIVersion, apiConfig)
			apiService.EnableLogger(os.Stdout)
		})

		It("Should do simple TTS Request", func() {
			var speechBuf []byte
			sh := func(r io.Reader)error {
				speechBuf = make([]byte, 100)
				n, err := r.Read(speechBuf)
				Ω(err).Should(Equal(io.EOF))
				Ω(n).Should(Equal(100))
				Ω(speechBuf).Should(Equal(speech))
				return nil
			}
			err := apiService.DoTts("Hello", sh)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
})
