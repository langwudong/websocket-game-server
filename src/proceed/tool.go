package proceed

func Contain(rooms []*Room, roomId string) (bool, *Room) {
	if rooms != nil {
		for _, room := range rooms {
			if room.ID == roomId {
				return true, room
			}
		}
	}
	return false, nil
}

func QueryRoom(rooms []*Room, player *Player) *Room {
	if rooms != nil {
		for _, room := range rooms {
			if room.ID == player.RoomID {
				return room
			}
		}
	}
	return nil
}

func RemoveRoom(rooms *[]*Room, room *Room) {
	for index, value := range *rooms {
		if value == room {
			*rooms = append((*rooms)[:index], (*rooms)[index+1:]...)
		}
	}
}
