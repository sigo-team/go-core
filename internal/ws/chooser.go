package ws

import (
	"sigo/internal/lib"
)

func (lb *Lobby) SetChooser(player *Player) {

	lb.Chooser = player
	lb.ChooserBC = player.Sender

	content := lib.Request{
		Type: "setChooser",
		Data: lib.Data{
			PlayerId: lb.Chooser.ID,
		},
	}.Marshall()

	lb.SendAll(&message{
		UserID:  lb.Khil.ID,
		Content: content,
	})
}
