package badamps

import (
	"bytes"
	"io"
	"regexp"
)

var badAmpFinder *regexp.Regexp = regexp.MustCompile("&[^ ;]* ")

// NextBus tends to be terrible and return non-entity ampersands.
// This function replaces any ampersand that looks like it's not
// an entity with &amp;.  This is done by finding ampersands that
// have a space before a semicolon.
// Note that this function could be considerably more efficient
// but this is how I thought to do it first and I don't want to 
// reimplement it.  Ideally, it'd be a replacement io.Reader that
// would buffer data after seeing a & and, failing to see a
// semicolon before seeing another space, would return the data
// it read into that buffer after emitting a &amp;.  This could
// be made a bit less dangerous by only reading up to the maximum
// length that an entity can be.
// For another day.  Probably not though.  This is plenty fast
// enough :D

func fixedBytes(badXML []byte) []byte {
	badAmps := badAmpFinder.FindAllIndex([]byte(badXML), -1)

	fixedXMLBuffer := &bytes.Buffer{}
	badXMLBuffer := bytes.NewReader([]byte(badXML))

	readSoFar := 0
	for _, elt := range badAmps {
		io.CopyN(fixedXMLBuffer, badXMLBuffer, int64(elt[0] - readSoFar))
		fixedXMLBuffer.Write([]byte("&amp;"))
		badXMLBuffer.ReadByte()
		readSoFar += elt[0] + 1
	}
	io.Copy(fixedXMLBuffer, badXMLBuffer)
	return fixedXMLBuffer.Bytes()
}
