package dbrepo

import (
	"context"
	"time"

	"github.com/shimon-git/booking-app/internal/models"
)

const timeout = 5 * time.Second

func (r *postgresDBRepo) InsertReservation(res models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	query := `insert into reservations (first_name, last_name, email, phone, start_date,
			  end_date, room_id, created_at, updated_at)
			  values($1,$2,$3,$4,$5,$6,$7,$8,$9);`
	_, err := r.DB.ExecContext(ctx, query,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.CreatedAt,
		res.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
