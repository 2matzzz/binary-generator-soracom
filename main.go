package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

var soracomUnifiedUDPEndpointFQDN = "uni.soracom.io:23080"

func main() {
	// [変数名]:[バイト列のインデックス]:[変数の型]:[型に依存した設定内容]
	// flag:0:bool:7 temp:1:int:13:/10 humid:3:uint:8:/100 lat::float:32 long:float:32
	//
	// flag:0:bool:7 [offset(required)]
	// temp:1:int:13:/10 [length(required)][offset(optional)][operations(optional)]
	// humid:3:uint:8:/100 [length(required)][offset(optional)][operations(optional)]
	// lat::float:32 [length(required)][offset(optional)][operations(optional)]
	// long:float:32 [length(required)][offset(optional)][operations(optional)]
	//
	// char ASCII 1 byte = 1 character

	// [flag  ] [temp        ]    [humid ] [lat                              ] [long                             ]
	// 10000000 00001111 11111000 01111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111
	// みたいなかんじになる

	var data = []byte{0x80, 0x07, 0xF8, 0x5F} // {"flag":true,"temp":25.5,"humid":0.95}
	var lat = []byte{0x42, 0x0E, 0xAD, 0xFF}  // {"lat":35.669918060302734}
	var lon = []byte{0x43, 0x0B, 0xBE, 0x86}  // {"long":139.74423217773438}

	for _, v := range lat {
		data = append(data, v)
	}

	for _, v := range lon {
		data = append(data, v)
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		panic(err)
	}

	conn, _ := net.Dial("udp", soracomUnifiedUDPEndpointFQDN)
	defer conn.Close()

	conn.Write(data)

	buffer := make([]byte, 1500)
	length, _ := conn.Read(buffer)
	fmt.Printf("Receive: %s", string(buffer[:length]))
}
