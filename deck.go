package decks

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type Card struct {
	Suit string
	Rank string
}

type Options struct {
	sorting      func(i, j int) bool
	shuffle      bool
	numJokers    int
	filterRanks  []string
	composedDeck int
}

type DeckOption func(*Options)

func WithSorting(less func(i, j int) bool) DeckOption {
	return func(o *Options) {
		o.sorting = less
	}
}

func WithShuffle() DeckOption {
	return func(o *Options) {
		o.shuffle = true
	}
}

func WithJokers(numJokers int) DeckOption {
	return func(o *Options) {
		o.numJokers = numJokers
	}
}

func WithFilterRanks(filter []string) DeckOption {
	return func(o *Options) {
		o.filterRanks = filter
	}
}

func WithComposedDeck(numDecks int) DeckOption {
	return func(o *Options) {
		o.composedDeck = numDecks
	}
}

type Deck struct {
	Cards []Card
}

func New(options ...DeckOption) []Card {
	var cards []Card
	suits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	opt := &Options{
		sorting:      nil,
		shuffle:      false,
		numJokers:    0,
		filterRanks:  nil,
		composedDeck: 1,
	}
	for _, option := range options {
		option(opt)
	}

	for _, suit := range suits {
		for rank := 1; rank <= 13; rank++ {
			if rank == 11 {
				cards = append(cards, Card{Rank: "J", Suit: suit})
			} else if rank == 12 {
				cards = append(cards, Card{Rank: "Q", Suit: suit})
			} else if rank == 13 {
				cards = append(cards, Card{Rank: "K", Suit: suit})
			} else {
				cards = append(cards, Card{Rank: strconv.Itoa(rank), Suit: suit})
			}
		}
	}
	cards = applyOptions(cards, opt)
	return cards
}

func applyOptions(cards []Card, opt *Options) []Card {
	if opt.sorting != nil {
		sort.Slice(cards, func(i, j int) bool {
			return opt.sorting(i, j)
		})
	}

	if opt.shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
	}

	for i := 0; i < opt.numJokers; i++ {
		cards = append(cards, Card{Rank: "Joker", Suit: "Joker"})
	}

	if opt.filterRanks != nil {
		filteredCards := make([]Card, 0, len(cards))
		for _, card := range cards {
			if !contains(opt.filterRanks, card.Rank) {
				filteredCards = append(filteredCards, card)
			}
		}
		cards = filteredCards
	}

	return cards

}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// func main() {
// 	New(
// 		WithShuffle(),
// 		WithJokers(2),
// 		WithFilterRanks([]string{"2", "3"}), // Filter out 2s and 3s
// 		WithComposedDeck(3),                 // Create a deck composed of 3 standard decks
// 	)

// }
