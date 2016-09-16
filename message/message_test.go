package message_test

import (
	"reflect"
	"testing"

	msg "github.com/UniversityRadioYork/bifrost-go/message"
)

// It feels like there should be more tests here, but since Message is
// essentially just a []string, why bother?

func TestPack(t *testing.T) {
	cases := []struct {
		msg     msg.Message
		want    []byte
		wantstr string
	}{
		// Res helper func, Unescaped command
		{
			msg.Res("uuid", msg.RsFload, "/this/is/a/file"),
			[]byte("uuid FLOAD /this/is/a/file\n"),
			"uuid FLOAD /this/is/a/file",
		},
		// Req helper func, Backslashes
		{
			msg.Req("uuid", msg.RqFload, `C:\silly\windows\is\silly`),
			[]byte(`uuid fload 'C:\silly\windows\is\silly'` + "\n"),
			"uuid fload C:\\silly\\windows\\is\\silly",
		},
		// Spaces
		{
			msg.Ack("uuid", msg.AckOk, "/home/the donald/01 The Nightfly.mp3"),
			[]byte("uuid ACK OK '/home/the donald/01 The Nightfly.mp3'\n"),
			"uuid ACK OK /home/the donald/01 The Nightfly.mp3",
		},
		// Single quotes
		{
			msg.Message{msg.RsOhai, "a'bar'b"},
			[]byte(`OHAI 'a'\''bar'\''b'` + "\n"),
			`OHAI a'bar'b`,
		},
		// Double quotes
		{
			msg.Message{msg.RsOhai, `a"bar"b`},
			[]byte(`OHAI 'a"bar"b'` + "\n"),
			`OHAI a"bar"b`,
		},
		// Single word (shouldn't ever be used)
		{
			msg.Message{msg.RsOhai},
			[]byte("OHAI\n"),
			"OHAI",
		},
	}

	for _, c := range cases {
		got := c.msg.Pack()
		gotstr := c.msg.String()
		if !reflect.DeepEqual(c.want, got) {
			t.Errorf("Message.Pack(%q) == %q, want %q", c.msg, got, c.want)
		}
		if gotstr != c.wantstr {
			t.Errorf("Message.String(%q) == %q, want %q", c.msg, gotstr, c.wantstr)
		}
	}
}
