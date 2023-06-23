package reposetory

import "github.com/shimon-git/booking-app/internal/models"

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) error
}
