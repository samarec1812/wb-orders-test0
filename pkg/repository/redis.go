package repository

import (
	"github.com/go-redis/redis"
)

type ConfigRedis struct {
	Addr     string
	Password string
	DB int
}



func NewRedisDB(cfg ConfigRedis) (*redis.Client, error) {
	//db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
	//	cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	//if err != nil {
	//	return nil, err
	//}
	//err = db.Ping()
	//if err != nil {
	//	return nil, err
	//}
	//DB, err := strconv.Atoi(cfg.DB)
	//if err != nil {
	//	return nil, fmt.Errorf("")
	//}

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		Password: cfg.Password,
		DB: cfg.DB,
	})

	err := redisClient.Ping().Err()
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}
