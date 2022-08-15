package domain

import (
	"testing"

	"happyday/abstract"
	"happyday/common"

	"github.com/stretchr/testify/assert"
)

func TestChangeName_Should_ReturnNameIsEmpty_When_NameIsEmpty(t *testing.T) {
	root := &aggregateRoot{}
	err := root.ChangeName("")

	assert.NotNil(t, err)
	assert.Equal(t, NameIsEmpty, err)

}

func TestChangeName_Should_ReturnNameTooLarge_When_NameLengthIsLargeThan255(t *testing.T) {
	root := &aggregateRoot{}
	err := root.ChangeName(common.RandString(256))

	assert.NotNil(t, err)
	assert.Equal(t, NameIsTooLarge, err)
}

func TestChangeName(t *testing.T) {
	root := &aggregateRoot{
		abstract.NewAggregateRoot(&State{}, 0),
	}

	name := common.RandString(100)
	err := root.ChangeName(name)

	assert.Nil(t, err)
	assert.Equal(t, name, root.State().Name())

	assert.Equal(t, 1, len(root.Events()))
	switch e := root.Events()[0].(type) {
	case NameChanged:
		assert.Equal(t, name, e.Name())
	default:
		assert.Fail(t, "invalid type")
	}
}

func TestChangeComment_Should_ReturnCommentLengthTooLarge_When_CommentLengthIsLargeThan1000(t *testing.T) {
	root := &aggregateRoot{}
	err := root.ChangeComment(common.RandString(1_001))

	assert.NotNil(t, err)
	assert.Equal(t, CommentTooLarge, err)
}

func TestChangeComment(t *testing.T) {
	root := &aggregateRoot{
		abstract.NewAggregateRoot(&State{}, 0),
	}

	comment := common.RandString(100)
	err := root.ChangeComment(comment)

	assert.Nil(t, err)
	assert.Equal(t, comment, root.State().Comment())
	assert.Equal(t, 1, len(root.Events()))
	switch e := root.Events()[0].(type) {
	case CommentChanged:
		assert.Equal(t, comment, e.Comment())
	default:
		assert.Fail(t, "invalid type")
	}
}

func TestChangePhones_Should_ReturnPhoneIsEmpty_When_HasNoPhone(t *testing.T) {
	root := &aggregateRoot{}
	err := root.ChangePhones([]Phone{})

	assert.NotNil(t, err)
	assert.Equal(t, PhonesIsEmpty, err)
}

func TestChangePhones_Should_ReturnPhoneLengthIsInvalid_When_PhoneNumberIsLessThan9(t *testing.T) {
	root := &aggregateRoot{}
	err := root.ChangePhones([]Phone{
		NewPhone(common.RandString(8)),
	})

	assert.NotNil(t, err)
	assert.Equal(t, PhoneLengthIsInvalid, err)
}

func TestChangePhones_Should_ReturnPhoneLengthIsInvalid_When_PhoneNumberIsGreaterThan11(t *testing.T) {
	root := &aggregateRoot{}
	err := root.ChangePhones([]Phone{
		NewPhone(common.RandString(12)),
	})

	assert.NotNil(t, err)
	assert.Equal(t, PhoneLengthIsInvalid, err)
}

func TestChangePhones_Should_ReturnPhoneIsInvalid_When_PhoneNumberIsInvalid(t *testing.T) {
	root := &aggregateRoot{}
	err := root.ChangePhones([]Phone{
		NewPhone(common.RandString(9)),
	})

	assert.NotNil(t, err)
	assert.Equal(t, PhoneNumberIsInvalid, err)
}

func TestChangePhones(t *testing.T) {
	root := &aggregateRoot{
		abstract.NewAggregateRoot(&State{}, 0),
	}

	phones := []Phone{
		NewPhone("123456789"),
		NewPhone("13123456780"),
		NewPhone("123456781"),
	}

	err := root.ChangePhones(phones)

	assert.Nil(t, err)
	assert.Equal(t, phones, root.State().Phones())
	assert.Equal(t, 1, len(root.Events()))
	switch e := root.Events()[0].(type) {
	case PhonesChanges:
		assert.Equal(t, phones, e.Phones())
	default:
		assert.Fail(t, "invalid type")
	}
}

func TestCreate_Should_ReturnNameIsEmpty_When_NameIsEmpty(t *testing.T) {
	root := &aggregateRoot{}
	err := root.Create("", "", []Phone{})

	assert.NotNil(t, err)
	assert.Equal(t, NameIsEmpty, err)
}

func TestCreate_Should_ReturnNameTooLarge_When_NameLengthIsLargeThan255(t *testing.T) {
	root := &aggregateRoot{}
	err := root.Create(common.RandString(256), "", []Phone{})

	assert.NotNil(t, err)
	assert.Equal(t, NameIsTooLarge, err)
}

func TestCreate_Should_ReturnCommentLengthTooLarge_When_CommentLengthIsLargeThan1000(t *testing.T) {
	root := &aggregateRoot{}
	err := root.Create(common.RandString(10), common.RandString(1_001), []Phone{})

	assert.NotNil(t, err)
	assert.Equal(t, CommentTooLarge, err)
}

func TestCreate_Should_ReturnPhoneIsEmpty_When_HasNoPhone(t *testing.T) {
	root := &aggregateRoot{}
	err := root.Create(common.RandString(10), "", []Phone{})

	assert.NotNil(t, err)
	assert.Equal(t, PhonesIsEmpty, err)
}

func TestCreate_Should_ReturnPhoneLengthIsInvalid_When_PhoneNumberIsLessThan9(t *testing.T) {
	root := &aggregateRoot{}
	err := root.Create(common.RandString(10), "",
		[]Phone{
			NewPhone(common.RandString(8)),
		})

	assert.NotNil(t, err)
	assert.Equal(t, PhoneLengthIsInvalid, err)
}

func TestCreate_Should_ReturnPhoneLengthIsInvalid_When_PhoneNumberIsGreaterThan11(t *testing.T) {
	root := &aggregateRoot{}
	err := root.Create(common.RandString(10), "",
		[]Phone{
			NewPhone(common.RandString(12)),
		})

	assert.NotNil(t, err)
	assert.Equal(t, PhoneLengthIsInvalid, err)
}

func TestCreate_Should_ReturnPhoneIsInvalid_When_PhoneNumberIsInvalid(t *testing.T) {
	root := &aggregateRoot{}
	err := root.Create(common.RandString(10), "",
		[]Phone{
			NewPhone(common.RandString(9)),
		})

	assert.NotNil(t, err)
	assert.Equal(t, PhoneNumberIsInvalid, err)
}

func TestCreate(t *testing.T) {
	root := &aggregateRoot{
		abstract.NewAggregateRoot(&State{}, 0),
	}

	name := common.RandString(10)
	comment := common.RandString(100)
	phones := []Phone{
		{number: "123456789"},
		{number: "13123456789"},
	}
	err := root.Create(name, comment, phones)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(root.Events()))
	assert.Equal(t, name, root.State().Name())
	assert.Equal(t, comment, root.State().Comment())
	assert.Equal(t, phones, root.State().Phones())

	switch e := root.Events()[0].(type) {
	case Created:
		assert.Equal(t, name, e.Name())
		assert.Equal(t, comment, e.Comment())
		assert.Equal(t, phones, e.Phones())
	default:
		assert.Fail(t, "expecting Created")
	}
}
