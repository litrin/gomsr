/*

	BSD 2-Clause License

	Copyright (c) 2020, Litrin Jiang
	All rights reserved.

	Redistribution and use in source and binary forms, with or without
	modification, are permitted provided that the following conditions are met:

	1. Redistributions of source code must retain the above copyright notice, this
	list of conditions and the following disclaimer.

	2. Redistributions in binary form must reproduce the above copyright notice,
	this list of conditions and the following disclaimer in the documentation
	and/or other materials provided with the distribution.

	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
	AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
	IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
	DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
	FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
	DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
	SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
	CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
	OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/

package gomsr

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func NewMSR(num CPUNum, offset int64) (*Operator, error) {
	var err error
	if os.Getuid() != 0 {
		err = errors.New("Permission deny, root permission is required")
		return nil, err
	}

	filename := fmt.Sprintf("/dev/cpu/%d/msr", num)
	_, err = os.Stat(filename)
	if err != nil {
		err = errors.New("msr file is not exist, run cmd `modprobe msr` before using")
		return nil, err
	}
	return &Operator{num, filename, offset, nil}, nil
}

func (msr *Operator) Read() error {

	fd, err := os.Open(msr.Filename)
	defer fd.Close()

	if err != nil {
		return err
	}

	fd.Seek(msr.Offset, io.SeekCurrent)
	content := make(Content, MSRLength)
	_, err = fd.Read(content)
	if err != nil {
		return err
	}
	msr.Value = content
	return nil
}

func (msr *Operator) Write(value Content) error {

	fd, err := os.OpenFile(msr.Filename, os.O_WRONLY, 0)
	defer fd.Close()

	if err != nil {
		return err
	}
	fd.Seek(msr.Offset, io.SeekCurrent)
	_, err = fd.Write(value)
	if err != nil {
		return err
	}

	return nil
}

func NewMSRContent(value uint64) Content {
	M := make(Content, MSRLength)
	for i := 0; i < MSRLength; i++ {
		M[i] = byte(value & 0xff)
		value >>= 8
	}
	return M
}

func (m Content) Value() uint64 {
	value := uint64(0)
	for i, v := range m {
		t := uint64(v)
		value += t << (i * 8)
	}

	return value
}

func (m Content) String() string {
	return fmt.Sprintf("0x%X", m.Value())
}
