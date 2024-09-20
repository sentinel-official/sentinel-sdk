package options

import (
	"encoding/base64"
	"errors"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// Constants for the Page struct fields.
const (
	NamePageCountTotal = "PageCountTotal"
	NamePageKey        = "PageKey"
	NamePageLimit      = "PageLimit"
	NamePageOffset     = "PageOffset"
	NamePageReverse    = "PageReverse"
)

// Default values for the Page fields.
const (
	DefaultPageCountTotal = false
	DefaultPageKey        = ""
	DefaultPageLimit      = 10
	DefaultPageOffset     = 0
	DefaultPageReverse    = false
)

// Flags for command-line options for Page.
const (
	FlagPageCountTotal = "page.count-total"
	FlagPageKey        = "page.key"
	FlagPageLimit      = "page.limit"
	FlagPageOffset     = "page.offset"
	FlagPageReverse    = "page.reverse"
)

// init function sets the default values for page-related parameters at package initialization.
func init() {
	SetDefault(NamePageCountTotal, DefaultPageCountTotal)
	SetDefault(NamePageKey, DefaultPageKey)
	SetDefault(NamePageLimit, DefaultPageLimit)
	SetDefault(NamePageOffset, DefaultPageOffset)
	SetDefault(NamePageReverse, DefaultPageReverse)
}

// Page defines a structure for holding pagination-related parameters.
type Page struct {
	CountTotal bool   `json:"count_total" toml:"count_total"` // CountTotal indicates whether to include total count in paged queries.
	Key        string `json:"key" toml:"key"`                 // Key is the base64-encoded key for page.
	Limit      uint64 `json:"limit" toml:"limit"`             // Limit is the maximum number of results per page.
	Offset     uint64 `json:"offset" toml:"offset"`           // Offset is the offset for page.
	Reverse    bool   `json:"reverse" toml:"reverse"`         // Reverse indicates whether to reverse the order of results in page.
}

// NewPage creates a new Page instance with default values.
func NewPage() *Page {
	return &Page{
		CountTotal: cast.ToBool(GetDefault(NamePageCountTotal)),
		Key:        cast.ToString(GetDefault(NamePageKey)),
		Limit:      cast.ToUint64(GetDefault(NamePageLimit)),
		Offset:     cast.ToUint64(GetDefault(NamePageOffset)),
		Reverse:    cast.ToBool(GetDefault(NamePageReverse)),
	}
}

// WithCountTotal sets whether to include total count in the page and returns the modified Page.
func (p *Page) WithCountTotal(v bool) *Page {
	p.CountTotal = v
	return p
}

// WithKey sets the base64-encoded key for the page and returns the modified Page.
func (p *Page) WithKey(v []byte) *Page {
	p.Key = base64.StdEncoding.EncodeToString(v)
	return p
}

// WithLimit sets the maximum number of results per page and returns the modified Page.
func (p *Page) WithLimit(v uint64) *Page {
	p.Limit = v
	return p
}

// WithOffset sets the offset for the page and returns the modified Page.
func (p *Page) WithOffset(v uint64) *Page {
	p.Offset = v
	return p
}

// WithReverse sets whether to reverse the order of results in the page and returns the modified Page.
func (p *Page) WithReverse(v bool) *Page {
	p.Reverse = v
	return p
}

// GetCountTotal returns whether to include total count in the page.
func (p *Page) GetCountTotal() bool {
	return p.CountTotal
}

// GetKey returns the decoded base64 key for the page.
func (p *Page) GetKey() []byte {
	v, err := base64.StdEncoding.DecodeString(p.Key)
	if err != nil {
		panic(err)
	}

	return v
}

// GetLimit returns the maximum number of results per page.
func (p *Page) GetLimit() uint64 {
	return p.Limit
}

// GetOffset returns the offset for the page.
func (p *Page) GetOffset() uint64 {
	return p.Offset
}

// GetReverse returns whether the order of results is reversed in the page.
func (p *Page) GetReverse() bool {
	return p.Reverse
}

// ValidatePageKey checks if the provided key is a valid base64-encoded string.
func ValidatePageKey(v string) error {
	if _, err := base64.StdEncoding.DecodeString(v); err != nil {
		return errors.New("key must be a valid base64-encoded string")
	}

	return nil
}

// ValidatePageLimit validates the Limit field.
func ValidatePageLimit(v uint64) error {
	if v == 0 {
		return errors.New("limit must be greater than zero")
	}

	return nil
}

// Validate validates all the fields of the Page struct.
func (p *Page) Validate() error {
	if err := ValidatePageKey(p.Key); err != nil {
		return err
	}
	if err := ValidatePageLimit(p.Limit); err != nil {
		return err
	}

	return nil
}

// PageRequest creates a new PageRequest with the configured options.
func (p *Page) PageRequest() *query.PageRequest {
	return &query.PageRequest{
		Key:        p.GetKey(),
		Offset:     p.GetOffset(),
		Limit:      p.GetLimit(),
		CountTotal: p.GetCountTotal(),
		Reverse:    p.GetReverse(),
	}
}

// GetPageCountTotalFromCmd retrieves whether to include total count from the command-line flags.
func GetPageCountTotalFromCmd(cmd *cobra.Command) (bool, error) {
	return cmd.Flags().GetBool(FlagPageCountTotal)
}

// GetPageKeyFromCmd retrieves the base64-encoded key from the command-line flags.
func GetPageKeyFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagPageKey)
}

