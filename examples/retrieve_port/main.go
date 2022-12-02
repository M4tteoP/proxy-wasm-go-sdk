// Copyright 2020-2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/binary"
	"errors"
	"math"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &httpContext{}
}

// Override types.DefaultPluginContext.
func (*httpContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{
		contextID: contextID,
	}
}

type httpContext struct {
	types.DefaultPluginContext
	contextID uint32
}

func (ctx *httpContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogInfo("OnPluginStart from Go!")
	return types.OnPluginStartStatusOK
}

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("OnHttpRequestHeaders")
	RetrieveAndPrint()
	return types.ActionContinue
}

func (ctx *httpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	proxywasm.LogInfo("OnHttpRequestBody")
	RetrieveAndPrint()
	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("OnHttpResponseHeaders")
	RetrieveAndPrint()
	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	proxywasm.LogInfo("OnHttpResponseBody")
	RetrieveAndPrint()
	return types.ActionContinue
}

func RetrieveAndPrint() {
	srcAddressRaw, err := proxywasm.GetProperty([]string{"source", "address"})
	if err != nil {
		proxywasm.LogInfof("failed to get source address: %v", err)
	}
	srcPortRaw, err := proxywasm.GetProperty([]string{"source", "port"})
	if err != nil {
		proxywasm.LogInfof("failed to get port address: %v", err)
	}
	dstAddressRaw, err := proxywasm.GetProperty([]string{"destination", "address"})
	if err != nil {
		proxywasm.LogInfof("failed to get destination address: %v", err)
	}
	dstPortRaw, err := proxywasm.GetProperty([]string{"destination", "port"})
	if err != nil {
		proxywasm.LogInfof("failed to get port address: %v", err)
	}
	srcAddress := string(srcAddressRaw)
	srcPort, _ := convertPort(srcPortRaw)
	dstAddress := string(dstAddressRaw)
	dstPort, _ := convertPort(dstPortRaw)
	proxywasm.LogInfof("\nSource Address: %v %s \nSource Port: %v %d\nDestination Address: %v %s\nDestination Port: %v %d\n", srcAddressRaw, srcAddress, srcPortRaw, srcPort, dstAddressRaw, dstAddress, dstPortRaw, dstPort)
}

func (ctx *httpContext) OnHttpRequestTrailers(numTrailers int) types.Action {
	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseTrailers(numTrailers int) types.Action {
	return types.ActionContinue
}

func (ctx *httpContext) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}

// UintToInt32 converts uint to int32.
func convertPort(b []byte) (int, error) {
	// 0 < Port number <= 65535, therefore the retrieved value should never exceed 16 bits
	// and correctly fit int (at least 32 bits in size)
	unsignedInt := binary.LittleEndian.Uint32(b)
	if unsignedInt > math.MaxInt32 {
		return 0, errors.New("port convertion")
	}

	return int(unsignedInt), nil
}
