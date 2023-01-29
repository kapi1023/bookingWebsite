package repository

import (
	"time"

	"github.com/kapi1023/bookingWebsite/internal/models"
)

type DataBaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailalilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error)
	SearchAvailalilityByDatesForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)
}
