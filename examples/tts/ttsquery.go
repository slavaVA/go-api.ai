package main
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
	"flag"
	"fmt"
	"github.com/slavaVA/go-api.ai"
	"os"
)

var (
	accessToken = flag.String("accessToken", "", "Client access token")
	queryText = flag.String("text", "", "Query text")
	lang = flag.String("lang","en-US", "Query language")
	fileName=flag.String("fileName","out.wav", "Out WAVE file name")
)


func main() {
	flag.Parse()
	if len(*accessToken)==0 {
		fmt.Println("Access token required!")
		return
	}

	fmt.Println("TTS API Endpoint : Lang=",*lang," text=",*queryText," accessToken=",*accessToken)

	exist,sl:=gapiai.IsLanguageSupport(*lang)
	if exist==false {
		fmt.Println("Language not supported: ",*lang)
		return
	}

	cfg:=&gapiai.ApiConfig{
		AccessToken:*accessToken,
		Lang:sl,
	}
	endPoint:=gapiai.DefaultTtsAPIEndpoint(cfg)

	endPoint.EnableLogger(os.Stdout)

	sh,err:=gapiai.NewSpeechToWaveFileHandler(*fileName)
	if err!=nil {
		panic(err)
	}
	if err:=endPoint.DoTts(*queryText,sh);err!=nil{
		panic(err)
	}
}
