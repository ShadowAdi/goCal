package migrations

import (
	"context"
	"fmt"
	"goCal/internal/db"
	"goCal/internal/logger"
)

func CreateUserTable() {
	ctx := context.Background()

	_, errEnum := db.Conn.Exec(ctx, `
	CREATE TYPE pronoun_enum AS ENUM ('he/him', 'she/her', 'they/them', 'other');
	CREATE TYPE date_format_enum AS ENUM ('DD/MM/YYYY', 'MM/DD/YYYY', 'YYYY-MM-DD');
	CREATE TYPE time_format_enum AS ENUM ('12h', '24h');
	`)

	if errEnum != nil {
		fmt.Printf("Error creating enums %s", errEnum)
		logger.Error("Failed to create user table: " + errEnum.Error())
	}

	_, err := db.Conn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	profileUrl TEXT,
	country TEXT NOT NULL,
    welcome_message TEXT DEFAULT 'Welcome to my scheduling page. Please follow the instructions to add an event to my calendar.',
	timezone TEXT DEFAULT 'UTC',
	pronouns pronoun_enum DEFAULT 'other',
	isVerified BOOLEAN DEFAULT False,
	date_format date_format_enum DEFAULT 'DD/MM/YYYY',
    time_format time_format_enum DEFAULT '24h',
	custom_link TEXT UNIQUE NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW() 
	)
	`)

	if err != nil {
		logger.Error("Failed to create user table: " + err.Error())
		fmt.Println("Failed to create users table:", err)
		return
	}
	fmt.Println("Users table created")
	logger.Info("Users table created")
}
