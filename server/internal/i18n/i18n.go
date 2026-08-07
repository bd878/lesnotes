package i18n

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"log/slog"
	"strings"
)

type translationsData struct {
	En map[string]string  `json:"en,omitempty"`
	Ru map[string]string  `json:"ru,omitempty"`
	De map[string]string  `json:"de,omitempty"`
	Fr map[string]string  `json:"fr,omitempty"`
}

type declinationsData struct {
	En map[string][]string `json:"en,omitempty"`
	Ru map[string][]string `json:"ru,omitempty"`
	De map[string][]string `json:"de,omitempty"`
	Fr map[string][]string `json:"fr,omitempty"`
}

var emptyDecl []string = make([]string, 3)

func (t translationsData) Get(code LangCode, key string) string {
	value := reflect.ValueOf(t)
	fieldValue := value.FieldByName(code.String())
	if !fieldValue.IsValid() {
		slog.Error("field value is invalid", slog.String("value", fieldValue.String()))
		return ""
	}

	dict, ok := fieldValue.Interface().(map[string]string)
	if !ok {
		slog.Error("translations: not ok", slog.String("context", "translationsData.Get"))
		return ""
	}

	text, ok := dict[key]
	if !ok {
		return ""
	}

	return text
}

func (d declinationsData) Get(code LangCode, key string) []string {
	value := reflect.ValueOf(d)
	fieldValue := value.FieldByName(code.String())
	if !fieldValue.IsValid() {
		slog.Error("declinations: field value is invalid", slog.String("value", fieldValue.String()))
		return emptyDecl
	}

	data, ok := fieldValue.Interface().(map[string][]string)
	if !ok {
		slog.Error("translations: not ok", slog.String("context", "declinationsData.Get"))
		return emptyDecl
	}

	texts, ok := data[key]
	if !ok {
		slog.Error("translations: cannot convert field value to map[string][]string", slog.String("context", "declinationsData.Get"))
		return emptyDecl
	}

	return texts
}

//go:embed texts.json
var textsFile []byte
var texts translationsData

//go:embed decls.json
var declsFile []byte
var decls declinationsData

func init() {
	if err := json.Unmarshal(textsFile, &texts); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(declsFile, &decls); err != nil {
		panic(err)
	}
}

type Translator interface {
	Code() LangCode
	Text(key string) string
	Decl(key string) []string
}

type LangTranslator string
func (t LangTranslator) Language(code Translator) string {
	return code.Text(string(t))
}

type SameText string
func (t SameText) Language(code Translator) string {
	return string(t)
}

type Text interface {
	Language(code Translator) string
}

type Decl interface {
	Language(code Translator) []string
}

type Translation struct {
	Ru string `json:"ru,omitempty"`
	En string `json:"en,omitempty"`
	De string `json:"de,omitempty"`
	Fr string `json:"fr,omitempty"`
}
func (t Translation) Language(code Translator) string {
	switch code.Code() {
	case LangRu:
		return t.Ru
	case LangEn:
		return t.En
	case LangFr:
		return t.Fr
	case LangDe:
		return t.De
	default:
		return t.En
	}
}

type LangCode string

var _ Translator = (LangCode)("")

// reflected on translations/declinations struct keys
const (
	LangRu LangCode = "Ru"
	LangEn LangCode = "En"
	LangDe LangCode = "De"
	LangFr LangCode = "Fr"
	LangUnknown LangCode = ""
)

var AcceptedLangs = []string{"ru", "en", "de", "fr"}

func LangFromString(code string) LangCode {
	switch code {
	case LangRu.String(), strings.ToLower(LangRu.String()):
		return LangRu
	case LangEn.String(), strings.ToLower(LangEn.String()):
		return LangEn
	case LangDe.String(), strings.ToLower(LangDe.String()):
		return LangDe
	case LangFr.String(), strings.ToLower(LangFr.String()):
		return LangFr
	default:
		return LangEn
	}
}

func Accepts(code string) bool {
	for _, lang := range AcceptedLangs {
		if strings.ToLower(code) == lang {
			return true
		}
	}
	return false
}

func (lang LangCode) String() string { return string(lang) }

func (lang LangCode) Code() LangCode { return lang }
func (lang LangCode) Text(key string) string { return texts.Get(lang, key) }
func (lang LangCode) Error(code int) string { return texts.Get(lang, fmt.Sprintf("%d", code)) }
func (lang LangCode) Decl(key string) []string { return decls.Get(lang, key) }