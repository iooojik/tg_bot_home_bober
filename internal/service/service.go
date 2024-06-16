package service

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"home_chief/internal/service/entity"
	"log/slog"
	"os"
	"time"
)

type BotService struct {
	db *sql.DB
}

func NewBotService() (*BotService, error) {
	db, err := sql.Open("mysql", viper.GetString("DB_DSN"))
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(10 * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	slog.Info("successfully connected to mysql")
	err = executeScripts([]string{"create_users_table.sql", "create_notif_table.sql"}, db)
	if err != nil {
		return nil, err
	}
	srv := &BotService{db: db}
	return srv, nil
}

func (bs *BotService) CheckUser(userId int) (string, error) {
	sqlResp, err := bs.db.Query(`SELECT * FROM users WHERE login =?;`, userId)
	if err != nil {
		return "", err
	}
	users, err := ReadRows[entity.UserRow](sqlResp)
	if err != nil {
		return "", err
	}
	if len(users) > 1 {
		panic("impossible")
	}
	sqlRes, err := bs.db.Exec(`INSERT INTO users (login) VALUES (?);`, userId)
	if err != nil {
		return "", err
	}
	lastId, err := sqlRes.LastInsertId()
	if err != nil {
		return "", err
	}
	slog.Info("created user with id", lastId)
	return "", nil
}

func (bs *BotService) ChangeDate(date, userId int) (string, error) {
	sqlRes, err := bs.db.Exec(`INSERT INTO notifications (num, user_login) VALUES (?, ?);`, userId, date)
	if err != nil {
		return "", err
	}
	lastId, err := sqlRes.LastInsertId()
	if err != nil {
		return "", err
	}
	slog.Info("created notification with id", lastId)
	return "", nil
}

func ReadRows[T any](rows *sql.Rows) ([]T, error) {
	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	cols := make([]*string, len(colNames))
	colPtrs := make([]interface{}, len(colNames))
	for i := 0; i < len(colNames); i++ {
		colPtrs[i] = &cols[i]
	}
	items := make([]T, 0)
	var myMap = make(map[string]any)
	for rows.Next() {
		scanErr := rows.Scan(colPtrs...)
		if scanErr != nil {
			return nil, scanErr
		}
		for i, col := range cols {
			myMap[colNames[i]] = col
		}
		rowItem := new(T)
		data, e := json.Marshal(myMap)
		if e != nil {
			return nil, e
		}
		e = json.Unmarshal(data, rowItem)
		if e != nil {
			return nil, e
		}
		items = append(items, *rowItem)
	}
	return items, nil
}

func executeScripts(scripts []string, cl *sql.DB) error {
	for _, scriptName := range scripts {
		scriptPath := viper.GetString("ASSETS_PATH") + "scripts/" + scriptName
		data, err := os.ReadFile(scriptPath)
		if err != nil {
			return err
		}
		_, err = cl.Query(string(data))
		if err != nil {
			return err
		}
	}
	return nil
}
