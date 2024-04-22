package ws

import "sigo/internal/api"

func (lb *Lobby) SetChooser(player *Player) {

	lb.Chooser = player
	lb.ChooserBC = player.Sender

	content := api.Request{
		Type: "setChooser",
		Data: api.Data{
			PlayerId: lb.Chooser.ID,
		},
	}.Marshall()

	lb.SendAll(&message{
		UserID:  lb.Khil.ID,
		Content: content,
	})
}
