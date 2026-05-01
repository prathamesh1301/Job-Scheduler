package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/redis/go-redis/v9"
)

type Job struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}


func StartWorker(ctx context.Context, rdb *redis.Client,workerName string) {
	log.Println("Worker started with name:", workerName)
	for {
		result, err := rdb.BRPop(ctx, 0, "job_queue").Result()
		log.Println("Job popped from queue:", result)
		if err != nil {
			log.Println("Error popping job from queue:", err)
			continue
		}
		jobJSON := result[1]

		var job Job
		json.Unmarshal([]byte(jobJSON), &job)

		processJob(job)
	}
}

func sendEmail(to, subject, body string) error {
    from := os.Getenv("SMTP_EMAIL")
    if from == "" {
        return fmt.Errorf("SMTP_EMAIL environment variable is not set")
    }
    password := os.Getenv("SMTP_PASSWORD")
    if password == "" {
        return fmt.Errorf("SMTP_PASSWORD environment variable is not set")
    }

    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    msg := []byte(fmt.Sprintf(
        "From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
        from, to, subject, body,
    ))

    auth := smtp.PlainAuth("", from, password, smtpHost)

    return smtp.SendMail(
        smtpHost+":"+smtpPort,
        auth,
        from,
        []string{to},
        msg,
    )
}

func processJob(job Job) {
    if err := sendEmail(job.Email, "Welcome to PrattyHub", fmt.Sprintf("Welcome %s to PrattyHub", job.Name)); err != nil {
        log.Printf("Failed to send email to %s: %v", job.Email, err)
        return
    }
    log.Printf("Successfully sent welcome email to %s", job.Email)
}