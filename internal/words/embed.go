package words

import _ "embed"

//go:embed english.json
var EnglishJSON []byte

//go:embed english_1k.json
var English1kJSON []byte

//go:embed quotes.json
var QuotesJSON []byte
