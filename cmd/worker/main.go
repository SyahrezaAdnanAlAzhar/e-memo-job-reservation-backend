package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/scheduler"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/websocket"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/pkg/database"

	redisClient "github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/pkg/redis"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// INIT CONNECTION
	db := database.Connect()
	defer db.Close()

	rdb := redisClient.Connect()
	defer rdb.Close()

	authRepo := repository.NewAuthRepository(rdb)
	ticketRepo := repository.NewTicketRepository(db)
	jobRepo := repository.NewJobRepository(db)

	hub := websocket.NewHub(authRepo)
	go hub.Run()

	// CREATE INSTANCE
	ticketReorderJob := scheduler.NewTicketReorderJob(db, ticketRepo, hub)
	jobReorderJob := scheduler.NewJobReorderJob(db, jobRepo, hub)

	// INIT SCHEDULER
	jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Could not load location Asia/Jakarta: %v", err)
	}
	c := cron.New(cron.WithLocation(jakartaLocation))

	// SCHEDULES
	// CRON FORMAT: "minute hour * * Day of the Week"
	// EXAMPLE: EVERY DAY 06:00, 14:00, AND 22:00
	// c.AddJob("0 6 * * *", ticketReorderJob)
	// c.AddJob("0 14 * * *", ticketReorderJob)
	// c.AddJob("0 22 * * *", ticketReorderJob)

	// c.AddJob("0 6 * * *", jobReorderJob)
	// c.AddJob("0 14 * * *", jobReorderJob)
	// c.AddJob("0 22 * * *", jobReorderJob)

	// FOR DEVELOPMENT ONLY
	c.AddJob("0 * * * *", ticketReorderJob)
	c.AddJob("0 * * * *", jobReorderJob)

	c.Start()
	log.Println("Cron job scheduler started.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down scheduler...")
	ctx := c.Stop()
	<-ctx.Done()
	log.Println("Scheduler gracefully stopped.")
}
