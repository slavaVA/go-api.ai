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
)

var _ = Describe("Model", func() {

	It("Should decode query response", func() {
		resultStr := `{
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
      "intentName": "TestIntentName"
    },
    "fulfillment": {
      "speech": "Some speech text",
      "messages": [
        {
          "type": 0,
          "speech": "Message speesh text"
        }
      ]
    },
    "score": 0.69
  },
  "status": {
    "code": 200,
    "errorType": "success"
  },
  "sessionId": "1"
}
`
		response := &QueryResponse{}
		err := response.Decode([]byte(resultStr))
		Ω(err).ShouldNot(HaveOccurred())

		Ω(response.SessionID).Should(Equal("1"))
		Ω(response.Status.Code).Should(Equal(200))

		Ω(response.Result.Parameters).Should(HaveLen(2))
		Ω(response.Result.Parameters["param1"]).Should(Equal("value1"))

		Ω(response.Result.Fulfillment.Messages).Should(HaveLen(1))

		Ω(response.Result.Fulfillment.Messages[0].Speech).Should(Equal("Message speesh text"))
	})
})
