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
	"encoding/json"
	"time"
)

type (
	//Query the following parameters are used as either query parameters in the URL or JSON keys in the POST body
	Query struct {
		Query         []string        `json:"query"`
		Confidence    []float32       `json:"confidence,omitempty"`
		Contexts      []DialogContext `json:"contexts,omitempty"`
		ResetContexts bool            `json:"resetContexts,omitempty"`
		Event         *Event          `json:"event,omitempty"`
		Timezone      string          `json:"timezone,omitempty"`
		Lang          string          `json:"lang"`
		SessionID     string          `json:"sessionId"`
		Entities      []Entity        `json:"entities,omitempty"`
		Location      *Location       `json:"location,omitempty"`
	}

	//DialogContext are strings that represent the current context of a user’s request. This is helpful for
	//differentiating phrases which may be vague or have different meanings depending on the user’s
	//preferences or geographic location, the current page in an app, or the topic of conversation.
	//
	//For example, if a user is listening to a music player application and finds a band that catches their
	//interest, they might say something like: “I want to hear more of them”. As a developer, you can include
	//the name of the band in the context with the request, so that the API.AI agent can process it more effectively.
	//
	//Or let’s say you’re a manufacturer of smart home devices, and you have an app that remotely controls
	//household appliances. A user might say, "Turn on the front door light", followed by “Turn it off”,
	//and the app will understand that the second phrase is still referring to the light. Now later,
	//if the user says, "Turn on the coffee machine", and follows this with “Turn it off”,
	//it will result in different action than before, because of the new context.
	DialogContext struct {
		Name       string                 `json:"name"`
		Parameters map[string]interface{} `json:"parameters"`
		Lifespan   int                    `json:"lifespan"`
	}

	//Event is a feature that allows you to invoke intents by an event name instead of a user query.
	//First, you define event names in intents. Then, you can trigger these intents by sending a
	///query request containing an "event" parameter.
	Event struct {
		Name string            `json:"name"`
		Data map[string]string `json:"data"`
	}

	//Location contains latitude and longitude values.
	//Example: {"latitude": 37.4256293, "longitude": -122.20539}
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	//Entity represent concepts and serve as a powerful tool for extracting parameter values from natural language inputs.
	//
	//The entities that are used in a particular agent will depend on the parameter values that are expected to be
	//returned as a result of agent functioning. In other words, a developer need not create entities for every
	//concept mentioned in the agent – only for those required for actionable data.
	//
	//There are 3 types of entities: system (defined by API.AI), developer (defined by a developer),
	//and user (built for each individual end-user in every request) entities. Furthermore, each of these can be
	//mapping (having reference values), enum type (having no reference values), or composite
	//(containing other entities with aliases and returning object type values).
	//
	//Array of entities that replace developer defined entities for this request only.
	//The entity(ies) need to exist in the developer console.
	Entity struct {
		Name    string        `json:"name"`
		Entries []EntityEntry `json:"entries"`
		Extend  bool          `json:"extend"`
		IsEnum  bool          `json:"isEnum"`
	}

	EntityEntry struct {
		Value    string   `json:"value"`
		Synonyms []string `json:"synonyms"`
	}

	//QueryResponse takes natural language text and information as JSON in the POST body and returns information as JSON.
	QueryResponse struct {
		ID        string       `json:"id"`
		Timestamp time.Time    `json:"timestamp"`
		Result    QueryResult  `json:"result"`
		Status    StatusObject `json:"status"`
		SessionID string       `json:"sessionId"`
	}

	QueryResult struct {
		Source           string                 `json:"source"`
		ResolvedQuery    string                 `json:"resolvedQuery"`
		Action           string                 `json:"action"`
		ActionIncomplete bool                   `json:"actionIncomplete"`
		Parameters       map[string]interface{} `json:"parameters"`
		Contexts         []DialogContext        `json:"contexts"`
		Fulfillment      Fulfillment            `json:"fulfillment"`
		Metadata         Metadata               `json:"metadata"`
	}

	Fulfillment struct {
		Speech      string     `json:"speech"`
		DisplayText string     `json:"displayText"`
		Source      string     `json:"source"`
		Data        string     `json:"data"`
		Messages    []Messages `json:"messages"`
	}

	Metadata struct {
		IntentID                  string `json:"intentId"`
		IntentName                string `json:"intentName"`
		WebhookUsed               string `json:"webhookUsed"`
		WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
	}

	//StatusObject is returned with every request and indicates if the request was successful.
	//If it is not successful, error information is included.
	//See Status and Error Codes for more information on the returned errors.
	StatusObject struct {
		Code         int    `json:"code"`
		ErrorType    string `json:"errorType"`
		ErrorID      string `json:"errorId"`
		ErrorDetails string `json:"errorDetails"`
	}

	//Messages is text response message object
	Messages struct {
		Type   int    `json:"type"`
		Speech string `json:"speech"`
	}

	//QueryAPIEndpoint is used to process natural language in the form of text. The query requests return structured data in JSON format with an action and parameters for that action.
	QueryAPIEndpoint interface {
		DoQuery(q Query) (*QueryResponse, error)
		TextRequest(sessionID string, text string) (*QueryResponse, error)
	}

	SupportedLang string
)

const (
	apiAiURL          = "https://api.api.ai/v1/"
	CurrentAPIVersion = "20150910"

	English          SupportedLang = "en"
	Russian          SupportedLang = "ru"
	German           SupportedLang = "de"
	Portuguese       SupportedLang = "pt"
	PortugueseBrazil SupportedLang = "pt-BR"
	Spanish          SupportedLang = "es"
	French           SupportedLang = "fr"
	Italian          SupportedLang = "it"
	Japanese         SupportedLang = "ja"
	Korean           SupportedLang = "ko"
	ChineseChina     SupportedLang = "zh-CN"
	ChineseHongKong  SupportedLang = "zh-HK"
	ChineseTaiwan    SupportedLang = "zh-TW"
)

func (status *StatusObject) IsSuccess() bool {
	return status.Code < 400
}

func (response *QueryResponse) Decode(data []byte) (err error) {
	err = json.Unmarshal(data, response)
	return
}
