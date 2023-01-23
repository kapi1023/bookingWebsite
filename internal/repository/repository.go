package repository

import "github.com/kapi1023/bookingWebsite/internal/models"

type DataBaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}
