package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/adjust/rmq/v3"
)

var (
	redisConnection rmq.Connection
	tweetsQueue rmq.Queue
)

type RedisCredentials struct {
    Protocol string
    Host     string
}

func loginToRedis(credentials RedisCredentials){
	var err error
	redisConnection, err = rmq.OpenConnection("pbr-th", credentials.Protocol, credentials.Host, 1, nil)
	if err != nil {
        log.Warn("Error getting Redis Client")
        log.Error(err)
    }
	log.Info("Logged in into Redis")
}

func openQueue(queueName string){
	var err error
	tweetsQueue, err = redisConnection.OpenQueue(queueName)
	if err != nil {
		log.Warn("Error getting Redis Queue")
        log.Error(err)
	}
	log.Info("Opened Redis Queue")
}

func publishBytesToQueue(message []byte)(error){
	return tweetsQueue.PublishBytes(message)
}