package bot

import (
	"context"
	"familyFormUi/config"
	"log"
)

func (m *Members) getMember(key string) (err error) {
	err = config.Rdb.HGetAll(context.Background(), key).Scan(m)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
func (m *Members) Add(key string) (err error) {
	err = config.Rdb.HSet(context.Background(), key, m).Err()
	if err != nil {
		return err
	}
	return
}
