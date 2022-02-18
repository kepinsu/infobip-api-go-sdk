package models

import (
	"fmt"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"mvdan.cc/xurls/v2"
)

const maxInteractiveListRows = 10

type BulkMessageResponse struct {
	Messages []MessageResponse `json:"messages"`
	BulkID   string            `json:"bulkId"`
}

type MessageResponse struct {
	To           string `json:"to"`
	MessageCount int32  `json:"messageCount"`
	MessageID    string `json:"messageId"`
	Status       Status `json:"status"`
}

type Status struct {
	GroupID     int32  `json:"groupId"`
	GroupName   string `json:"groupName"`
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Action      string `json:"action"`
}

type TemplatesResponse struct {
	Templates []TemplateResponse `json:"templates"`
}

type TemplateResponse struct {
	ID                string            `json:"ID"`
	BusinessAccountID int64             `json:"businessAccountID"`
	Name              string            `json:"name"`
	Language          string            `json:"language"`
	Status            string            `json:"status"`
	Category          string            `json:"category"`
	Structure         TemplateStructure `json:"structure"`
}

type TemplateStructure struct {
	Header  *TemplateHeader  `json:"header,omitempty"`
	Body    string           `json:"body" validate:"required"`
	Footer  string           `json:"footer,omitempty" validate:"lte=60"`
	Buttons []TemplateButton `json:"buttons,omitempty" validate:"omitempty,min=1,max=3,dive"`
	Type    string           `json:"type,omitempty" validate:"oneof=TEXT MEDIA UNSUPPORTED"`
}

type TemplateHeader struct {
	Format string `json:"format,omitempty" validate:"oneof=TEXT IMAGE VIDEO DOCUMENT LOCATION"`
	Text   string `json:"text" validate:"lte=60"`
}

