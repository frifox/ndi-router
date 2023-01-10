package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"sync"
)

type SQLite struct {
	//Address string

	log *zap.SugaredLogger
	db  *sql.DB
	sync.Mutex
}

func (s *SQLite) Init(logger *zap.SugaredLogger) (err error) {
	s.log = logger.Named("sqlite")

	log := s.log.Named("Init")
	log.Infow("starting")

	//s.Address = "file:ndi-router.db?_mutex=full"

	// create db file
	s.Open()

	_, err = s.db.Exec(`create table if not exists matrix ("id" integer, "output" string, "input" string, primary key (id))`)
	if err != nil {
		log.Errorw("db.Exec(create table)", "err", err)
		return err
	}

	return nil
}

func (s *SQLite) Open() {
	log := s.log.Named("Open")

	var err error
	s.db, err = sql.Open("sqlite3", conf.SqliteDsn)
	if err != nil {
		log.Errorw("sql.Open", "err", err)
	}
}
func (s *SQLite) Close() {
	log := s.log.Named("Close")

	err := s.db.Close()
	if err != nil {
		log.Errorw("sql.Close", "err", err)
	}
}

func (s *SQLite) UpdateMatrix(output string, input string) {
	log := s.log.Named("UpdateMatrix()")

	// delete old config
	del, err := s.db.Prepare("delete from matrix where output = ?")
	_, err = del.Exec(output)
	if err != nil {
		log.Errorw("del", "err", err)
		return
	}

	// save new config
	insert, err := s.db.Prepare("insert into matrix (output, input) values (?, ?)")
	res, err := insert.Exec(output, input)
	if err != nil {
		log.Errorw("insert", "err", err)
		return
	}

	n, err := res.RowsAffected()
	if err != nil {
		log.Errorw("res.RowsAffected", "err", err)
		return
	}

	log.Infow("saved", "rows", n)
}

func (s *SQLite) GetMatrix(output string) (input string) {
	log := s.log.Named("GetMatrix()")

	query, err := s.db.Prepare("select input from matrix where output = ?")
	rows, err := query.Query(output)
	if err != nil {
		log.Errorw("select", "err", err)
		return ""
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&input)
		if err != nil {
			log.Errorw("can't Scan", "err", err)
			continue
		}

		return input
	}

	return ""
}
