package baps3

import (
	"errors"
	"sort"
)

// Feature is the type for known feature flags.
type Feature int

const (
	/* Feature constants.
	 *
	 * When adding to this, also add the string equivalent to ftStrings.
	 */

	// FtUnknown represents an unknown feature.
	FtUnknown Feature = iota
	// FtFileLoad represents the FileLoad standard feature.
	FtFileLoad
	// FtPlayStop represents the PlayStop standard feature.
	FtPlayStop
	// FtSeek represents the Seek standard feature.
	FtSeek
	// FtEnd represents the End standard feature.
	FtEnd
	// FtTimeReport represents the TimeReport standard feature.
	FtTimeReport
	// FtPlaylist represents the Playlist standard feature.
	FtPlaylist
	// FtPlaylistAutoAdvance represents the Playlist.AutoAdvance feature.
	FtPlaylistAutoAdvance
	// FtPlaylistTextItems represents the Playlist.TextItems feature.
	FtPlaylistTextItems
)

// Yes, a global variable.
// Go can't handle constant arrays.
var ftStrings = []string{
	"<UNKNOWN FEATURE>",    // FtUnknown
	"FileLoad",             // FtFileLoad
	"PlayStop",             // FtPlayStop
	"Seek",                 // FtSeek
	"End",                  // FtEnd
	"TimeReport",           // FtTimeReport
	"Playlist",             // FtPlaylist
	"Playlist.AutoAdvance", // FtPlaylistAutoAdvance
	"Playlist.TextItems",   // FtPlaylistTextItems
}

// IsUnknown returns whether word represents a feature unknown to Bifrost.
func (word Feature) IsUnknown() bool {
	return word == FtUnknown
}

func (word Feature) String() string {
	return ftStrings[int(word)]
}

// LookupFeature finds the equivalent Feature for a string.
// If the message word is not known to Bifrost, it will return FtUnknown.
func LookupFeature(word string) Feature {
	// This is O(n) on the size of ftStrings, which is unfortunate, but
	// probably ok.
	for i, str := range ftStrings {
		if str == word {
			return Feature(i)
		}
	}

	return FtUnknown
}

// FeatureSet is a set of features. Go figure.
type FeatureSet map[Feature]struct{}

// FeatureSetFromMsg returns a populated FeatureSet from a RsFeatures
func FeatureSetFromMsg(msg *Message) (fs FeatureSet, err error) {
	if msg.Word() != RsFeatures {
		err = errors.New("Message is not a FEATURES message")
		return
	}
	fs = make(FeatureSet)
	for _, featurestr := range msg.Args() {
		f := LookupFeature(featurestr)
		if f == FtUnknown {
			err = errors.New("Unknown feature: " + featurestr)
		}
		fs[f] = struct{}{}
	}
	return
}

// AddFeature adds a feature to a featureset
// The given featureset is returned to enable chaining
func (fs FeatureSet) AddFeature(feat Feature) FeatureSet {
	fs[feat] = struct{}{}
	return fs
}

// DelFeature deletes a feature from a featureset
// The given featureset is returned to enable chaining
func (fs FeatureSet) DelFeature(feat Feature) FeatureSet {
	delete(fs, feat)
	return fs
}

// ToMessage converts a featureset into a RsFeatures message
func (fs FeatureSet) ToMessage() (msg *Message) {
	// Sort features alphabetically to make output deterministic
	// Otherwise tests are a pain in the arse
	var featstrings []string
	for f := range fs {
		featstrings = append(featstrings, f.String())
	}
	sort.Strings(featstrings)
	msg = NewMessage(RsFeatures)
	for _, featstring := range featstrings {
		msg.AddArg(featstring)
	}
	return
}
