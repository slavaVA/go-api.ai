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

	"bytes"
	"encoding/binary"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("Wavefile", func() {
	It("Should write speech to file", func() {
		soundData := make([]byte, 100)
		r := bytes.NewReader(soundData)
		dir, err := ioutil.TempDir("", "wavedir")
		Ω(err).ShouldNot(HaveOccurred())

		defer os.RemoveAll(dir)

		tmpfn := filepath.Join(dir, "test.wav")

		sh, err := NewSpeeshToWaveFileHandler(tmpfn, r)
		Ω(err).ShouldNot(HaveOccurred())

		sh(r, r.Size())

		bf, err := ioutil.ReadFile(tmpfn)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(bf).Should(HaveLen(148))

		Ω(string(bf[:4])).Should(Equal("RIFF"))
		Ω(binary.LittleEndian.Uint32(bf[4:8])).Should(Equal(uint32(140)))

		Ω(string(bf[8:12])).Should(Equal("WAVE"))
		Ω(binary.LittleEndian.Uint32(bf[12:16])).Should(Equal(uint32(132)))

		Ω(string(bf[16:20])).Should(Equal("fmt "))
		Ω(binary.LittleEndian.Uint32(bf[20:24])).Should(Equal(uint32(16)))

	})
})
