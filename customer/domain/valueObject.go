package domain

type (
	Phone struct {
		number string
	}
)

func NewPhone(number string) Phone {
	return Phone{number: number}
}

func (phone Phone) Number() string {
	return phone.number
}
