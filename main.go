package main

import (
	"context"
	"fmt"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"os"
	"strconv"
	"time"

	immuclient "github.com/codenotary/immudb/pkg/client"
)

const pool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const username = "immudb"
const password = "immudb"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Structure of key-value pairs
type kventry struct {
	key string
	val string
}

// Handles random generation of strings
// Returns a string of length = l
func randStr(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = pool[seededRand.Intn(len(pool))]
	}

	return string(bytes)
}

// Populates a slice of kventries
// Returns key-value list/slice
func genvals(n int) []kventry {
	// Make a slice of kventries with space for n entries
	kvl := make([]kventry, n)
	for i := 0; i < n; i++ {
		kvl[i].key = randStr(9)
		kvl[i].val = randStr(42)
	}

	return kvl
}

// Inserts values into database using VerifiedSet
func verifiedInsert(ctx context.Context, client immuclient.ImmuClient, kvl []kventry ) {
	for i := 0; i < len(kvl); i++ {
		tx, err := client.VerifiedSet(ctx, []byte(kvl[i].key), []byte(kvl[i].val))
		if err != nil {
			log.Fatal("Error inserting: " + err.Error())
		}
		fmt.Printf("Inserting key: %s; val: %s; at index [%d]\n", kvl[i].key, kvl[i].val, tx.Id)
	}
}

// connect to the immudb database
func connect() (immuclient.ImmuClient, context.Context) {
	client, err := immuclient.NewImmuClient(immuclient.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	lr, err := client.Login(ctx, []byte(username), []byte(password))
	if err != nil {
		log.Fatal(err)
	}

	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	return client, ctx
}

// Validate command-line argument input
func parseArgs(argv []string) int {
	if len(argv) > 2 {
		log.Fatal("Too many arguments!")
	}

	i, err := strconv.Atoi(argv[1])
	if i <= 0 {
		log.Fatal("Invalid input value\n")
	}

	if err != nil {
		log.Fatal("Invalid input value" + err.Error())
	}


	return i
}

func main() {
	// Get number of entries to insert from user
	n := parseArgs(os.Args)
	fmt.Printf("Number of entries requested: %d\n", n)

	// Create list of key-value pairs of length = n
	kvl := genvals(n)

	// initiate connection to the database and grab client and context
	client, ctx := connect()

	// now insert the list into immudb
	start := time.Now()
	verifiedInsert(ctx, client, kvl)
	elapsed := time.Since(start)

	fmt.Printf("%d Entries were inserted within %.2f seconds\n", n, elapsed.Seconds())
}
