package mysql

import (
	"database/sql"
	"errors"
	"fmt"
)

// Config is the configuration structure for the MySQL asset index.
//
//swagger:ignore
type Config struct {
	// Host is the host name for the MySQL connection
	Host string `json:"host" yaml:"host" env:"MYSQL_HOST"`
	// Port is the TCP port number for the MySQL connection.
	Port int `json:"port" yaml:"port" env:"MYSQL_PORT"`
	// Username is the user for authentication.
	Username string `json:"user" yaml:"user" env:"MYSQL_USER"`
	// Password is the password for the Username
	Password string `json:"password" yaml:"password" env:"MYSQL_PASSWORD"`
	// Database is the database name the asset index will be stored in.
	// This database will be created if it does not exist.
	Database string `json:"db" yaml:"db" env:"MYSQL_DATABASE"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Username == "" {
		return errors.New("username cannot be empty")
	}
	if c.Database == "" {
		return errors.New("database cannot be empty")
	}
	return nil
}

// ConnectString returns a MySQL connect string from the settings.
func (c *Config) ConnectString() string {
	host := c.Host
	if c.Port > 0 {
		host = fmt.Sprintf("%s:%d", host, c.Port)
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", c.Username, c.Password, host, c.Database)
}

func (c *Config) CreateMySQLDB() (*sql.DB, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	db, err := sql.Open(
		"mysql",
		c.ConnectString(),
	)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(
		`CREATE DATABASE IF NOT EXISTS ` + c.Database,
	); err != nil {
		return nil, fmt.Errorf("failed to create database (%w)", err)
	}
	return db, nil
}
