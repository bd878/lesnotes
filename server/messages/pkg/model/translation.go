package model

type (
	Translation struct {
		// TODO UserID
		MessageID           int64               `json:"message"`
		Lang                string              `json:"lang"`
		Title               string              `json:"title"`
		Text                string              `json:"text"`
		CreatedAt           string              `json:"created_at"`
		UpdatedAt           string              `json:"updated_at"`
	}

	TranslationPreview struct {
		MessageID           int64               `json:"message"`
		Lang                string              `json:"lang"`
		Title               string              `json:"title"`
		CreatedAt           string              `json:"created_at"`
		UpdatedAt           string              `json:"updated_at"`
	}

	SendTranslationRequest struct {
		MessageID           int64               `json:"message"`
		Lang                string              `json:"lang"`
		Title               string              `json:"title"`
		Text                string              `json:"text"`
	}

	SendTranslationResponse struct {
		Description         string              `json:"description"`
	}

	UpdateTranslationRequest struct {
		MessageID           int64               `json:"message"`
		Lang                string              `json:"lang"`
		Title               *string             `json:"title,omitempty"`
		Text                *string             `json:"text,omitempty"`
	}

	UpdateTranslationResponse struct {
		Description         string              `json:"description"`
	}

	DeleteTranslationRequest struct {
		MessageID           int64               `json:"message"`
		Lang                string              `json:"lang"`
	}

	DeleteTranslationResponse struct {
		Description         string              `json:"description"`
	}

	ReadTranslationRequest struct {
		MessageID           int64               `json:"message"`
		Lang                string              `json:"lang"`
		Name                *string             `json:"name,omitempty"`
	}

	ReadTranslationResponse struct {
		Translation         *Translation        `json:"translation"`
	}

	ListTranslationsRequest struct {
		MessageID           int64               `json:"message"`
		Name                string              `json:"name,omitempty"`
	}

	ListTranslationsResponse struct {
		Translations        []*Translation      `json:"translations"`
	}
)