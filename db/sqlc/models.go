// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type Role string

const (
	RoleProvider Role = "provider"
	RoleClient   Role = "client"
	RolePartner  Role = "partner"
	RoleAdmin    Role = "admin"
)

func (e *Role) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Role(s)
	case string:
		*e = Role(s)
	default:
		return fmt.Errorf("unsupported scan type for Role: %T", src)
	}
	return nil
}

type NullRole struct {
	Role  Role `json:"role"`
	Valid bool `json:"valid"` // Valid is true if Role is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRole) Scan(value interface{}) error {
	if value == nil {
		ns.Role, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Role.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Role), nil
}

type StatusOrders string

const (
	StatusOrdersPending   StatusOrders = "pending"
	StatusOrdersAccepted  StatusOrders = "accepted"
	StatusOrdersCompleted StatusOrders = "completed"
	StatusOrdersCancelled StatusOrders = "cancelled"
)

func (e *StatusOrders) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = StatusOrders(s)
	case string:
		*e = StatusOrders(s)
	default:
		return fmt.Errorf("unsupported scan type for StatusOrders: %T", src)
	}
	return nil
}

type NullStatusOrders struct {
	StatusOrders StatusOrders `json:"status_orders"`
	Valid        bool         `json:"valid"` // Valid is true if StatusOrders is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStatusOrders) Scan(value interface{}) error {
	if value == nil {
		ns.StatusOrders, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.StatusOrders.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStatusOrders) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.StatusOrders), nil
}

type StatusPayment string

const (
	StatusPaymentPending   StatusPayment = "pending"
	StatusPaymentCompleted StatusPayment = "completed"
	StatusPaymentFailed    StatusPayment = "failed"
)

func (e *StatusPayment) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = StatusPayment(s)
	case string:
		*e = StatusPayment(s)
	default:
		return fmt.Errorf("unsupported scan type for StatusPayment: %T", src)
	}
	return nil
}

type NullStatusPayment struct {
	StatusPayment StatusPayment `json:"status_payment"`
	Valid         bool          `json:"valid"` // Valid is true if StatusPayment is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStatusPayment) Scan(value interface{}) error {
	if value == nil {
		ns.StatusPayment, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.StatusPayment.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStatusPayment) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.StatusPayment), nil
}

type StatusSubscription string

const (
	StatusSubscriptionActive   StatusSubscription = "active"
	StatusSubscriptionInactive StatusSubscription = "inactive"
	StatusSubscriptionExpired  StatusSubscription = "expired"
)

func (e *StatusSubscription) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = StatusSubscription(s)
	case string:
		*e = StatusSubscription(s)
	default:
		return fmt.Errorf("unsupported scan type for StatusSubscription: %T", src)
	}
	return nil
}

type NullStatusSubscription struct {
	StatusSubscription StatusSubscription `json:"status_subscription"`
	Valid              bool               `json:"valid"` // Valid is true if StatusSubscription is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStatusSubscription) Scan(value interface{}) error {
	if value == nil {
		ns.StatusSubscription, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.StatusSubscription.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStatusSubscription) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.StatusSubscription), nil
}

type City struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Favorite struct {
	ID         int64     `json:"id"`
	ClientID   int64     `json:"client_id"`
	ProviderID int64     `json:"provider_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Message struct {
	ID         int64        `json:"id"`
	SenderID   int64        `json:"sender_id"`
	ReceiverID int64        `json:"receiver_id"`
	Content    string       `json:"content"`
	IsRead     sql.NullBool `json:"is_read"`
	CreatedAt  time.Time    `json:"created_at"`
}

type Order struct {
	ID                 int64            `json:"id"`
	ClientID           int64            `json:"client_id"`
	CategoryID         int64            `json:"category_id"`
	ServiceID          sql.NullInt64    `json:"service_id"`
	Status             NullStatusOrders `json:"status"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	ProviderAccepted   sql.NullBool     `json:"provider_accepted"`
	ProviderMessage    sql.NullString   `json:"provider_message"`
	ClientMessage      sql.NullString   `json:"client_message"`
	OrderDate          sql.NullTime     `json:"order_date"`
	SelectedProviderID sql.NullInt64    `json:"selected_provider_id"`
	SubtitleCategoryID sql.NullInt64    `json:"subtitle_category_id"`
}

