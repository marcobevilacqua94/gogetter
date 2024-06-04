package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/couchbase/gocb/v2"
)

// const rootPEM = `
// -----BEGIN CERTIFICATE-----
// MIIDDDCCAfSgAwIBAgIIF9VKJ87vtsQwDQYJKoZIhvcNAQELBQAwJDEiMCAGA1UE
// AxMZQ291Y2hiYXNlIFNlcnZlciBjNmE2MDVhNzAeFw0xMzAxMDEwMDAwMDBaFw00
// OTEyMzEyMzU5NTlaMCQxIjAgBgNVBAMTGUNvdWNoYmFzZSBTZXJ2ZXIgYzZhNjA1
// YTcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCaL24ymdVip/7tnS6a
// 7qHkC0bA87Cv3Aapnw6aWVPWZLm/oxUi3/6JCrqZBgJYBV/YiJTkdMQ6AU/6mLDS
// dn/5ohI4M+5QZgwOF1HYREMZmJ/3K39w4EwLmHTRKqtkft7RhZe0r1G01pJT6RPc
// pJGvoqn67KtibqIAy683VmA2XSgmS8MYvsW8f3U7rtMG2vqlKMBvzYkJZuCU9bQ3
// TE80Vkypg2XHGOFpYJLCxm2FuRVQ1WQpTnGP+vjN5Aed6NCNzUFChdvTuUVi4pZL
// uwaO7026em4N/T/Zn27SkNZ/SvwK9JVw/dgIQCV26sF5O/GENezm93c10GqZw2OP
// k5ffAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBhjAPBgNVHRMBAf8EBTADAQH/MB0G
// A1UdDgQWBBR0HQjyRtmZQTE40UJd89ot/O4sHzANBgkqhkiG9w0BAQsFAAOCAQEA
// B4D7FROAimUQhnTLeShFVNICa9z/wYw3i4wH8UgKhFO9kEsVTg2+qFkAVy2o+xCf
// sRhH9kYHsLidGc6s38qUpF/gsH3UoTTFO3j/WS+2G4kZ6yQPiQy0DaoR0aVdVx1a
// TjdzwgGCu4NJFAcvW0xOqPVyZxGTOrLz8tiNcdykw7CICO4/cmTzQHkidGh4qVyh
// OguDXwNefzj6LkTxOpZODKPmxT5Udw0DUBBWB16ebOItUUZtd+JqYziLsPAAC18a
// FCp+9KQg/ZtYA/O/OOBik6H+PkbZPSQUEBDFhHmJvYHzeMaEFeVFGMHGD7wwLPqg
// 9cUMo085986p5luEIIrA1w==
// -----END CERTIFICATE-----
// `

func main() {
	//roots := x509.NewCertPool()
	//roots.AppendCertsFromPEM([]byte(rootPEM))

	// Uncomment following line to enable logging
	// gocb.SetLogger(gocb.VerboseStdioLogger())

	username := "Administrator"
	password := "password"
	ip := "127.0.0.1"
	bucketName := "sample"
	scopeName := "_default"
	collectionName := "_default"
	concurrency := "16"

	args := os.Args

	username = args[1]
	password = args[2]
	ip = args[3]
	bucketName = args[4]
	scopeName = args[5]
	collectionName = args[6]
	concurrency = args[7]

	// Update this to your cluster details
	connectionString := "couchbase://" + ip

	gocb.VerboseStdioLogger()

	// For a secure cluster connection, use `couchbases://<your-cluster-ip>` instead.
	cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
		// SecurityConfig: gocb.SecurityConfig{
		// 	TLSRootCAs: roots,
		// },
	})
	if err != nil {
		log.Fatal(err)
	}

	bucket := cluster.Bucket(bucketName)

	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	scope := bucket.Scope(scopeName)
	collection := scope.Collection(collectionName)

	var wg sync.WaitGroup

	numOps, err := strconv.Atoi(concurrency)

	// Genera un numero di telefono mobile dell'Inghilterra
	generateUKMobileNumber := func() string {
		ukMobilePrefixes := []string{"071", "073", "074", "075", "077", "078", "079"}
		prefix := ukMobilePrefixes[randomInt(len(ukMobilePrefixes))]
		remainingDigits := randomDigits(8)
		phoneNumber := prefix + remainingDigits
		return "+44" + phoneNumber
	}

	// Genera un numero di telefono mobile dell'Italia
	generateITMobileNumber := func() string {
		itMobilePrefixes := []string{"32", "33", "34", "36", "37", "38", "39"}
		prefix := itMobilePrefixes[randomInt(len(itMobilePrefixes))]
		remainingDigits := randomDigits(8)
		phoneNumber := prefix + remainingDigits
		return "+39" + phoneNumber
	}

	for i := 0; i < numOps; i++ {
		if randomInt(1) < 5 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for true {
					collection.Get(generateUKMobileNumber(), nil)
				}
			}()
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for true {
					collection.Get(generateITMobileNumber(), nil)
				}
			}()
		}
	}

	wg.Wait()

}

// randomInt returns a random integer from 0 to n-1.
func randomInt(n int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))
	return int(num.Int64())
}

// randomDigits generates a string of random digits of the given length.
func randomDigits(length int) string {
	var digits strings.Builder
	for i := 0; i < length; i++ {
		digit := randomInt(10)
		digits.WriteString(fmt.Sprint(digit))
	}
	return digits.String()
}
