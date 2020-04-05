package main

import (
	"flag"
	"github.com/go-redis/redis/v7"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

type Node struct {
	url string
	id string
}

type JobAssignment struct {
	id string
	jobs []string
}

type State string

const MASTER State = "master"
const WORKER State = "worker"

type Job struct {
	id string
	name string
	// other job related stuff
}

const electionKey = "__election_key__"
const topologyKey = "__topology__"
const jobAssignmentKey = "__job_assignment__"
const jobPrefixKey = "job:"

var redisClient *redis.Client

func initLogger() {
	logger, _ = zap.NewProduction()
	defer func(){ _ = logger.Sync()}()
	sugar = logger.Sugar()
}

func connectRedis() {
	redisClient= redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		sugar.Fatalf("error while creating redis client : %v", err.Error())
	}
	sugar.Info("successfully connected to redis")
}

func closeRedisConnection() {
	if redisClient != nil {
		sugar.Info("disconnecting from redis")
		err := redisClient.Close()
		if err != nil {
			sugar.Fatalf("error while closing redis connection : %v", err.Error())
		}
	}
}

// this generates uniq id for a node
func generateID() string { return ksuid.New().String() }

// this runs the election process
func startElection(id string) (chan State, error) { return nil, nil }

// this fetches node topology from redis
func fetchTopology() ([]Node, error) { return nil, nil }

// this does health check and returns true if healthy , false otherwise
func healthCheck(node Node) bool { return false }

// fetch job assignments from redis for the current node, id is current nodes' id
func fetchJobAssignment(id string) JobAssignment { return JobAssignment{} }

// fetches assigned jobs' details from redis
func fetchJobDetails(jobs []string) ([]Job, error) { return nil, nil }

// this process a job
func processJob(job Job) {
	sugar.Infof("processing job : %s:%s", job.id, job.name)
}

func main() {
	initLogger()

	port := flag.String("port", "", "http port, ex -port=:8000")
	flag.Parse()

	if *port == "" {
		sugar.Error("missing command line arg")
		flag.Usage()
		os.Exit(-1)
	}

	id := generateID()
	sugar.Infof("starting worker: id - %v", id)
	connectRedis()
	defer closeRedisConnection()

	go startServer(*port)

	signalCh := make(chan os.Signal, 1)
	exitProcess := make(chan struct{})
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func () {
		<- signalCh
		// perform cleanup
		sugar.Info("received os interrupt, exiting")
		exitProcess <- struct{}{}
		return
	}()
	<-exitProcess
}

func startServer(port string) {
	http.HandleFunc("/health", healthEndpointHandler)
	sugar.Infof("starting server at : %s", port)
	http.ListenAndServe(port, nil)
}

func healthEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}