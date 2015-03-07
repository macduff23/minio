/*
 * Mini Object Storage, (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package erasure1

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestSingleWrite(c *C) {
	var testBuffer bytes.Buffer
	testData := "Hello, World"
	testHeader := DataHeader{
		Key:              "testobj",
		ChunkIndex:       1,
		OriginalLength:   uint32(len(testData)),
		EncoderK:         8,
		EncoderM:         8,
		EncoderTechnique: Cauchy,
	}

	err := Write(&testBuffer, testHeader.Key, testHeader.ChunkIndex, testHeader.OriginalLength, testHeader.EncoderK, testHeader.EncoderM, testHeader.EncoderTechnique, bytes.NewBufferString(testData))
	c.Assert(err, IsNil)

	actualVersion := make([]byte, 4)
	_, err = testBuffer.Read(actualVersion)
	c.Assert(err, IsNil)
	c.Assert(binary.LittleEndian.Uint32(actualVersion), DeepEquals, uint32(1))

	actualHeader := DataHeader{}

	decoder := gob.NewDecoder(&testBuffer)
	decoder.Decode(&actualHeader)

	c.Assert(actualHeader, DeepEquals, testHeader)

	var actualData bytes.Buffer

	dataLength, err := io.Copy(&actualData, &testBuffer)
	c.Assert(dataLength, Equals, int64(len(testData)))
	c.Assert(actualData.Bytes(), DeepEquals, []byte(testData))
	c.Assert(err, IsNil)
}

func (s *MySuite) TestReadWrite(c *C) {
	var testBuffer bytes.Buffer
	testData := "Hello, World"
	testHeader := DataHeader{
		Key:              "testobj",
		ChunkIndex:       1,
		OriginalLength:   uint32(len(testData)),
		EncoderK:         8,
		EncoderM:         8,
		EncoderTechnique: Cauchy,
	}

	err := Write(&testBuffer, testHeader.Key, testHeader.ChunkIndex, testHeader.OriginalLength, testHeader.EncoderK, testHeader.EncoderM, testHeader.EncoderTechnique, bytes.NewBufferString(testData))
	c.Assert(err, IsNil)

	header, err := ReadHeader(&testBuffer)
	c.Assert(err, IsNil)
	c.Assert(header, DeepEquals, testHeader)
	c.Assert(testBuffer.Len(), Equals, len(testData))
	c.Assert(testBuffer.Bytes(), DeepEquals, []byte(testData))
}
