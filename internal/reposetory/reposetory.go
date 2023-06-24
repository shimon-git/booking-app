package reposetory

import (
	"time"

	"github.com/shimon-git/booking-app/internal/models"
)

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	CheckAvialibilityForDatesByRoomID(startDate, endDate time.Time, roomID int) (bool, error)
	CheckAvialibilityByDatesForAllRooms(startDate, endDate time.Time) ([]models.Room, error)
}
