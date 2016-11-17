package responses

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/perthgophers/puddle/messagerouter"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"
)

// BokBok is a smart chat bot
type BokBok struct {
	db        *bolt.DB
	Channel   string
	prefixLen int
}

// NewBokBok makes a BoKBok and opens the bolt database
func NewBokBok(prefixLen int) *BokBok {
	var err error = nil
	bkbk := new(BokBok)
	bkbk.db, err = bolt.Open("./markovchains.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	bkbk.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("all"))
		return nil
	})

	bkbk.prefixLen = prefixLen

	return bkbk
}

// RespondHere allows the user to apply a channel to respond to
// Triggered by !bokbok <channel>
func (bkbk *BokBok) RespondHere(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	bkbk.Channel = cr.Message.Channel

	w.Write("Channel has been set")

	return nil
}

// Processes Message builds or amends the Markov Chain for "all" and the individual user
func (bkbk *BokBok) ProcessMessage(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	bkbk.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(cr.Username))
		return nil
	})
	ch := bkbk.Chain(cr.Username)
	ch.BuildFromString(cr.Text)

	allch := bkbk.Chain("all")
	allch.BuildFromString(cr.Text)

	bkbk.SaveChains(ch, allch)
	return nil
}

//Respond reponds to message
func (bkbk *BokBok) Respond(w messagerouter.ResponseWriter) {
	allch := bkbk.Chain("all")
	w.Write(fmt.Sprintf("@%s: %s", cr.Username, allch.Generate()))
}

// MaybeRespond might respond, or it might not, for top kek
func (bkbk *BokBok) MaybeRespond(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	if cr.Message.Channel == cr.Message.User {
		bkbk.Respond(w)
		return nil
	}
	if cr.Message.User != "0" && cr.Message.Channel != bkbk.Channel {
		return nil
	}
	if bkbk.YesNo() {
		bkbk.Respond(w)
	}

	return nil
}

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	Username  string
	Chain     map[string][]string
	PrefixLen int
}

// Marshal encodes a chain to json.
func (c *Chain) Marshal() ([]byte, error) {
	enc, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// Decode decodes json to Chain
func (c *Chain) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return err
	}
	return nil
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func (bkbk *BokBok) NewChain(username string) *Chain {
	c := new(Chain)
	c.PrefixLen = bkbk.prefixLen
	c.Username = username
	c.Chain = make(map[string][]string)

	return c
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.PrefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		c.Chain[key] = append(c.Chain[key], s)
		p.Shift(s)
	}
}

// BuildFromString ...
func (c *Chain) BuildFromString(s string) {
	p := make(Prefix, c.PrefixLen)
	for _, v := range strings.Split(s, " ") {
		key := p.String()
		c.Chain[key] = append(c.Chain[key], v)
		p.Shift(v)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate() string {
	p := make(Prefix, c.PrefixLen)
	var words []string
	for {
		choices := c.Chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

// Chain retreives chain from bkbk.db. Creates username bucket if it doesn't exist
func (bkbk *BokBok) Chain(username string) *Chain {
	var ch *Chain = bkbk.NewChain(username)
	var data []byte
	bkbk.db.View(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(username))

		b := tx.Bucket([]byte(username))
		log.Println("BUCKET", b)
		data = b.Get([]byte("chain"))
		return nil
	})
	if data != nil {
		ch.Unmarshal(data)
	} else {
		bkbk.NewChain(username)
	}

	return ch
}

func (bkbk *BokBok) YesNo() bool {
	n := rand.Intn(6-1) + 1
	fmt.Println(n)
	return n == 3
}

// SaveChains saves chains to the database
func (bkbk *BokBok) SaveChains(chains ...*Chain) error {
	err := bkbk.db.Update(func(tx *bolt.Tx) error {
		for _, ch := range chains {
			b := tx.Bucket([]byte(ch.Username))
			if b == nil {
				return fmt.Errorf("Can't retrieve bucket for %s", ch.Username)
			}
			j, err := ch.Marshal()
			if err != nil {
				return fmt.Errorf("Can't encode chain: %s", err)
			}
			b.Put([]byte("chain"), j)
		}

		return nil
	})
	return err
}

func init() {
	rand.Seed(time.Now().Unix())
	bkbk := NewBokBok(2)
	Handle("*", bkbk.ProcessMessage)
	Handle("*", bkbk.MaybeRespond)
	Handle("!bokbok", bkbk.RespondHere)
}
