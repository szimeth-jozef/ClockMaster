package period

import (
	"fmt"
	"time"
)

type InvoicePeriod struct {
	Year  int
	Month int
}

func New() InvoicePeriod {
	now := time.Now()

	return InvoicePeriod{
		Year:  now.Year(),
		Month: int(now.Month()),
	}
}

func FromPart(year, month int) (InvoicePeriod, error) {
	if year < 0 || month < 1 || month > 12 {
		return InvoicePeriod{}, fmt.Errorf("invalid year or month")
	}

	return InvoicePeriod{
		Year:  year,
		Month: month,
	}, nil
}

func FromString(period string) (InvoicePeriod, error) {
	var year, month int
	_, err := fmt.Sscanf(period, "%d-%02d", &year, &month)
	if err != nil {
		return InvoicePeriod{}, err
	}

	return FromPart(year, month)
}

func (p InvoicePeriod) Copy() InvoicePeriod {
	return InvoicePeriod{
		Year:  p.Year,
		Month: p.Month,
	}
}

func (p *InvoicePeriod) Next() {
	if p.Month == 12 {
		p.Year++
		p.Month = 1
	} else {
		p.Month++
	}
}

func (p *InvoicePeriod) GetNext() InvoicePeriod {
	copy := p.Copy()
	copy.Next()
	return copy
}

func (p *InvoicePeriod) Previous() {
	if p.Month == 1 {
		p.Year--
		p.Month = 12
	} else {
		p.Month--
	}
}

func (p *InvoicePeriod) GetPrevious() InvoicePeriod {
	copy := p.Copy()
	copy.Previous()
	return copy
}

func (p InvoicePeriod) String() string {
	return fmt.Sprintf("%d-%02d", p.Year, p.Month)
}
