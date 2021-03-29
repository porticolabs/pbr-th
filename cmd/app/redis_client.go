package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/adjust/rmq/v3"
)

var (
	redisConnection rmq.Connection // Redis service acive connection
	tweetsQueue rmq.Queue // Redis queue to publish the messages to
)

type RedisCredentials struct {
    Protocol string // Redis service protocol
    Host     string // Redis service host: hostname:port
}

// Uses the connection data provided to login
//  into Redis service, and store the connection
//  in a local global variable
func loginToRedis(credentials RedisCredentials){
	var err error
	redisConnection, err = rmq.OpenConnection("pbr-th", credentials.Protocol, credentials.Host, 1, nil)
	if err != nil {
        log.Warn("Error getting Redis Client")
        log.Error(err)
    }
	log.Info("Logged in into Redis")
}

// Uses the Redis connection to prepare the 
//   Queue to be used for publishing
//   curated tweets
func openQueue(queueName string){
	var err error
	tweetsQueue, err = redisConnection.OpenQueue(queueName)
	if err != nil {
		log.Warn("Error getting Redis Queue")
        log.Error(err)
	}
	log.Info("Opened Redis Queue")
}

// Publish a tweet in json format on the Redis queue
func publishBytesToQueue(message []byte)(error){
	return tweetsQueue.PublishBytes(message)
}