type TemplateButton struct {
	Type        string `json:"type" validate:"required,oneof=PHONE_NUMBER URL QUICK_REPLY"`
	Text        string `json:"text" validate:"required,lte=200"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	URL         string `json:"url,omitempty" validate:"omitempty,url"`
}

type TemplateCreate struct {
	Name      string            `json:"name" validate:"required"`
	Language  string            `json:"language" validate:"required"`
	Category  string            `json:"category" validate:"required"`
	Structure TemplateStructure `json:"structure" validate:"required"`
}

func (t *TemplateCreate) Validate() error {
	validate = validator.New()
	validate.RegisterStructValidation(templateCreateValidation, TemplateCreate{})
	validate.RegisterStructValidation(templateHeaderValidation, TemplateHeader{})
	validate.RegisterStructValidation(templateButtonValidation, TemplateButton{})
	return validate.Struct(t)
}

func templateCreateValidation(sl validator.StructLevel) {
	template, _ := sl.Current().Interface().(TemplateCreate)
	validateTemplateName(sl, template)
	validateTemplateLanguage(sl, template)
	validateCategory(sl, template)
}

func validateTemplateName(sl validator.StructLevel, template TemplateCreate) {
	if !isSnakeCase(template.Name) {
		sl.ReportError(template.Name, "name", "Name", "namenotsnakecase", "")
	}
}

func validateTemplateLanguage(sl validator.StructLevel, template TemplateCreate) {
	switch template.Language {
	case "af", "sq", "ar", "az", "bn", "bg", "ca", "zh_CN", "zh_HK", "zh_TW", "hr", "cs", "da", "nl", "en", "en_GB",
		"en_US", "et", "fil", "fi", "fr", "de", "el", "gu", "ha", "he", "hi", "hu", "id", "ga", "it", "ja", "kn", "kk",
		"ko", "lo", "lv", "lt", "mk", "ms", "ml", "mr", "nb", "fa", "pl", "pt_BR", "pt_PT", "pa", "ro", "ru", "sr",
		"sk", "sl", "es", "es_AR", "es_ES", "es_MX", "sw", "sv", "ta", "te", "th", "tr", "uk", "ur", "uz", "vi",
		"unknown":
	default:
		sl.ReportError(template.Language, "language", "Language", "invalidlanguage", "")
	}
}

func validateCategory(sl validator.StructLevel, template TemplateCreate) {
	switch template.Category {
	case "ACCOUNT_UPDATE", "PAYMENT_UPDATE", "PERSONAL_FINANCE_UPDATE", "SHIPPING_UPDATE",
		"RESERVATION_UPDATE", "ISSUE_RESOLUTION", "APPOINTMENT_UPDATE", "TRANSPORTATION_UPDATE",
		"TICKET_UPDATE", "ALERT_UPDATE", "AUTO_REPLY":
	default:
		sl.ReportError(template.Category, "category", "Category", "invalidcategory", "")
	}
}

func templateHeaderValidation(sl validator.StructLevel) {
	header, _ := sl.Current().Interface().(TemplateHeader)
	if header.Format == "TEXT" && header.Text == "" {
		sl.ReportError(header.Text, "text", "Text", "missingtext", "")
	}
}

func templateButtonValidation(sl validator.StructLevel) {
	button, _ := sl.Current().Interface().(TemplateButton)
	switch button.Type {
	case "PHONE_NUMBER":
		if button.PhoneNumber == "" {
			sl.ReportError(button.PhoneNumber, "phoneNumber", "PhoneNumber", "required", "")
		}
	case "URL":
		if button.URL == "" {
			sl.ReportError(button.URL, "url", "URL", "required", "")
		}
	}
}

type MessageCommon struct {
	From         string `json:"from" validate:"required,lte=24"`
	To           string `json:"to" validate:"required,lte=24"`
	MessageID    string `json:"messageId,omitempty" validate:"lte=50"`
	CallbackData string `json:"callbackData,omitempty" validate:"lte=4000"`
	NotifyURL    string `json:"notifyUrl,omitempty" validate:"omitempty,url,lte=2048"`
}

type TemplateMessages struct {
	Messages []TemplateMessage `json:"messages" validate:"required,min=1,dive"`
	BulkID   string            `json:"bulkId,omitempty" validate:"lte=100"`
}

func (t *TemplateMessages) Validate() error {
	validate = validator.New()
	validate.RegisterStructValidation(templateMsgsValidation, TemplateMessageContent{})
	validate.RegisterStructValidation(templateHeaderSendValidation, TemplateHeaderSend{})
	validate.RegisterStructValidation(templateDataValidation, TemplateData{})
	validate.RegisterStructValidation(templateButtonSendValidation, TemplateButtonSend{})
	return validate.Struct(t)
}

func templateMsgsValidation(sl validator.StructLevel) {
	content, _ := sl.Current().Interface().(TemplateMessageContent)

	if !isSnakeCase(content.TemplateName) {
		sl.ReportError(content.TemplateName, "templateName", "TemplateName", "templatenamenotsnakecase", "")
	}
}

func isSnakeCase(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && r != '_' {
			return false
		}
	}
	return true
}

func templateHeaderSendValidation(sl validator.StructLevel) {
	header, _ := sl.Current().Interface().(TemplateHeaderSend)
	switch header.Type {
	case "TEXT":
		if header.Placeholder == "" {
			sl.ReportError(header.Placeholder, "placholder", "Placeholder", "missingplaceholder", "")
		}
	case "DOCUMENT":
		if header.MediaURL == "" {
			sl.ReportError(header.MediaURL, "mediaUrl", "MediaURL", "missingmediaurl", "")
		}
		if header.Filename == "" {
			sl.ReportError(header.Filename, "filename", "Filename", "missingfilename", "")
		}
	case "VIDEO", "IMAGE":
		if header.MediaURL == "" {
			sl.ReportError(header.MediaURL, "mediaUrl", "MediaURL", "missingmediaurl", "")
		}
	case "LOCATION":
		if header.Latitude == nil {
			sl.ReportError(header.Latitude, "latitude", "Latitude", "missinglatitude", "")
		}
		if header.Longitude == nil {
			sl.ReportError(header.Longitude, "longitude", "Longitude", "missinglongitude", "")
		}
	}
}

func templateDataValidation(sl validator.StructLevel) {
	templateData, _ := sl.Current().Interface().(TemplateData)
	validateTemplateButtonLength(sl, templateData)
	validateTemplateButtonTypes(sl, templateData)
}

func validateTemplateButtonLength(sl validator.StructLevel, templateData TemplateData) {
	if len(templateData.Buttons) > 1 && templateData.Buttons[0].Type == "URL" {
		sl.ReportError(templateData.Buttons, "buttons", "Buttons", "dynamicurlcountoverone", "")
	}
}

func validateTemplateButtonTypes(sl validator.StructLevel, templateData TemplateData) {
	types := map[string]int{"QUICK_REPLY": 0, "URL": 0}
	for _, button := range templateData.Buttons {
		types[button.Type]++
	}
	if types["QUICK_REPLY"] > 0 && types["URL"] > 0 {
		sl.ReportError(templateData.Buttons, "buttons", "Buttons", "bothreplyandurlpresent", "")
	}
}

func templateButtonSendValidation(sl validator.StructLevel) {
	button, _ := sl.Current().Interface().(TemplateButtonSend)
	if button.Type == "QUICK_REPLY" && len(button.Parameter) > 128 {
		sl.ReportError(button.Parameter, "parameter", "Parameter", "parametertoolong", "")
	}
}

type TemplateMessage struct {
	MessageCommon
	Content     TemplateMessageContent `json:"content" validate:"required"`
	SMSFailover *SMSFailover           `json:"smsFailover,omitempty"`
}

type TemplateMessageContent struct {
	TemplateName string       `json:"templateName" validate:"required,lte=512"`
	TemplateData TemplateData `json:"templateData" validate:"required"`
	Language     string       `json:"language" validate:"required"`
}

type TemplateData struct {
	Body    TemplateBody         `json:"body" validate:"required"`
	Header  *TemplateHeaderSend  `json:"header,omitempty"`
	Buttons []TemplateButtonSend `json:"buttons,omitempty" validate:"omitempty,max=3,dive"`
}

type TemplateBody struct {
	Placeholders []string `json:"placeholders" validate:"required,dive,gte=1"`
}

type TemplateHeaderSend struct {
	Type        string   `json:"type" validate:"required,oneof=TEXT DOCUMENT IMAGE VIDEO LOCATION"`
	Placeholder string   `json:"placeholder,omitempty"`
	MediaURL    string   `json:"mediaUrl,omitempty" validate:"omitempty,url,lte=2048"`
	Filename    string   `json:"filename,omitempty" validate:"lte=240"`
	Latitude    *float32 `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude   *float32 `json:"longitude,omitempty" validate:"omitempty,longitude"`
}

