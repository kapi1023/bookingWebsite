package dbrepo

import (
	"errors"
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
	if id > 2 {
		return room, errors.New("Room not found")
	}
	return room, nil
}

func (m *testDBRepo) GetUserById(id string) (models.User, error) {
	var user models.User

	return user, nil
}
func (m *testDBRepo) UpdateUser(user models.User) error {
	return nil
}
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	return 1, "", nil
}

// Returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}
