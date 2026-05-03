package main

import (
	"auth/internals/redis"
	"auth/internals/store"
	"context"

	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"
)

type Job struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Retries   int
	LastError string `json:"error"`
}

var MaxRetry int = 3

func StartWorker(ctx context.Context, rdb *redis.Redis, workerName string, store *store.Store) {
	log.Println("Worker started with name:", workerName)
	for {
		result, err := rdb.Client.BRPop(ctx, 0, workerName).Result()
		log.Println("Job popped from queue:", result)
		if err != nil {
			log.Println("Error popping job from queue:", err)
			continue
		}
		jobJSON := result[1]

		var job Job
		json.Unmarshal([]byte(jobJSON), &job)
		host, _ := os.Hostname()
		log.Println("Worker:", host, "processing job:", job.Name)

		if err := processJob(ctx, job, store,rdb); err != nil {
			job.Retries++
			job.LastError = err.Error()
			if job.Retries <= MaxRetry {
				delay := time.Second * time.Duration(1<<job.Retries)

				go func(j Job, d time.Duration) {
					time.Sleep(d)

					jobJSON, err := json.Marshal(j)
					if err != nil {
						log.Println("Marshal error:", err)
						return
					}

					if err := rdb.Client.LPush(ctx, workerName, jobJSON).Err(); err != nil {
						log.Println("Requeue failed:", err)
					}
				}(job, delay)
			} else {
				jobJSON, _ := json.Marshal(job)
				rdb.Client.LPush(ctx, "failed_jobs", jobJSON)
			}
		}
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

func processJob(ctx context.Context, job Job, store *store.Store,rdb *redis.Redis) error {
	idempotencyKey := job.Email + job.Name + "welcome"
	cacheKey := "cache:" + idempotencyKey

    exists, err := rdb.Client.Exists(ctx, cacheKey).Result()
    if err != nil {
        log.Println("Redis error:", err)
    }

    if exists == 1 {
        log.Println("Cache hit:Job already processed skipping", job.Email)
        return nil
    }

	dbexists, err := store.Idempotency.Exists(context.Background(), idempotencyKey)
	if err != nil {
		return err
	}
	if dbexists {
		log.Println("Job already processed skipping", job.Email)
		return nil
	}
	if err := sendEmail(job.Email, "Welcome to PrattyHub", fmt.Sprintf("Welcome %s to PrattyHub", job.Name)); err != nil {
		log.Printf("Failed to send email to %s: %v", job.Email, err)
		return err
	}
	store.Idempotency.Create(ctx, idempotencyKey)
	rdb.Client.SetNX(ctx, cacheKey, "1", 24*time.Hour)
	log.Printf("Successfully sent welcome email to %s", job.Email)
	return nil
}