// GetPageLimitFromCmd retrieves the limit from the command-line flags.
func GetPageLimitFromCmd(cmd *cobra.Command) (uint64, error) {
	return cmd.Flags().GetUint64(FlagPageLimit)
}

// GetPageOffsetFromCmd retrieves the offset from the command-line flags.
func GetPageOffsetFromCmd(cmd *cobra.Command) (uint64, error) {
	return cmd.Flags().GetUint64(FlagPageOffset)
}

// GetPageReverseFromCmd retrieves whether to reverse the order of results from the command-line flags.
func GetPageReverseFromCmd(cmd *cobra.Command) (bool, error) {
	return cmd.Flags().GetBool(FlagPageReverse)
}

// SetFlagPageCountTotal sets the flag for the count total field in the given command.
func SetFlagPageCountTotal(cmd *cobra.Command) {
	value := cast.ToBool(GetDefault(NamePageCountTotal))
	cmd.Flags().Bool(FlagPageCountTotal, value, "Include total count in paged queries.")
}

// SetFlagPageKey sets the flag for the base64-encoded key in the given command.
func SetFlagPageKey(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NamePageKey))
	cmd.Flags().String(FlagPageKey, value, "Base64-encoded key for page.")
}

// SetFlagPageLimit sets the flag for the results per page limit in the given command.
func SetFlagPageLimit(cmd *cobra.Command) {
	value := cast.ToUint64(GetDefault(NamePageLimit))
	cmd.Flags().Uint64(FlagPageLimit, value, "Maximum number of results per page.")
}

// SetFlagPageOffset sets the flag for the page offset in the given command.
func SetFlagPageOffset(cmd *cobra.Command) {
	value := cast.ToUint64(GetDefault(NamePageOffset))
	cmd.Flags().Uint64(FlagPageOffset, value, "Offset for the page.")
}

// SetFlagPageReverse sets the flag for the reverse order in the given command.
func SetFlagPageReverse(cmd *cobra.Command) {
	value := cast.ToBool(GetDefault(NamePageReverse))
	cmd.Flags().Bool(FlagPageReverse, value, "Reverse the order of results in the page.")
}

// SetPageFlags sets all page-related flags for the command.
func SetPageFlags(cmd *cobra.Command) {
	SetFlagPageCountTotal(cmd)
	SetFlagPageKey(cmd)
	SetFlagPageLimit(cmd)
	SetFlagPageOffset(cmd)
	SetFlagPageReverse(cmd)
}

// NewPageFromCmd creates a new Page object from the command-line flags.
func NewPageFromCmd(cmd *cobra.Command) (*Page, error) {
	// Retrieve count total from the command flags
	countTotal, err := GetPageCountTotalFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve key from the command flags
	key, err := GetPageKeyFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve limit from the command flags
	limit, err := GetPageLimitFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve offset from the command flags
	offset, err := GetPageOffsetFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve reverse from the command flags
	reverse, err := GetPageReverseFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Return a new Page object with the retrieved values
	return &Page{
		CountTotal: countTotal,
		Key:        key,
		Limit:      limit,
		Offset:     offset,
		Reverse:    reverse,
	}, nil
}
