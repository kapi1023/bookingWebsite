package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kapi1023/bookingWebsite/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertsReservations into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	var newId int

	stmt := `insert into reservations(first_name,last_name,email,phone,start_date,end_date,room_id,created_at,updated_at)
			values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	stmt := `insert into room_restrictions(start_date,end_date,room_id,reservation_id,created_at,updated_at,restriction_id)
			values($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomId,
		r.ReservationId,
		time.Now(),
		time.Now(),
		r.RestrictionId,
	)
	if err != nil {
		return err
	}
	return nil
}

// SearchAvailalilityByDatesByRoomId return true if availability exist for roomId and false if no availability exist for roomId
func (m *postgresDBRepo) SearchAvailalilityByDatesByRoomId(start, end time.Time, roomId int) (bool, error) {
	var numRows int
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	query := `select count(id) from room_restrictions where room_id=$1 and  $2<end_date and $3<start_date;`

	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}

	return false, nil

}

// SearchAvailalilityByForAllRooms returns a slice of avaible rooms if any avaible for date range
func (m *postgresDBRepo) SearchAvailalilityByDatesForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	query := `select r.id,r.room_name from rooms r where r.id not in(select room_id from room_restrictions rr where $1<rr.end_date and $2>rr.start_date);`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.Id,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {

		return rooms, err
	}
	fmt.Println("-----------------------")
	return rooms, nil
}

func (m *postgresDBRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var room models.Room
	query := `select id,room_name,created_at,updated_at from rooms where Id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.Id,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}
	return room, nil
}

// GetUserById returns the user by id
func (m *postgresDBRepo) GetUserById(id string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	query := `select id,first_name,last_name,email,password,access_level,created_at,updated_at from users where id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	var user models.User
	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil

}

// UpdateUser updates user in database
func (m *postgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	query := `update user set first_name = $1, last_name = $2, email = $3,access_level = $4, updated_at = $5 from user`

	_, err := m.DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

// Authenicate authenicate a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	query := `select id, password from users where email = $1`
	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPassword, nil
}

// Returns a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var reservations []models.Reservation
	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at,rm.id,rm.room_name
	from reservations r
	left join rooms rm on rm.id = r.room_id
	order by r.start_date asc
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.Id,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.Id,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}
