package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	Id      int `json:"id"`
	UserId  int `json:"userId"`
	EventId int `json:"eventId"`
}

func (am *AttendeeModel) Insert(a *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO attendees (event_id,user_id ) VALUES ($1, $2) RETURNING id`

	err := am.DB.QueryRowContext(ctx, query, a.EventId, a.UserId).Scan(&a.Id)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (am *AttendeeModel) GetByEvent(eventId int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT u.id, u.name, u.email FROM users u
		JOIN attendees a ON u.id = a.user_id
		WHERE a.event_id = $1`

	rows, err := am.DB.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (am *AttendeeModel) GetByEventAndUser(eventId, userId int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM attendees WHERE event_id = $1 AND user_id = $2`

	var a Attendee
	err := am.DB.QueryRowContext(ctx, query, eventId, userId).Scan(&a.Id, &a.UserId, &a.EventId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &a, nil
}

func (am *AttendeeModel) Delete(userId, eventId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM attendees WHERE user_id = $1 AND event_id = $2`

	_, err := am.DB.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return err
	}

	return nil
}

func (am *AttendeeModel) GetEventsByUserId(userId int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT e.id, e.name, e.description, e.date, e.location 
		FROM events e
		JOIN attendees a ON e.id = a.event_id
		WHERE a.user_id = $1`

	rows, err := am.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event

	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.Id, &e.Name, &e.Description, &e.Date, &e.Location); err != nil {
			return nil, err
		}
		events = append(events, &e)
	}

	return events, nil
}
