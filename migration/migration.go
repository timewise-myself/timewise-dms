package main

import (
	"dbms/config"
	"dbms/database"
	"github.com/spf13/viper"
	"github.com/timewise-team/timewise-models/models"
	"log"
)

func main() {
	viper.AddConfigPath("..")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		return
	}

	viper.AutomaticEnv()
	cfg := &config.Config{
		ServerPort: viper.GetString("WEB.PORT"),
		DBUser:     viper.GetString("DB.USERNAME"),
		DBPassword: viper.GetString("DB.PASSWORD"),
		DBName:     viper.GetString("DB.NAME"),
		DBHost:     viper.GetString("DB.HOST"),
		DBPort:     viper.GetString("DB.PORT"),
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
		return
	}

	// Migrate the schema
	err = db.AutoMigrate(
		//&models.TwUser{},
		//&models.TwUserEmail{},
		//&models.TwWorkspace{},
		//&models.TwWorkspaceUser{},
		//&models.TwBoardColumn{},
		//&models.TwSchedule{},
		//&models.TwScheduleParticipant{},
		//&models.TwScheduleLog{},
		//&models.TwComment{},
		//&models.TwRecurrenceException{},
		&models.TwReminder{},
		//&models.TwWorkspaceLog{},
		//&models.TwNotificationSettings{},
		//&models.TwNotifications{},
		//&models.TwDocument{},
	)
	if err != nil {
		log.Fatalf("Could not migrate schema: %v", err)
		return
	} else {
		log.Println("Migration success")
	}
}
