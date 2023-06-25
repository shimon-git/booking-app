package dbrepo

import (
	"errors"
	"time"

	"github.com/shimon-git/booking-app/internal/models"
)

// InsertReservation - insert a reservation into the database
func (r *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction = insert a room restriction into the DB.
func (r *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	return nil
}

// CheckAvialibilityForDatesByRoomID - return true if avialibility exist, otherwise return false
func (r *testDBRepo) CheckAvialibilityForDatesByRoomID(startDate, endDate time.Time, roomID int) (bool, error) {

	return false, nil
}

// CheckAvialibilityByDatesForAllRooms - return a slice of avilable rooms for a given date
func (r *testDBRepo) CheckAvialibilityByDatesForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID - gets a room by ID.
func (r *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room,errors.New("Some error")
	}
	
	return room, nil
}
