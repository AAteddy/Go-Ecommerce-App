package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AAteddy/go-ecommerce-app/internal/pkg/errors"
)

// StringArray is a custom type for PostgreSQL text arrays
type StringArray []string

// Value implements the driver.Valuer interface for database serialization
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	// Format as PostgreSQL array literal: {"value1","value2"}
	quotedValues := make([]string, len(a))
	for i, v := range a {
		quotedValues[i] = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	}
	return "{" + strings.Join(quotedValues, ",") + "}", nil
}

// Scan implements the sql.Scanner interface for database deserialization
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan StringArray: expected string, got %T", value)
	}

	// Remove curly braces
	str = strings.Trim(str, "{}")
	if str == "" {
		*a = StringArray{}
		return nil
	}

	// Split by comma, handling quoted values
	*a = parsePostgresArray(str)
	return nil
}

// parsePostgresArray parses PostgreSQL array literal format
func parsePostgresArray(s string) StringArray {
	var result []string
	var current strings.Builder
	inQuotes := false
	escapeNext := false

	for _, char := range s {
		if escapeNext {
			current.WriteRune(char)
			escapeNext = false
			continue
		}

		switch char {
		case '\\':
			escapeNext = true
		case '"':
			inQuotes = !inQuotes
		case ',':
			if !inQuotes {
				result = append(result, current.String())
				current.Reset()
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// MarshalJSON implements json.Marshaler
func (a StringArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(a))
}

// UnmarshalJSON implements json.Unmarshaler
func (a *StringArray) UnmarshalJSON(data []byte) error {
	var strArray []string
	if err := json.Unmarshal(data, &strArray); err != nil {
		return err
	}
	*a = strArray
	return nil
}

// Order represents an e-commerce order.
type Order struct {
	ID         uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	ProductIDs StringArray `json:"product_id" gorm:"type:text[]; not null"`
	Quantity   int         `json:"quantity" gorm:"not null;default:1"`
	Total      float64     `json:"total" gorm:"type:decimal(10,2);not null;default:0.00"`
	Status     string      `json:"status" gorm:"size:50;not null;default:'pending'"`
	CreatedAt  time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

// NewOrder validates and creates a new order.
func NewOrder(userID uuid.UUID, productIDs []string, total float64, quantity int) (*Order, error) {
	// Validate input
	if userID == uuid.Nil || len(productIDs) == 0 || total <= 0 || quantity <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return &Order{
		ID:         uuid.New(),
		UserID:     userID,
		ProductIDs: StringArray(productIDs),
		Quantity:   quantity,
		Total:      total,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}
