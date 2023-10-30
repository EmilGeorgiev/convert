package convert_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/EmilGeorgiev/convert"
	"github.com/stretchr/testify/assert"
)

func TestConvertSrcWithPtrFieldsToDstWithPtrFields(t *testing.T) {
	id := int64(22)
	name := "ivan"
	isStarred := true
	ownerType := "user"
	userName := "Emil"

	p := Template{
		ID:         &id,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  &ownerType,
		User:       &User{FirstName: "Emil"},
	}

	id2 := int64(999)
	got := &TemplateEntity{
		ID: &id2,
	}
	convert.SrcToDst(p, &got)

	idZero := int64(0)
	want := TemplateEntity{
		ID:         &id,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  nil,
		User:       &UserEntity{ID: &idZero, FirstName: &userName},
	}

	assert.Equal(t, want, *got)
}

func TestConvertSrcWithPtrFieldsToDstWithNonPtrFields(t *testing.T) {
	id := int64(22)
	name := "ivan"
	isStarred := true
	ownerType := int64(1)
	userName := "Emil"

	p := TemplateEntity{
		ID:         &id,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  &ownerType,
		User:       &UserEntity{FirstName: &userName},
	}

	var got NewTemplate
	convert.SrcToDst(p, &got)

	want := NewTemplate{
		ID:         id,
		Name:       name,
		Categories: []string{"11", "22"},
		IsStarred:  isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  1,
		User:       UserEntity{FirstName: &userName},
	}

	assert.Equal(t, want, got)
}

func TestConvertSrcWithNonPtrFieldsToDstWithPrtFields(t *testing.T) {
	id := int64(22)
	name := "ivan"
	isStarred := false

	p := NewTemplate{
		ID:         id,
		Name:       name,
		Categories: []string{"11", "22"},
		IsStarred:  isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  88123,
	}

	var got Template
	convert.SrcToDst(p, &got)

	want := Template{
		ID:         &id,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  nil,
		User:       &User{FirstName: ""},
	}

	assert.Equal(t, want, got)
}

func TestConvertSrcWithNonPtrFieldsToNDstWithNonPtrFields(t *testing.T) {
	name := "Emil"
	p := NewTemplate{
		ID:         int64(22),
		Name:       "ivan",
		Categories: []string{"11", "22"},
		IsStarred:  true,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  88123,
		User:       UserEntity{FirstName: &name},
	}

	var got NewTemplate
	convert.SrcToDst(p, &got)

	want := NewTemplate{
		ID:         int64(22),
		Name:       "ivan",
		Categories: []string{"11", "22"},
		IsStarred:  true,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  88123,
		User:       UserEntity{FirstName: &name},
	}

	assert.Equal(t, want, got)
}

func TestConvertPointerSrcToPointerDst(t *testing.T) {
	id := int64(22)
	name := "ivan"
	isStarred := true
	ownerType := "user"

	p := Template{
		ID:         &id,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  &ownerType,
	}

	var got TemplateEntity
	convert.SrcToDst(&p, &got)

	want := TemplateEntity{
		ID:         &id,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  nil,
		//User: &UserEntity{},
	}

	assert.Equal(t, want, got)
}

func TestConvertSrcWithNotSetTimeFieldsToDst(t *testing.T) {
	id := int64(22)
	name := "ivan"
	isStarred := true
	ownerType := "apikey"

	p := Template{
		ID:         &id,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		OwnerType:  &ownerType,
	}

	var got TemplateEntity
	convert.SrcToDst(p, &got)

	want := TemplateEntity{
		ID:         &id,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		OwnerType:  nil,
		//User: &UserEntity{},
	}

	assert.Equal(t, want, got)
	assert.Zero(t, got.CreatedAt)
}

func TestConvertSrcDstInterface(t *testing.T) {
	id := int64(22)
	name := "ivan"
	isStarred := true
	ownerType := "apikey"
	userName := "Emil"
	userID := int64(7878)

	p := Template{
		ID:         &id,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  &ownerType,
		User:       &User{ID: userID, FirstName: userName},
	}

	var got Printer
	got = &TemplateEntity{}
	convert.SrcToDst(p, got)

	want := &TemplateEntity{
		ID:         &id,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  nil,
		User:       &UserEntity{ID: &userID, FirstName: &userName},
	}

	assert.Equal(t, want, got)
}

