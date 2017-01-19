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
	"bufio"
	"encoding/binary"
	"io"
	"os"
)

type (
	chunkHeader struct {
		ChunkName [4]byte
		ChunkSize uint32
	}

	waveFormat struct {
		AudioFormat   uint16
		NumChannels   uint16
		SampleRate    uint32
		BytesPerSec   uint32
		BlockAlign    uint16
		BitsPerSample uint16
	}

	waveFileWriter struct {
		wavFile *os.File
		writer  *bufio.Writer
	}
)

var (
	riffChunkName       = []byte{'R', 'I', 'F', 'F'}
	waveFormatChunkName = []byte{'W', 'A', 'V', 'E'}
	fmtChunkName        = []byte{'f', 'm', 't', ' '}
	dataChunkName       = []byte{'d', 'a', 't', 'a'}

	pcm16WaveFormat = waveFormat{
		1,
		1,
		8000,
		16000,
		2,
		16,
	}
)

func createWavFile(wavFilePath string) (*waveFileWriter, error) {
	out, err := os.Create(wavFilePath)
	if err != nil {
		return nil, err
	}
	fw := &waveFileWriter{
		wavFile: out,
		writer:  bufio.NewWriter(out),
	}
	return fw, nil
}

func (wf *waveFileWriter) writeHeader(fmt waveFormat, soundLength uint32) error {

	fmtLength := uint32(binary.Size(fmt))
	var riffLength uint32
	riffLength = soundLength + 8 + 8 + fmtLength + 8

	if err := wf.writeChunkSize(riffChunkName, riffLength); err != nil {
		return err
	}
	riffLength -= 8
	if err := wf.writeChunkSize(waveFormatChunkName, riffLength); err != nil {
		return err
	}
	if err := wf.writeChunk(fmtChunkName, fmt); err != nil {
		return err
	}
	if err := wf.writeChunkSize(dataChunkName, soundLength); err != nil {
		return err
	}
	if err := wf.writer.Flush(); err != nil {
		return err
	}
	return nil
}

func (wf *waveFileWriter) writeChunkSize(name []byte, size uint32) error {
	h := chunkHeader{
		ChunkSize: size,
	}
	copy(h.ChunkName[:], name[:])
	return binary.Write(wf.wavFile, binary.LittleEndian, h)
}

func (wf *waveFileWriter) writeChunk(name []byte, data interface{}) error {
	err := wf.writeChunkSize(name, uint32(binary.Size(data)))
	if err != nil {
		return err
	}
	return binary.Write(wf.wavFile, binary.BigEndian, data)
}

func (wf *waveFileWriter) writeSoundData(r io.Reader) (n int64, err error) {
	wf.writer.Flush()
	n, err = io.Copy(wf.wavFile, r)
	return
}

func (wf *waveFileWriter) close() error {
	return wf.wavFile.Close()
}

func NewSpeeshToWaveFileHandler(wavFilePath string, reader io.Reader) (SpeechHandler, error) {
	wf, err := createWavFile(wavFilePath)
	if err != nil {
		return nil, err
	}
	return func(r io.Reader, length int64) {
		defer wf.close()
		wf.writeHeader(pcm16WaveFormat, uint32(length))
		wf.writeSoundData(r)

	}, nil
}
