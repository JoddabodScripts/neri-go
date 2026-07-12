package nerimity

import "context"

// ButtonComponentType is the kind of interactive component in a button
// response.
type ButtonComponentType string

const (
	// ComponentText is a read-only block of text.
	ComponentText ButtonComponentType = "text"
	// ComponentDropdown is a selectable dropdown.
	ComponentDropdown ButtonComponentType = "dropdown"
	// ComponentInput is a single-line text input.
	ComponentInput ButtonComponentType = "input"
)

// DropdownItem is one selectable option in a dropdown component.
type DropdownItem struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// ButtonComponent is one component of a button response modal. Build one with
// TextComponent, DropdownComponent or InputComponent rather than by hand.
type ButtonComponent struct {
	ID          string              `json:"id"`
	Type        ButtonComponentType `json:"type"`
	Content     string              `json:"content,omitempty"`     // text
	Label       string              `json:"label,omitempty"`       // dropdown, input
	Items       []DropdownItem      `json:"items,omitempty"`       // dropdown
	Placeholder string              `json:"placeholder,omitempty"` // input
}

// TextComponent returns a read-only text component.
func TextComponent(id, content string) ButtonComponent {
	return ButtonComponent{ID: id, Type: ComponentText, Content: content}
}

// DropdownComponent returns a dropdown component.
func DropdownComponent(id, label string, items []DropdownItem) ButtonComponent {
	return ButtonComponent{ID: id, Type: ComponentDropdown, Label: label, Items: items}
}

// InputComponent returns a single-line text input component.
func InputComponent(id, label, placeholder string) ButtonComponent {
	return ButtonComponent{ID: id, Type: ComponentInput, Label: label, Placeholder: placeholder}
}

// ButtonResponse is the payload sent back when responding to a button click.
type ButtonResponse struct {
	// Content is the response message text.
	Content string `json:"content,omitempty"`
	// Components are the interactive components to show.
	Components []ButtonComponent `json:"components,omitempty"`
	// Title overrides the response modal's title.
	Title string `json:"title,omitempty"`
	// ButtonLabel overrides the response's confirm button label.
	ButtonLabel string `json:"buttonLabel,omitempty"`
}

// MessageButton is a button-click interaction delivered with the
// MessageButtonClick event.
type MessageButton struct {
	client *Client

	// ID is the clicked button's ID (as set in ButtonOption).
	ID string
	// UserID is the ID of the user who clicked.
	UserID string
	// MessageID is the ID of the message the button is on.
	MessageID string
	// ChannelID is the ID of the channel the message is in.
	ChannelID string
	// Type is "button_click" or "modal_click".
	Type string
	// Data carries submitted modal component values, if any.
	Data map[string]string
	// User is the clicking user if cached, else nil.
	User *User
	// Channel is the channel if cached, else nil.
	Channel *Channel
	// Message is the message if cached, else nil (see Partial).
	Message *Message
	// Partial reports whether Message could not be resolved from cache; call
	// Fetch to populate it.
	Partial bool
}

func newMessageButton(client *Client, p messageButtonClickPayload) *MessageButton {
	b := &MessageButton{
		client:    client,
		ID:        p.ButtonID,
		UserID:    p.UserID,
		MessageID: p.MessageID,
		ChannelID: p.ChannelID,
		Type:      p.Type,
		Data:      p.Data,
		Partial:   true,
	}
	b.User, _ = client.users.get(p.UserID)
	b.Channel, _ = client.channels.get(p.ChannelID)
	if msg, ok := client.messages.get(p.MessageID); ok {
		b.Message = msg
		b.Partial = false
	}
	return b
}

// Fetch populates Message by fetching it from the API if it was not cached.
func (b *MessageButton) Fetch(ctx context.Context) error {
	if !b.Partial {
		return nil
	}
	msg, err := b.client.fetchMessage(ctx, b.ChannelID, b.MessageID)
	if err != nil {
		return err
	}
	b.Message = msg
	b.Partial = false
	return nil
}

// Respond replies to the button click. Provide content, components, and/or
// title and button-label overrides via resp.
func (b *MessageButton) Respond(ctx context.Context, resp ButtonResponse) error {
	return b.client.buttonCallback(ctx, b.ChannelID, b.MessageID, b.ID, b.UserID, resp)
}