type TemplateButtonSend struct {
	Type      string `json:"type" validate:"required,oneof=QUICK_REPLY URL"`
	Parameter string `json:"parameter" validate:"required"`
}

type SMSFailover struct {
	From string `json:"from" validate:"required,lte=24"`
	Text string `json:"text" validate:"required,lte=4096"`
}

type TextMessage struct {
	MessageCommon
	Content TextContent `json:"content" validate:"required"`
}

func (t *TextMessage) Validate() error {
	validate = validator.New()
	validate.RegisterStructValidation(previewURLValidation, TextContent{})
	return validate.Struct(t)
}

type TextContent struct {
	Text       string `json:"text" validate:"required,gte=1,lte=4096"`
	PreviewURL bool   `json:"previewURL,omitempty"`
}

func previewURLValidation(sl validator.StructLevel) {
	content, _ := sl.Current().Interface().(TextContent)
	containsURL := xurls.Relaxed().FindString(content.Text)
	if content.PreviewURL && containsURL == "" {
		sl.ReportError(content.Text, "text", "Text", "missingurlintext", "")
	}
}

type DocumentMessage struct {
	MessageCommon
	Content DocumentContent `json:"content" validate:"required"`
}

func (t *DocumentMessage) Validate() error {
	validate = validator.New()
	return validate.Struct(t)
}

type DocumentContent struct {
	MediaURL string `json:"mediaUrl" validate:"required,url,lte=2048"`
	Caption  string `json:"caption,omitempty" validate:"lte=3000"`
	Filename string `json:"filename,omitempty" validate:"lte=240"`
}

type ImageMessage struct {
	MessageCommon
	Content ImageContent `json:"content" validate:"required"`
}

func (t *ImageMessage) Validate() error {
	validate = validator.New()
	return validate.Struct(t)
}

type ImageContent struct {
	MediaURL string `json:"mediaUrl" validate:"required,url,lte=2048"`
	Caption  string `json:"caption,omitempty" validate:"lte=3000"`
}

