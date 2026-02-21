package main

import "time"

type application struct {
	config        config
	store         storage.Storage
	redis         redis.Redis
	authenticator auth.Authenticator
}

type config struct {
	addr     string
	db       dbConfig
	env      string
	apiURL   string
	auth     authConfig
	redisCfg redisConfig
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	user string
	pass string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
