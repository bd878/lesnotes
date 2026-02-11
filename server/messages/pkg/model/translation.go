package model

type (
	Translation struct {
		MessageID           int64               `json:"message"`
		Lang                string              `json:"lang"`
		Title               string              `json:"title"`
		Text                string              `json:"text"`
		CreatedAt           int64               `json:"created_at"`
		UpdatedAt           int64               `json:"updated_at"`
	}

	TranslationPreview struct {
		MessageID           int64               `json:"message"`
		Lang                string              `json:"lang"`
		Title               string              `json:"title"`
		CreatedAt           int64               `json:"created_at"`
		UpdatedAt           int64               `json:"updated_at"`
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
		Name                string              `json:"name,omitempty"`
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