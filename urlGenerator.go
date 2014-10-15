package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

type urlGenerator struct {
}

const (
	LETTER_LIST = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// @see http://stackoverflow.com/questions/742013/how-to-code-a-url-shortener
// I'm actually not looking for the most efficient implementation.
// I'm actually trying to make it short and hard to guess since I'm using string
// lookups in the table

func (my *urlGenerator) Hash(plaintext string) []byte {
	hash := sha256.New()

	// Add a random seed of the current platform nanoseconds prefix to the hash
	// This may be broken on platforms that don't have that resolution
	rand.Seed(int64(time.Now().Nanosecond()))
	buffer := make([]byte, 32)
	r := rand.Int63()
	written := binary.PutVarint(buffer, r)
	if written == 0 {
		log.Fatalf("Buffer not big enough for %v", r)
	}
	hash.Write(buffer)

	// Now Hash using the plaintext password
	hash.Write([]byte(plaintext))
	result := hash.Sum(nil)
	return result
}

// Turn the raw bytes into uint64s
func (my *urlGenerator) ToWords(raw []byte) []uint64 {
	var words = []uint64{}
	for i := 0; i < len(raw); i += 8 {
		end := int(math.Min(float64(len(raw)), float64(i+8)))
		words = append(words, binary.LittleEndian.Uint64(raw[i:end]))
	}
	return words
}

// Take a uint64 and then using the letter list mod the value until
// we are negative
func (my *urlGenerator) Stringify(word uint64, list string) string {
	digits := []rune{}
	lenList := uint64(len(list))
	for word > 0 {
		remainder := word % lenList
		word /= lenList
		digits = append(digits, rune(list[remainder]))
		//		fmt.Printf("%v %v %v\n", digits, remainder, word)
		//		time.Sleep(time.Second / 10)
	}
	return string(digits)
}

func (my *urlGenerator) swap(array []uint64, a, b int) {
	array[a], array[b] = array[b], array[a]
}

// Take a plain text string as input
// 1. Hash it using sha256 + a random seed based on current invocation time
// 2. Using two random 8 byte words out of the 32 bytes
// 3. Convert those into strings using modulo against a fixed list of characters
func (my *urlGenerator) Generate(plaintext string) string {
	hash := my.Hash(plaintext)
	words := my.ToWords(hash)
	currentEnd := len(words)
	// pick 2 non-repeating words at random
	// uses a Fisher-Yates shuffle with a Durstenfeld implementation
	// @see http://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	index := rand.Intn(currentEnd)
	my.swap(words, index, currentEnd-1)
	currentEnd--
	index = rand.Intn(currentEnd)
	my.swap(words, index, currentEnd-1)

	url := fmt.Sprintf("%s%s", my.Stringify(words[3], LETTER_LIST), my.Stringify(words[2], LETTER_LIST))
	return url
}
