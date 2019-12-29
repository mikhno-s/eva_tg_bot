package app

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mikhno-s/eva_tg_bot/app/scheme"
)

type carsDB struct {
	ConnString string
	Conn       *sqlx.DB
}

func (DB *carsDB) Init() error {
	conn, err := sqlx.Open("postgres", DB.ConnString)
	if err != nil {
		return err
	}
	DB.Conn = conn
	return nil
}

func (DB *carsDB) Close() error {
	return DB.Conn.Close()
}

func (DB *carsDB) insertEvacuatedCars(cars []*scheme.Car) error {
	err := DB.Conn.Ping()
	if err != nil {
		return err
	}

	_, err = DB.Conn.Query("SELECT 'evacuated_cars'::regclass")
	if err != nil {
		return err
	}

	txn, err := DB.Conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := DB.Conn.PrepareNamed(`INSERT INTO evacuated_cars VALUES (:date, :model, :license_plate, :vin, :entity_id)`)

	for _, c := range cars {
		_, err = stmt.Exec(c)
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (DB *carsDB) getLastEvacuatedCar() (*scheme.Car, error) {
	car := scheme.Car{}
	err := DB.Conn.Get(&car, "SELECT date, model, license_plate, vin, entity_id FROM evacuated_cars ORDER BY id DESC limit 1")
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &car, nil
}

func (DB *carsDB) getLastReadedMessageID(chatID int64) (int64, error) {
	row := DB.Conn.QueryRow("SELECT last_message_id FROM telegram_last_chat_messages WHERE chat_id=$1",
		strconv.Itoa(int(chatID)))

	var messageID int64
	err := row.Scan(&messageID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, nil
		}
		return 0, err
	}
	return messageID, nil
}

func (DB *carsDB) SetLastReadedMessageID(chatID int64, messageID string) {
	DB.Conn.MustExec("INSERT INTO telegram_last_chat_messages VALUES($1, $2) ON CONFLICT (chat_id) DO UPDATE SET last_message_id=$2",
		strconv.Itoa(int(chatID)), messageID)
	// DB.Conn.MustExec("UPDATE telegram_last_chat_messages SET last_message_id=$1 WHERE chat_id=$2",
	// 	messageID, strconv.Itoa(int(chatID)))
}

// func (DB *carsDB) InsertLastReadedMessageID(chatID int64, messageID string) {
// 	DB.Conn.MustExec("INSERT INTO telegram_last_chat_messages VALUES($1, $2)",
// 		strconv.Itoa(int(chatID)), messageID)
// }

func (DB *carsDB) getAllCars() ([]*scheme.Car, error) {
	cars := []*scheme.Car{}
	err := DB.Conn.Select(&cars, "SELECT date, model, license_plate, vin, entity_id FROM evacuated_cars ORDER BY id ASC")
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return cars, nil
		}
		return nil, err
	}

	return cars, nil
}
