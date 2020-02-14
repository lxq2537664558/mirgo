package mir

import (
	"sync"

	"github.com/yenkeia/mirgo/common"
)

func DetectMapVersion(input []byte) byte {
	//c# custom map format
	if (input[2] == 0x43) && (input[3] == 0x23) {
		return 100
	}
	//wemade mir3 maps have no title they just start with blank bytes
	if input[0] == 0 {
		return 5
	}
	//shanda mir3 maps start with title: (C) SNDA, MIR3.
	if (input[0] == 0x0F) && (input[5] == 0x53) && (input[14] == 0x33) {
		return 6
	}

	//wemades antihack map (laby maps) title start with: Mir2 AntiHack
	if (input[0] == 0x15) && (input[4] == 0x32) && (input[6] == 0x41) && (input[19] == 0x31) {
		return 4
	}

	//wemades 2010 map format i guess title starts with: Map 2010 Ver 1.0
	if (input[0] == 0x10) && (input[2] == 0x61) && (input[7] == 0x31) && (input[14] == 0x31) {
		return 1
	}

	//shanda's 2012 format and one of shandas(wemades) older formats share same header info, only difference is the filesize
	if (input[4] == 0x0F) && (input[18] == 0x0D) && (input[19] == 0x0A) {
		W := int(input[0] + (input[1] << 8))
		H := int(input[2] + (input[3] << 8))
		if len(input) > (52 + (W * H * 14)) {
			return 3
		} else {
			return 2
		}
	}

	//3/4 heroes map format (myth/lifcos i guess)
	if (input[0] == 0x0D) && (input[1] == 0x4C) && (input[7] == 0x20) && (input[11] == 0x6D) {
		return 7
	}

	return 0
}

func GetMapV1(bytes []byte) *Map {
	offset := 21
	w := common.BytesToUint16(bytes[offset : offset+2])
	offset += 2
	xor := common.BytesToUint16(bytes[offset : offset+2])
	offset += 2
	h := common.BytesToUint16(bytes[offset : offset+2])
	width := int(w ^ xor)
	height := int(h ^ xor)

	m := NewMap(width, height)

	offset = 54
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			p := common.Point{X: uint32(i), Y: uint32(j)}
			c := new(Cell)
			c.Map = m
			c.Point = p
			c.Objects = new(sync.Map)
			if (common.BytesToUint32(bytes[offset:offset+4])^0xAA38AA38)&0x20000000 != 0 {
				c.Attribute = common.CellAttributeHighWall
			}
			if ((common.BytesToUint16(bytes[offset+6:offset+8]) ^ xor) & 0x8000) != 0 {
				c.Attribute = common.CellAttributeLowWall
			}
			if c.Attribute == common.CellAttributeWalk {
				m.SetCell(p, c)
			}
			offset += 15
		}
	}
	m.Width = width
	m.Height = height
	return m
}
