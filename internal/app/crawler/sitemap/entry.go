package sitemap

import "net/url"

type Frequency string

const (
	FrequencyDaily   Frequency = "daily"
	FrequencyMonthly Frequency = "monthly"
)

type Entry struct {
	Location        url.URL
	ChangeFrequency Frequency
	Priority        float64
}
