package session

import "log"

func (s *Session) Begin() (err error) {
	log.Printf("start tx")
	s.tx, err = s.db.Begin()
	if err != nil {
		log.Printf("begin tx err: %s", err.Error())
		return err
	}
	return
}
func (s *Session) Commit() error {
	log.Printf("commit tx")
	if err := s.tx.Commit(); err != nil {
		log.Printf("commit tx err: %s", err.Error())
		return err
	}
	return nil
}
func (s *Session) Rollback() error {
	log.Printf("rollback tx")
	if err := s.tx.Rollback(); err != nil {
		log.Printf("rollback tx err: %s", err.Error())
		return err
	}
	return nil
}
