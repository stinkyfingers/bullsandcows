package game

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Game struct {
	ID           string           `json:"id"`
	Created      time.Time        `json:"created,omitempty"`
	Session      *session.Session `json:"-"`
	ConnectionID string           `json:"connectionId,omitempty"`
	Rounds       []Round          `json:"rounds,omitempty"`
	Code         Sequence         `json:"code,omitempty"`
	Concluded    bool             `json:"concluded,omitempty"`
	Won          bool             `json:"won,omitempty"`
	TTL          int64            `json:"ttl"`
}

type Sequence [4]Peg

type Peg struct {
	Color string `json:"color"`
	Int   int    `json:"int"`
}

type Round struct {
	Sequence Sequence `json:"sequence,omitempty"`
	Bulls    int      `json:"bulls,omitempty"`
	Cows     int      `json:"cows,omitempty"`
}

var (
	Red    = Peg{"red", 1}
	Blue   = Peg{"blue", 2}
	Green  = Peg{"green", 3}
	Yellow = Peg{"yellow", 4}
	Purple = Peg{"purple", 5}
	Orange = Peg{"orange", 6}
)

var tableName = "bullsandcows"

const maxRounds = 12

func NewGame(sess *session.Session, connectionID string) *Game {
	g := &Game{
		ID:           createID(),
		Created:      time.Now(),
		Session:      sess,
		ConnectionID: connectionID,
		Rounds:       []Round{},
		Code:         *randomCode(),
		TTL:          time.Now().Add(time.Hour * 24).Unix(),
	}
	return g
}

func NewSequence() *Sequence {
	s := Sequence([4]Peg{})
	return &s
}

func (g *Game) PlayRound(s Sequence) error {
	err := g.validateSequence(&s)
	if err != nil {
		return err
	}

	cowGroup := make(map[Peg]int)
	recheckPegs := make(map[Peg]int)
	var bulls, cows int

	// bulls
	for i, peg := range s {
		if peg.Int == g.Code[i].Int {
			bulls++
			continue
		}
		cowGroup[g.Code[i]]++
		recheckPegs[peg]++
	}

	// cows
	for key, quantity := range recheckPegs {
		if cowQuantity, ok := cowGroup[key]; !ok {
			continue
		} else {
			if quantity < cowQuantity || quantity == cowQuantity {
				cows += quantity
			}
			if quantity > cowQuantity {
				cows += cowQuantity
			}
		}
	}

	g.Rounds = append(g.Rounds, Round{
		Sequence: s,
		Bulls:    bulls,
		Cows:     cows,
	})
	g.Won = bulls == 4
	g.Concluded = (len(g.Rounds) == maxRounds || g.Won)
	return g.Put()
}

func (g *Game) validateSequence(sequence *Sequence) error {
	intMap := map[int]Peg{
		1: Red, 2: Blue, 3: Green, 4: Yellow, 5: Purple, 6: Orange,
	}
	strMap := map[string]Peg{
		"red": Red, "blue": Blue, "green": Green, "yellow": Yellow, "purple": Purple, "orange": Orange,
	}
	for i := range sequence {
		if sequence == nil || (sequence[i].Int == 0 && sequence[i].Color == "") {
			return errors.New("incomplete sequence")
		}
		if sequence[i].Int == 0 {
			sequence[i] = strMap[sequence[i].Color]
		}
		if sequence[i].Color == "" {
			sequence[i] = intMap[sequence[i].Int]
		}
	}
	return nil
}

func randomCode() *Sequence {
	sequence := NewSequence()
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		sequence[i] = getPegs()[rand.Intn(6)]
	}
	return sequence
}

func getPegs() []Peg {
	return []Peg{Red, Blue, Green, Yellow, Purple, Orange}
}

func createID() string {
	// TODO - assure not in DB
	rand.Seed(time.Now().UnixNano())
	idInt := rand.Intn(9999)
	return strconv.Itoa(idInt)
}

/* DB */

func (g *Game) Put() error {
	item, err := dynamodbattribute.MarshalMap(g)
	if err != nil {
		return err
	}

	_, err = dynamodb.New(g.Session).PutItem(&dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	return err
}

func Get(sess *session.Session, gameID string) (Game, error) {
	var g Game
	id, err := dynamodbattribute.Marshal(gameID)
	if err != nil {
		return g, err
	}

	resp, err := dynamodb.New(sess).GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       map[string]*dynamodb.AttributeValue{"id": id},
	})
	if err != nil {
		return g, err
	}
	if resp == nil {
		return g, errors.New("nil response")
	}
	if resp.Item == nil {
		return g, errors.New("game not found")
	}
	err = dynamodbattribute.UnmarshalMap(resp.Item, &g)
	if err != nil {
		return g, err
	}
	g.Session = sess
	return g, nil
}

func (g *Game) Delete() error {
	id, err := dynamodbattribute.Marshal(g.ID)
	if err != nil {
		return err
	}
	_, err = dynamodb.New(g.Session).DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key:       map[string]*dynamodb.AttributeValue{"id": id},
	})
	return err
}
