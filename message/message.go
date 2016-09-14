package message

import (
	"strings"
	"unicode"
)

// Requests
const (
	// RqDump denotes a 'dump' request message.
	RqDump = "dump"

	RqFload = "fload"

	RqEject = "eject"

	RqPlay = "play"

	RqStop = "stop"

	RqEnd = "end"

	RqPos = "pos"
)

// Responses
const (
	// RsAck denotes a message with the 'ACK' response.
	RsAck = "ACK"

	// RsOhai denotes a message with the 'OHAI' response.
	RsOhai = "OHAI"

	// RsIama denotes a message with the 'IAMA' response.
	RsIama = "IAMA"

	RsFload = "FLOAD"

	//RsEject denotes a message with the 'EJECT' response.
	RsEject = "EJECT"

	RsPlay = "PLAY"

	RsStop = "STOP"

	RsEnd = "END"

	RsPos = "POS"
)

const (
	// AckOk denotes an ACK message with the 'OK' type.
	AckOk = "OK"

	// AckWhat denotes an ACK message with the 'WHAT' type.
	AckWhat = "WHAT"

	// AckFail denotes an ACK message with the 'FAIL' type.
	AckFail = "FAIL"

	// Tag that indicates a broadcast message.
	TagBroadcast = "!"
)

type Message []string

// Req constructs a request command.
func Req(tag, reqType string, params ...string) Message {
	req := Message{tag, string(reqType)}
	return append(req, params...)
}

// Res constructs a response command.
func Res(tag, resType string, params ...string) Message {
	res := Message{tag, resType}
	return append(res, params...)
}

// Ack constructs an 'ACK' response command, with the type of ACK and message,
// followed by the original request command.
func Ack(tag, ackType, msg string, origCmd Message) Message {
	resp := Message{tag, RsAck, ackType, msg}
	return append(resp, origCmd...)
}

// IsBroadcast checks the tag for the broadcast identifier.
func (m Message) IsBroadcast() bool {
	return m[0] == TagBroadcast
}

func escapeArgument(input string) string {
	return "'" + strings.Replace(input, "'", `'\''`, -1) + "'"
}

// Pack outputs the given Message as raw bytes representing a BAPS3 message.
// These bytes can be sent down a TCP connection to a BAPS3 server, providing
// they are terminated using a line-feed character.
// Note that this will panic if used on an empty Message.
func (m Message) Pack() []byte {
	outstr := m[0]
	for _, a := range m[1:] {
		// Escape arg if needed
		for _, c := range a {
			if c < unicode.MaxASCII && (unicode.IsSpace(c) || strings.ContainsRune(`'"\`, c)) {
				a = escapeArgument(a)
				break
			}
		}
		outstr += " " + a
	}
	outstr += "\n"
	return []byte(outstr)
}

// Converts a message into a string representation. Note that it doesn't escape
// the arguments, so is likely only useful for logging and debugging.
func (m Message) String() string {
	outstr := m[0]
	for _, s := range m[1:] {
		outstr += " " + s
	}
	return outstr
}