func TestConvertSrcDstDifferentTypes(t *testing.T) {
	u := User{
		Permissions:       []uint{1, 2},
		SimplePermissions: []uint{10, 11},
	}

	var got UserEntity
	convert.SrcToDst(u, &got)

	id := int64(0)
	name := ""
	want := UserEntity{
		ID:                &id,
		FirstName:         &name,
		Permissions:       []uint{1, 2},
		SimplePermissions: []uint{10, 11},
	}

	assert.Equal(t, want, got)
}

func TestConvertPointerZeroValueDstPtrFields(t *testing.T) {
	name := "ivan"
	isStarred := true
	ownerType := "apikey"

	p := Template{
		ID:         nil,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  &ownerType,
	}

	var got TemplateEntity
	convert.SrcToDst(p, &got)

	want := TemplateEntity{
		ID:         nil,
		Name:       &name,
		Categories: []string{"11", "22"},
		IsStarred:  &isStarred,
		CreatedAt:  time.Unix(11111, 0),
		OwnerType:  nil,
	}

	assert.Equal(t, want, got)
}

func TestConvertSrcWithNilFieldsToDstWithNotNilFields(t *testing.T) {
	id := int64(22)
	name := "ivan"
	isStarred := true
	userName := "Emil"

	p := Template{
		ID:         nil,
		Name:       nil,
		Categories: nil,
		IsStarred:  nil,
		OwnerType:  nil,
		User:       nil,
	}

	got := TemplateEntity{
		ID:         &id,
		Name:       &name,
		Categories: []string{"full", "basic"},
		IsStarred:  &isStarred,
		User:       &UserEntity{FirstName: &userName},
	}
	convert.SrcToDst(p, &got)

	want := TemplateEntity{
		ID:         &id,
		Name:       &name,
		Categories: []string{"full", "basic"},
		IsStarred:  &isStarred,
		User:       &UserEntity{FirstName: &userName},
	}

	assert.Equal(t, want, got)
}

func TestConvertSrcWithEmptyFieldSliceToDstWithNotEmptyFieldsSlice(t *testing.T) {
	p := Template{
		Categories: []string{},
	}

	got := TemplateEntity{
		Categories: []string{"full", "basic"},
	}
	convert.SrcToDst(p, &got)

	want := TemplateEntity{
		Categories: []string{},
	}

	assert.Equal(t, want, got)
}

func TestConvertSrcWithEmptyFieldSliceToDstWithNilFieldsSlice(t *testing.T) {
	p := Template{
		Categories: []string{},
	}

	got := TemplateEntity{
		Categories: nil,
	}
	convert.SrcToDst(p, &got)

	want := TemplateEntity{
		Categories: []string{},
	}

	assert.Equal(t, want, got)
}

func TestConvertSrcWithNilSliceToDstWithNotEmptyFieldsSlice(t *testing.T) {
	p := Template{
		Categories: nil,
	}

	got := TemplateEntity{
		Categories: []string{"full", "basic"},
	}
	convert.SrcToDst(p, &got)

	want := TemplateEntity{
		Categories: []string{"full", "basic"},
	}

	assert.Equal(t, want, got)
}

func TestConvertAliasOfPrimitiveTypes(t *testing.T) {
	i := MyInt(56)
	var got MyIntEntity
	convert.SrcToDst(i, &got)

	want := MyIntEntity(56)
	assert.Equal(t, want, got)
}

func TestConvertPrimitiveType(t *testing.T) {
	i := 100
	var got int64
	convert.SrcToDst(i, &got)

	assert.EqualValues(t, 100, got)
}

func TestConvertPrimitiveOfDifferentTypes(t *testing.T) {
	i := "100"
	var got int64
	convert.SrcToDst(i, &got)
	assert.EqualValues(t, 0, got)
}

