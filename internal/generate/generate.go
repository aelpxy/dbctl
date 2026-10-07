// Package generate creates random names and passwords for new databases.
package generate

import (
	"crypto/rand"
	mrand "math/rand/v2"
	"strings"
)

const adjectives = `autumn hidden bitter misty silent empty dry dark summer icy delicate quiet white cool
spring winter patient twilight dawn crimson wispy weathered blue billowing broken cold damp falling
frosty green long late lingering bold little morning muddy old red rough still small sparkling
throbbing shy wandering withered wild black young holy solitary fragrant aged snowy proud floral
restless divine polished ancient purple lively nameless`

const nouns = `waterfall river breeze moon rain wind sea morning snow lake sunset pine shadow leaf dawn
glitter forest hill cloud meadow sun glade bird brook butterfly bush dew dust field fire flower
firefly feather grass haze mountain night pond darkness snowflake silence sound sky shape surf
thunder violet water wildflower wave resonance wood dream cherry tree fog frost voice paper frog
smoke star`

// Name returns a random four-word name such as "misty-river-bold-pine".
func Name() string {
	adj := strings.Fields(adjectives)
	noun := strings.Fields(nouns)

	return strings.Join([]string{pick(adj), pick(noun), pick(adj), pick(noun)}, "-")
}

// Password returns a random password with 128 bits of entropy.
func Password() string {
	return rand.Text()
}

func pick(words []string) string {
	return words[mrand.IntN(len(words))]
}
