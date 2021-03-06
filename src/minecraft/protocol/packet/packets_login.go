package packet

import (
	chat2 "net/hyren/nyrah/minecraft/chat"
	codecs2 "net/hyren/nyrah/minecraft/protocol/codecs"
)

type LoginStart struct {
	Username codecs2.String
}

func (_ LoginStart) ID() int { return 0x00 }

type LoginSuccess struct {
	UUID     codecs2.String
	Username codecs2.String
}

func (_ LoginSuccess) ID() int { return 0x02 }

type LoginDisconnect struct {
	Chat chat2.TextComponent
}

func (_ LoginDisconnect) ID() int { return 0x00 }