func TestConvertSrcWithAliasToDstWithoutAlias(t *testing.T) {
	ui := UUID("eee")
	b := Book{
		ID:     UUID("ooo"),
		Pages:  Pages(34),
		UserID: &ui,
	}

	var actual BookEntity
	convert.SrcToDst(b, &actual)

	expected := BookEntity{
		ID:     "ooo",
		Pages:  34,
		UserID: "eee",
	}
	assert.Equal(t, expected, actual)
}

func TestConvertSrcWithoutAliasToDstWithAlias(t *testing.T) {
	b := BookEntity{
		ID:     "ooo",
		Pages:  34,
		UserID: "eee",
	}

	var actual Book
	convert.SrcToDst(b, &actual)

	ui := UUID("eee")
	expected := Book{
		ID:     UUID("ooo"),
		Pages:  Pages(34),
		UserID: &ui,
	}
	assert.Equal(t, expected, actual)
}

func TestConvertSrcWithoutAliasToDstWithSliceOfAlias(t *testing.T) {
	b := BookEntity{
		ID:     "ooo",
		Pages:  34,
		UserID: "eee",
	}

	var actual Book
	convert.SrcToDst(b, &actual)

	ui := UUID("eee")
	expected := Book{
		ID:     UUID("ooo"),
		Pages:  Pages(34),
		UserID: &ui,
	}
	assert.Equal(t, expected, actual)
}

func TestConvertSrcWithAliasSliceToDstWithSliceWithoutAlias(t *testing.T) {
	q := Query{Filters: []string{"id", "name", "date"}}

	var actual QueryEntity
	convert.SrcToDst(q, &actual)

	expected := QueryEntity{Filters: []string{"id", "name", "date"}}

	assert.Equal(t, actual, expected)
}

func TestConvertSrcWithSlicesAToDstWithSliceAndAlias(t *testing.T) {
	q := QueryEntity{Filters: []string{"id", "name", "date"}}

	var actual Query
	convert.SrcToDst(q, &actual)

	expected := Query{Filters: []string{"id", "name", "date"}}

	assert.Equal(t, actual, expected)
}

//func TestConvertSrcWithSliceOfStructToDstWithSliceOfStruct(t *testing.T) {
//	src := MyTemplates{
//		Templates: []Template{
//			{Categories: []string{"111", "222"}},
//			{Categories: []string{"333", "444"}},
//		},
//	}
//
//	var actual MyNewTemplates
//	convert.SrcToDst(src, &actual)
//
//	expected := MyNewTemplates{
//		Templates: []NewTemplate{
//			{Categories: []string{"111", "222"}},
//			{Categories: []string{"333", "444"}},
//		},
//	}
//
//	assert.Equal(t, expected, actual)
//}

type MyTemplates struct {
	Templates []Template
}

type MyNewTemplates struct {
	Templates []NewTemplate
}

type User struct {
	ID                int64
	FirstName         string
	SimplePermissions []uint
	Permissions       Permissions
}

type UserEntity struct {
	ID                *int64
	FirstName         *string
	Permissions       Permissions
	SimplePermissions []uint
}

type Template struct {
	ID         *int64
	Name       *string
	Categories []string
	IsStarred  *bool
	CreatedAt  time.Time
	OwnerType  *string
	User       *User
}

type NewTemplate struct {
	ID         int64
	ExternalID string
	Name       string
	Categories []string
	IsStarred  bool
	CreatedAt  time.Time
	OwnerType  int64
	User       UserEntity
}

type TemplateEntity struct {
	ID         *int64
	Name       *string
	Categories []string
	IsStarred  *bool
	CreatedAt  time.Time
	OwnerType  *int64
	User       *UserEntity
}

type Permissions []uint

func (e TemplateEntity) Print() {
	fmt.Println("hello")
}

type Printer interface {
	Print()
}

type MyInt int64

type MyIntEntity int64

type UUID string

type Pages int

type Book struct {
	ID     UUID
	Pages  Pages
	UserID *UUID
}

type BookEntity struct {
	ID     string
	Pages  int
	UserID string
}

type Filters []string

type Query struct {
	Filters Filters
}

type QueryEntity struct {
	Filters []string
}
