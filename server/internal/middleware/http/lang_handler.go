package middleware

import (
	"context"
	"net/http"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/i18n"
	"github.com/bd878/gallery/server/third_party/accept"
)

type LangContextKey struct {}

func Language(next Handler) Handler {
	return handler(func(w http.ResponseWriter, req *http.Request) (err error) {
		var lang i18n.LangCode

		languageHint := req.Header.Get("X-Language")
		if i18n.Accepts(languageHint) {
			lang = i18n.LangFromString(languageHint)
		} else {
			languages := req.Header.Get("Accept-Language")
			preferredLang, err := accept.Negotiate(languages, i18n.AcceptedLangs...)
			if err != nil {
				logger.Errorw("lang middleware", "error", err)
				lang = i18n.LangEn
			} else {
				lang = i18n.LangFromString(preferredLang)
			}
		}

		req = req.WithContext(context.WithValue(req.Context(), LangContextKey{}, lang))

		return next.Handle(w, req)
	})
}