type AudioMessage struct {
	MessageCommon
	Content AudioContent `json:"content" validate:"required"`
}

func (t *AudioMessage) Validate() error {
	validate = validator.New()
	return validate.Struct(t)
}

type AudioContent struct {
	MediaURL string `json:"mediaUrl" validate:"required,url,lte=2048"`
}

type VideoMessage struct {
	MessageCommon
	Content VideoContent `json:"content" validate:"required"`
}

func (t *VideoMessage) Validate() error {
	validate = validator.New()
	return validate.Struct(t)
}

type VideoContent struct {
	MediaURL string `json:"mediaUrl" validate:"required,url,lte=2048"`
	Caption  string `json:"caption,omitempty" validate:"lte=3000"`
}

type StickerMessage struct {
	MessageCommon
	Content StickerContent `json:"content" validate:"required"`
}

func (t *StickerMessage) Validate() error {
	validate = validator.New()
	return validate.Struct(t)
}

type StickerContent struct {
	MediaURL string `json:"mediaUrl" validate:"required,url,lte=2048"`
}

type LocationMessage struct {
	MessageCommon
	Content LocationContent `json:"content" validate:"required"`
}

func (t *LocationMessage) Validate() error {
	validate = validator.New()
	return validate.Struct(t)
}

type LocationContent struct {
	Latitude  *float32 `json:"latitude" validate:"required,latitude"`
	Longitude *float32 `json:"longitude" validate:"required,longitude"`
	Name      string   `json:"name" validate:"lte=1000"`
	Address   string   `json:"address" validate:"lte=1000"`
}

type ContactMessage struct {
	MessageCommon
	Content ContactContent `json:"content" validate:"required"`
}

func (t *ContactMessage) Validate() error {
	validate = validator.New()
	validate.RegisterStructValidation(birthdayValidation, Contact{})
	return validate.Struct(t)
}

func birthdayValidation(sl validator.StructLevel) {
	contact, _ := sl.Current().Interface().(Contact)
	if contact.Birthday == "" {
		return
	}
	_, err := time.Parse("2006-01-02", contact.Birthday)
	if err != nil {
		sl.ReportError(contact.Birthday, "birthday", "Contact", "invalidbirthdayformat", "")
	}
}

type ContactContent struct {
	Contacts []Contact `json:"contacts" validate:"required,dive"`
}

type Contact struct {
	Addresses []ContactAddress `json:"addresses,omitempty" validate:"omitempty,dive"`
	Birthday  string           `json:"birthday,omitempty"`
	Emails    []ContactEmail   `json:"emails,omitempty" validate:"omitempty,dive"`
	Name      ContactName      `json:"name" validate:"required"`
	Org       ContactOrg       `json:"org,omitempty"`
	Phones    []ContactPhone   `json:"phones,omitempty" validate:"omitempty,dive"`
	Urls      []ContactURL     `json:"urls,omitempty" validate:"omitempty,dive"`
}

type ContactAddress struct {
	Street      string `json:"street,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Zip         string `json:"zip,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	Type        string `json:"type,omitempty" validate:"omitempty,oneof=HOME WORK"`
}

type ContactEmail struct {
	Email string `json:"email,omitempty" validate:"omitempty,email"`
	Type  string `json:"type,omitempty" validate:"omitempty,oneof=HOME WORK"`
}

type ContactName struct {
	FirstName     string `json:"firstName" validate:"required"`
	LastName      string `json:"lastName,omitempty"`
	MiddleName    string `json:"middleName,omitempty"`
	NameSuffix    string `json:"nameSuffix,omitempty"`
	NamePrefix    string `json:"namePrefix,omitempty"`
	FormattedName string `json:"formattedName" validate:"required"`
}

type ContactOrg struct {
	Company    string `json:"company,omitempty"`
	Department string `json:"department,omitempty"`
	Title      string `json:"title,omitempty"`
}

type ContactPhone struct {
	Phone string `json:"phone,omitempty"`
	Type  string `json:"type,omitempty" validate:"omitempty,oneof=CELL MAIN IPHONE HOME WORK"`
	WaID  string `json:"waId,omitempty"`
}

