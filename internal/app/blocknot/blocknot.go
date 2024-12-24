package blocknot

import (
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Blocknot struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Blocknot {
	return &Blocknot{
		db: db,
	}
}

func (s *Blocknot) Open() error {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=midiy password=Qwe230405180405 dbname=my_chat port=5433 sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Blocknot) Close() {
	sqlDB, err := s.db.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}
