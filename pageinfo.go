package web

import (
	"net/http"
	"time"
)

// ChangeFrequency tells how frequently the page is likely to change
type ChangeFrequency int

const (
	ChangeFrequencyUnknown ChangeFrequency = iota
	ChangeFrequencyNever
	ChangeFrequencyAlways
	ChangeFrequencyHourly
	ChangeFrequencyDaily
	ChangeFrequencyWeekly
	ChangeFrequencyMonthly
	ChangeFrequencyYearly
)

func (f ChangeFrequency) String() string {
	s := []string{"", "never", "always", "hourly", "weekly", "monthly", "yearly"}
	n := int(f)

	if f > ChangeFrequencyUnknown &&
		n < len(s) {
		return s[n]
	}

	return ""
}

// Router can check for the existance of the requested resource
type RouterPageInfo interface {
	PageInfo(*http.Request) (interface{}, bool)
}

// PageInfo can tell how frequencly the page is likely to change
type PageInfoChangeFrequency interface {
	ChangeFrequency() string
}

// PageInfo can tell what Path was requested
type PageInfoLocation interface {
	Location() string
}

// PageInfo can tell the recommended Path to access this resource
type PageInfoCanonical interface {
	Canonical() string
}

// PageInfo can tell the relevance of this resource
type PageInfoPriority interface {
	Priority() float32
}

// PageInfo can tell when the package was last modified
type PageInfoLastModified interface {
	LastModified() time.Time
}

// PageInfo can tell the mime-types supported by the resource
type PageInfoMimeType interface {
	MimeType() []string
}

// PageInfo can tell the languages supported by the resource
type PageInfoLanguage interface {
	Language() []string
}

// PageInfo can tell the supported methods supported by the resource
type PageInfoMethods interface  {
	Method() []string
}

// PageInfo can tell what handler to use to request the resource
type PageInfoHandler() interface {
	Handler() http.Handler
}