type ContactURL struct {
	URL  string `json:"url,omitempty" validate:"omitempty,url"`
	Type string `json:"type,omitempty" validate:"omitempty,oneof=HOME WORK"`
}

type InteractiveButtonsMessage struct {
	MessageCommon
	Content InteractiveButtonsContent `json:"content" validate:"required"`
}

func (t *InteractiveButtonsMessage) Validate() error {
	validate = validator.New()
	validate.RegisterStructValidation(buttonHeaderValidation, InteractiveButtonsHeader{})
	return validate.Struct(t)
}

func buttonHeaderValidation(sl validator.StructLevel) {
	header, _ := sl.Current().Interface().(InteractiveButtonsHeader)
	switch header.Type {
	case "TEXT":
		if header.Text == "" {
			sl.ReportError(header.Text, "text", "Text", "missingtext", "")
		}
	case "VIDEO", "IMAGE", "DOCUMENT":
		if header.MediaURL == "" {
			sl.ReportError(header.MediaURL, "mediaUrl", "MediaURL", "missingmediaurl", "")
		}
	}
}

type InteractiveButtonsContent struct {
	Body   InteractiveButtonsBody    `json:"body" validate:"required"`
	Action InteractiveButtons        `json:"action" validate:"required"`
	Header *InteractiveButtonsHeader `json:"header,omitempty" validate:"omitempty"`
	Footer *InteractiveButtonsFooter `json:"footer,omitempty"`
}

type InteractiveButtonsBody struct {
	Text string `json:"text" validate:"required,lte=1024"`
}

type InteractiveButtons struct {
	Buttons []InteractiveButton `json:"buttons" validate:"required,min=1,max=3,dive"`
}

type InteractiveButton struct {
	Type  string `json:"type" validate:"required,oneof=REPLY"`
	ID    string `json:"id" validate:"required,lte=256"`
	Title string `json:"title" validate:"required,lte=20"`
}

type InteractiveButtonsHeader struct {
	Type     string `json:"type" validate:"required,oneof=TEXT VIDEO IMAGE DOCUMENT"`
	Text     string `json:"text,omitempty" validate:"lte=60"`
	MediaURL string `json:"mediaUrl,omitempty" validate:"omitempty,url,lte=2048"`
	Filename string `json:"filename,omitempty" validate:"lte=240"`
}

type InteractiveButtonsFooter struct {
	Text string `json:"text" validate:"required,lte=60"`
}

type InteractiveListMessage struct {
	MessageCommon
	Content InteractiveListContent `json:"content" validate:"required"`
}

func (t *InteractiveListMessage) Validate() error {
	validate = validator.New()
	validate.RegisterStructValidation(interactiveListActionValidation, InteractiveListAction{})
	return validate.Struct(t)
}

func interactiveListActionValidation(sl validator.StructLevel) {
	action, _ := sl.Current().Interface().(InteractiveListAction)
	validateRowCount(sl, action)
	validateDuplicateRows(sl, action)
	validateListSectionTitles(sl, action)
}

func validateRowCount(sl validator.StructLevel, action InteractiveListAction) {
	var rowCount int
	for _, section := range action.Sections {
		rowCount += len(section.Rows)
		if rowCount > maxInteractiveListRows {
			sl.ReportError(
				action.Sections,
				"sections",
				"Sections",
				"rowcountovermax",
				"",
			)
		}
	}
}

func validateDuplicateRows(sl validator.StructLevel, action InteractiveListAction) {
	rowIDs := make(map[string]int)
	for _, section := range action.Sections {
		for _, row := range section.Rows {
			rowIDs[row.ID]++
			if rowIDs[row.ID] > 1 {
				sl.ReportError(
					action.Sections,
					"sections",
					"Sections",
					fmt.Sprintf("duplicaterowID%s", row.ID),
					"",
				)
			}
		}
	}
}

func validateListSectionTitles(sl validator.StructLevel, action InteractiveListAction) {
	if len(action.Sections) > 1 {
		for _, section := range action.Sections {
			if section.Title == "" {
				sl.ReportError(
					action.Sections,
					"sections",
					"Sections",
					"missingtitlemultiplesections",
					"",
				)
			}
		}
	}
}

