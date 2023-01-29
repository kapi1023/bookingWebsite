package dbrepo

import (
	"time"

	"github.com/kapi1023/bookingWebsite/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertsReservations into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	return nil
}

// SearchAvailalilityByDatesByRoomId return true if availability exist for roomId and false if no availability exist for roomId
func (m *testDBRepo) SearchAvailalilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {

	return false, nil

}

// SearchAvailalilityByForAllRooms returns a slice of avaible rooms if any avaible for date range
func (m *testDBRepo) SearchAvailalilityByDatesForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room

	return room, nil
}