type OrderResponse struct {
	ID           int64          `json:"id"`
	OrderID      int64          `json:"order_id"`
	ProviderID   int64          `json:"provider_id"`
	Message      sql.NullString `json:"message"`
	OfferedPrice sql.NullString `json:"offered_price"`
	IsSelected   sql.NullBool   `json:"is_selected"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type PartnerStatistic struct {
	ID                  int64         `json:"id"`
	PartnerID           int64         `json:"partner_id"`
	ProvidersAttracted  sql.NullInt32 `json:"providers_attracted"`
	TotalSubscriptions  sql.NullInt32 `json:"total_subscriptions"`
	ActiveSubscriptions sql.NullInt32 `json:"active_subscriptions"`
	UpdatedAt           time.Time     `json:"updated_at"`
}

type Payment struct {
	ID            int64             `json:"id"`
	UserID        int64             `json:"user_id"`
	Amount        string            `json:"amount"`
	PaymentMethod string            `json:"payment_method"`
	Status        NullStatusPayment `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type PendingSubscription struct {
	ID               int64         `json:"id"`
	PaymentID        int64         `json:"payment_id"`
	UserID           int64         `json:"user_id"`
	SubscriptionType string        `json:"subscription_type"`
	StartDate        time.Time     `json:"start_date"`
	EndDate          time.Time     `json:"end_date"`
	OriginalPrice    string        `json:"original_price"`
	FinalPrice       string        `json:"final_price"`
	PromoCodeID      sql.NullInt64 `json:"promo_code_id"`
	IsUpdate         bool          `json:"is_update"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type PromoCode struct {
	ID                 int64         `json:"id"`
	PartnerID          int64         `json:"partner_id"`
	Code               string        `json:"code"`
	DiscountPercentage int32         `json:"discount_percentage"`
	ValidUntil         time.Time     `json:"valid_until"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	MaxUsages          sql.NullInt32 `json:"max_usages"`
	CurrentUsages      sql.NullInt32 `json:"current_usages"`
}

type Review struct {
	ID         int64     `json:"id"`
	ClientID   int64     `json:"client_id"`
	ProviderID int64     `json:"provider_id"`
	Rating     int32     `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Service struct {
	ID                 int64          `json:"id"`
	ProviderID         int64          `json:"provider_id"`
	CategoryID         int64          `json:"category_id"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	Price              string         `json:"price"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	Subcategory        sql.NullString `json:"subcategory"`
	Country            sql.NullString `json:"country"`
	City               sql.NullString `json:"city"`
	District           sql.NullString `json:"district"`
	SubtitleCategoryID sql.NullInt64  `json:"subtitle_category_id"`
}

type ServiceCategory struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	Slug        string    `json:"slug"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Subscription struct {
	ID               int64                  `json:"id"`
	ProviderID       int64                  `json:"provider_id"`
	StartDate        time.Time              `json:"start_date"`
	EndDate          time.Time              `json:"end_date"`
	Status           NullStatusSubscription `json:"status"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	SubscriptionType sql.NullString         `json:"subscription_type"`
	Price            sql.NullString         `json:"price"`
	PromoCodeID      sql.NullInt64          `json:"promo_code_id"`
	OriginalPrice    sql.NullString         `json:"original_price"`
}

type SubtitleCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SupportMessage struct {
	ID        int64     `json:"id"`
	SenderID  int64     `json:"sender_id"`
	Subject   string    `json:"subject"`
	Messages  string    `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID               int64          `json:"id"`
	Username         string         `json:"username"`
	Email            string         `json:"email"`
	PasswordHash     string         `json:"password_hash"`
	PasswordChangeAt time.Time      `json:"password_change_at"`
	Role             NullRole       `json:"role"`
	Country          sql.NullString `json:"country"`
	City             sql.NullString `json:"city"`
	District         sql.NullString `json:"district"`
	Phone            string         `json:"phone"`
	Whatsapp         string         `json:"whatsapp"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	PhotoUrl         []byte         `json:"photo_url"`
	Description      sql.NullString `json:"description"`
	IsVerified       sql.NullBool   `json:"is_verified"`
	IsBlocked        sql.NullBool   `json:"is_blocked"`
}