type InteractiveListContent struct {
	Body   InteractiveListBody    `json:"body" validate:"required"`
	Action InteractiveListAction  `json:"action" validate:"required"`
	Header *InteractiveListHeader `json:"header,omitempty" validate:"omitempty"`
	Footer *InteractiveListFooter `json:"footer,omitempty"`
}

type InteractiveListBody struct {
	Text string `json:"text" validate:"required,lte=1024"`
}

type InteractiveListAction struct {
	Title    string                   `json:"title" validate:"required,lte=20"`
	Sections []InteractiveListSection `json:"sections" validate:"required,min=1,max=10,dive"`
}

type InteractiveListSection struct {
	Title string       `json:"title,omitempty" validate:"lte=24"`
	Rows  []SectionRow `json:"rows" validate:"required,min=1,max=10,dive"`
}

type SectionRow struct {
	ID          string `json:"id" validate:"required,lte=200"`
	Title       string `json:"title" validate:"required,lte=24"`
	Description string `json:"description,omitempty" validate:"lte=72"`
}

type InteractiveListHeader struct {
	Type string `json:"type" validate:"required,oneof=TEXT"`
	Text string `json:"text" validate:"required,lte=60"`
}

type InteractiveListFooter struct {
	Text string `json:"text" validate:"required,lte=60"`
}

type InteractiveProductMessage struct {
	MessageCommon
	Content InteractiveProductContent `json:"content" validate:"required"`
}

func (t *InteractiveProductMessage) Validate() error {
	validate = validator.New()
	return validate.Struct(t)
}

type InteractiveProductContent struct {
	Action InteractiveProductAction  `json:"action" validate:"required"`
	Body   *InteractiveProductBody   `json:"body,omitempty"`
	Footer *InteractiveProductFooter `json:"footer,omitempty"`
}

type InteractiveProductAction struct {
	CatalogID         string `json:"catalogId" validate:"required"`
	ProductRetailerID string `json:"productRetailerId" validate:"required"`
}

type InteractiveProductBody struct {
	Text string `json:"text" validate:"required,lte=1024"`
}

type InteractiveProductFooter struct {
	Text string `json:"text" validate:"required,lte=60"`
}

type InteractiveMultiproductMessage struct {
	MessageCommon
	Content InteractiveMultiproductContent `json:"content" validate:"required"`
}

func (t *InteractiveMultiproductMessage) Validate() error {
	validate = validator.New()
	validate.RegisterStructValidation(multiproductActionValidation, InteractiveMultiproductAction{})
	return validate.Struct(t)
}

func multiproductActionValidation(sl validator.StructLevel) {
	action, _ := sl.Current().Interface().(InteractiveMultiproductAction)
	validateMultiproductSectionTitles(sl, action)
}

func validateMultiproductSectionTitles(sl validator.StructLevel, action InteractiveMultiproductAction) {
	if len(action.Sections) > 1 {
		for _, section := range action.Sections {
			if section.Title == "" {
				sl.ReportError(
					action.Sections,
					"sections",
					"Sections",
					"missingtitlemultiplesections",
					"",
				)
			}
		}
	}
}

type InteractiveMultiproductContent struct {
	Header InteractiveMultiproductHeader  `json:"header" validate:"required"`
	Body   InteractiveMultiproductBody    `json:"body" validate:"required"`
	Action InteractiveMultiproductAction  `json:"action" validate:"required"`
	Footer *InteractiveMultiproductFooter `json:"footer,omitempty"`
}

type InteractiveMultiproductHeader struct {
	Type string `json:"type" validate:"required,oneof=TEXT"`
	Text string `json:"text" validate:"required,lte=60"`
}

type InteractiveMultiproductBody struct {
	Text string `json:"text" validate:"required,lte=1024"`
}

type InteractiveMultiproductAction struct {
	CatalogID string                           `json:"catalogId" validate:"required"`
	Sections  []InteractiveMultiproductSection `json:"sections" validate:"required,min=1,max=10,dive"`
}

type InteractiveMultiproductSection struct {
	Title              string   `json:"title,omitempty" validate:"lte=24"`
	ProductRetailerIDs []string `json:"productRetailerIds" validate:"required,min=1"`
}

type InteractiveMultiproductFooter struct {
	Text string `json:"text" validate:"required,lte=60"`
}
