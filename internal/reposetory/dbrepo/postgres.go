package dbrepo

import (
	"context"
	"time"

	"github.com/shimon-git/booking-app/internal/models"
)

const timeout = 5 * time.Second

// InsertReservation - insert a reservation into the database
func (r *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var reservationID int
	query := `insert into reservations (first_name, last_name, email, phone, start_date,
			  end_date, room_id, created_at, updated_at)
			  values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id;`
	err := r.DB.QueryRowContext(ctx, query,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.CreatedAt,
		res.UpdatedAt).Scan(&reservationID)
	if err != nil {
		return 0, err
	}
	return reservationID, nil
}

// InsertRoomRestriction = insert a room restriction into the DB.
func (r *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `insert into room_restrictions
	(start_date,end_date,room_id,reservation_id,created_at, updated_at,restriction_id)
	values ($1,$2,$3,$4,$5,$6,$7);`
	_, err := r.DB.ExecContext(ctx, query,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		res.CreatedAt,
		res.UpdatedAt,
		res.RestrictionID,
	)
	if err != nil {
		return err
	}
	return nil
}

// CheckAvialibilityForDatesByRoomID - return true if avialibility exist, otherwise return false
func (r *postgresDBRepo) CheckAvialibilityForDatesByRoomID(startDate, endDate time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var numOfRows int
	query := `select count(id) from room_restrictions
			  where room_id = $1
			  and $2 <= end_date and $3 >= start_date;`
	err := r.DB.QueryRowContext(ctx, query, startDate, endDate).Scan(&numOfRows)
	if err != nil {
		return false, err
	}
	return numOfRows == 0, nil
}

// CheckAvialibilityByDatesForAllRooms - return a slice of avilable rooms for a given date
func (r *postgresDBRepo) CheckAvialibilityByDatesForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var rooms []models.Room
	query := `select r.id, r.room_name from rooms r
	where r.id not in (
		select room_id from room_restrictions rr 
		where $1 <= rr.end_date 
		and $2 >= rr.start_date);`
	rows, err := r.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.RoomName); err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

// GetRoomByID - gets a room by ID.
func (r *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var room models.Room
	query := `select id,room_name,created_at,updated_at from rooms where id = $1;`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt)
	if err != nil {
		return room, err
	}
	return room, nil
}
