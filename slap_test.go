package slap

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestCrud(t *testing.T) {
	type some struct {
		ID       string
		Address  string `slap:"index"`
		Name     string
		Universe int64
		Age      int `slap:"index"`
		Life     bool
		Range    []byte
		Money    float64 `slap:"index"`
		When     time.Time
	}

	tm := time.Now().Round(0)

	tbl1 := some{
		Address:  "St Leonards",
		Name:     "Jim",
		Universe: 424242,
		Age:      60,
		Life:     true,
		Range:    []byte("some bytes"),
		Money:    32.42,
		When:     tm,
	}

	tbl2 := some{
		Address:  "St Leonards",
		Name:     "Tom",
		Universe: 999,
		Age:      46,
		Life:     true,
		Range:    []byte("some bytes"),
		Money:    36.06,
	}

	tbl3 := some{
		Address:  "Jersey St",
		Universe: 1000,
		Age:      25,
		Life:     false,
		Range:    []byte("more bytes"),
		Money:    0.42,
	}

	tbl4 := some{
		Address:  "Romsey St",
		Universe: 1001,
		Age:      46,
		Life:     true,
		Range:    []byte("if any bytes"),
		Money:    100.01,
	}

	piv := New("/tmp/badger", "sparkle")
	defer piv.db.Close()

	sl := []some{tbl1, tbl2, tbl3, tbl4}
	ws := []string{"one", "two"}
	var err error

	t.Run("test create", func(t *testing.T) {
		piv.db.DropAll()

		_, err = piv.Create(&tbl1)
		if err != nil {
			t.Error(err)
		}
		_, err = piv.Create(tbl1)
		if err == nil {
			t.Error("must return error")
		}
		_, err = piv.Create(sl)
		if err == nil {
			t.Error("must return error")
		}
		_, err = piv.Create(ws)
		if err == nil {
			t.Error("must return error")
		}
		_, err = piv.Create(&ws)
		if err == nil {
			t.Error("must return error")
		}
		_, err = piv.Create(&sl)
		if err != nil {
			t.Error(err)
		}
		_, err = piv.Create("test")
		if err == nil {
			t.Error("must return error")
		}
	})

	t.Run("test read", func(t *testing.T) {
		piv.db.DropAll()

		id, err := piv.Create(&tbl1)
		if err != nil {
			t.Fatal(err)
		}
		if id == nil {
			t.Fatal("id should not be nil")
		}
		if len(id) != 1 {
			t.Fatal("id should have 1 element")
		}

		res, err := piv.Read(&some{}, []string{}, id...)
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("res should not be nil")
		}
		if len(res) != 1 {
			t.Fatal("res should have 1 element")
		}

		if res[0].(some).Address != "St Leonards" {
			t.Error("invalid read")
		}

		if res[0].(some).Name != "Jim" {
			t.Error("invalid read")
		}

		if res[0].(some).Universe != 424242 {
			t.Error("invalid read")
		}

		if res[0].(some).When != tm {
			t.Error("invalid read")
		}

		if res[0].(some).ID != id[0] {
			t.Error("invalid read")
		}

		if res[0].(some).Age != 60 {
			t.Error("invalid read")
		}

		if res[0].(some).Life != true {
			t.Error("invalid read")
		}

		if string(res[0].(some).Range) != "some bytes" {
			t.Error("invalid read")
		}

		if res[0].(some).Money != 32.42 {
			t.Error("invalid read")
		}

		res, err = piv.Read(&some{}, []string{"Name", "Address"}, id...)
		if err != nil {
			t.Fatal(err)
		}

		if len(res) != 1 {
			t.Fatal("res should have 1 element")
		}

		if res[0].(some).Address != "St Leonards" {
			t.Error("invalid read")
		}

		if res[0].(some).Name != "Jim" {
			t.Error("invalid read")
		}

		if res[0].(some).Universe != 0 {
			t.Error("invalid read")
		}

		if res[0].(some).When == tm {
			t.Error("invalid read")
		}

		if res[0].(some).Age != 0 {
			t.Error("invalid read")
		}

		if res[0].(some).Life != false {
			t.Error("invalid read")
		}

		if string(res[0].(some).Range) != "" {
			t.Error("invalid read")
		}

		if res[0].(some).Money != 0.0 {
			t.Error("invalid read")
		}

		if res[0].(some).ID != id[0] {
			t.Error("invalid read")
		}

		ids, err := piv.Create(&sl)
		if err != nil {
			t.Error(err)
		}

		res2, err := piv.Read(&some{}, []string{}, ids...)
		if err != nil {
			t.Fatal(err)
		}

		if len(res2) != 4 {
			t.Error("invalid read")
		}
	})

	t.Run("test update", func(t *testing.T) {
		piv.db.DropAll()

		id, err := piv.Create(&tbl4)
		if err != nil {
			t.Fatal(err)
		}

		res, err := piv.Read(&some{}, []string{}, id[0])
		if err != nil {
			t.Fatal(err)
		}

		if res[0].(some).Name != "" {
			t.Error("invalid field read")
		}

		if res[0].(some).Address != "Romsey St" {
			t.Error("invalid field read")
		}

		if res[0].(some).ID != id[0] {
			t.Error("invalid field read")
		}

		err = piv.Update(&some{Name: "Ruslan", Address: "Jersey St", ID: "blah"}, id[0])
		if err != nil {
			t.Error(err)
		}

		res, err = piv.Read(&some{}, []string{}, id[0])
		if err != nil {
			t.Fatal(err)
		}

		if res[0].(some).Name != "Ruslan" {
			t.Error("invalid field update")
		}

		if res[0].(some).Address != "Jersey St" {
			t.Error("invalid field update")
		}

		if res[0].(some).ID != id[0] {
			t.Error("invalid field update")
		}
	})

	t.Run("test delete", func(t *testing.T) {
		piv.db.DropAll()

		id, err := piv.Create(&tbl3)
		if err != nil {
			t.Fatal(err)
		}

		res, err := piv.Read(&some{}, []string{}, id[0])
		if err != nil {
			t.Fatal(err)
		}

		if res[0].(some).Age != 25 {
			t.Error("invalid field read")
		}

		err = piv.Delete(&some{}, id[0])
		if err != nil {
			t.Fatal(err)
		}

		res, err = piv.Read(&some{}, []string{}, id[0])
		if !errors.Is(err, ErrNoRecord) {
			t.Fatal("should display correct error")
		}
		if res == nil {
			t.Fatal("result should not be nil")
		}
		if len(res) != 0 {
			t.Fatal("res should have 0 element")
		}
	})

	t.Run("test model", func(t *testing.T) {

		m, err := model(&tbl1, true)
		if err != nil {
			t.Fatal(err)
		}
		v, err := m.values(&tbl1)
		if err != nil {
			t.Error(err)
		}
		if v["Address"].(string) != "St Leonards" {
			t.Error("value conversion")
		}
		if v["Money"].(float64) != 32.42 {
			t.Error("value conversion")
		}

		m, err = model(&tbl3, false)
		if err != nil {
			t.Fatal(err)
		}
		v, err = m.values(&tbl3)
		if err != nil {
			t.Error(err)
		}
		if _, ok := v["Name"]; ok {
			t.Error("zero value present")
		}

		m, err = model(&tbl4, true)
		if err != nil {
			t.Fatal(err)
		}
		v, err = m.values(&tbl4)
		if err != nil {
			t.Error(err)
		}
		if v["Name"].(string) != "" {
			t.Error("value conversion")
		}

		type s1 struct{ Name string }
		m, err = model(&s1{}, true)
		if !errors.Is(err, ErrNoPrimaryID) {
			t.Fatal("must retuen correct error")
		}
	})

	t.Run("test where", func(t *testing.T) {
		piv.db.DropAll()

		_, err = piv.Create(&sl)
		if err != nil {
			t.Error(err)
		}

		res, err := piv.where(some{Address: "Romsey St"})

		rd, err := piv.Read(&some{}, []string{}, res...)
		if err != nil {
			t.Fatal(err)
		}
		if rd == nil {
			t.Fatal("res should not be nil")
		}

		if rd[0].(some).Money != 100.01 {
			t.Error("invalid read field")
		}

	})

	t.Run("test select", func(t *testing.T) {
		piv.db.DropAll()

		_, err := piv.Create(&sl)
		if err != nil {
			t.Error(err)
		}

		res, err := piv.Select(&some{Address: "St Leonards"}, []string{})
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("res should not be nil")
		}
		if len(res) != 2 {
			t.Fatal("res should have 2 elements")
		}

		res, err = piv.Select(&some{Address: "St Leonards", Age: 46}, []string{})
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("res should not be nil")
		}
		if len(res) != 1 {
			t.Fatal("res should have 1 elements")
		}

		res, err = piv.Select(&some{Address: "Romsey St"}, []string{})
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("res should not be nil")
		}
		if len(res) != 1 {
			t.Fatal("res should have 1 elements")
		}

	})
}

func TestEncoding(t *testing.T) {
	ss := "Hello, World"
	v, err := toBytes(ss)
	if err != nil {
		t.Error(err)
	}
	r, err := fromBytes(v, "string")
	if err != nil {
		t.Error(err)
	}
	if ss != r.(string) {
		t.Error("invalid conversion")
	}
	si := 42
	v, err = toBytes(si)
	if err != nil {
		t.Error(err)
	}
	r, err = fromBytes(v, "int")
	if err != nil {
		t.Error(err)
	}
	if si != r.(int) {
		t.Error("invalid conversion")
	}
	bl := true
	v, err = toBytes(bl)
	if err != nil {
		t.Error(err)
	}
	r, err = fromBytes(v, "bool")
	if err != nil {
		t.Error(err)
	}
	if bl != r.(bool) {
		t.Error("invalid conversion")
	}
	bs := []byte("some bytes")
	v, err = toBytes(bs)
	if err != nil {
		t.Error(err)
	}
	r, err = fromBytes(v, "[]uint8")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(bs, r.([]byte)) {
		t.Error("invalid conversion")
	}
	fl := 32.54
	v, err = toBytes(fl)
	if err != nil {
		t.Error(err)
	}
	r, err = fromBytes(v, "float64")
	if err != nil {
		t.Error(err)
	}
	if fl != r.(float64) {
		t.Error("invalid conversion")
	}
}

func TestTime(t *testing.T) {
	piv := New("/tmp/badger", "sparkle")
	defer piv.db.Close()
	piv.db.DropAll()
	w := time.Now().Round(0)

	type tmc struct {
		ID   string
		When time.Time
	}

	b := tmc{When: w}

	id, err := piv.Create(&b)
	if err != nil {
		t.Error(err)
	}

	res, err := piv.Read(&tmc{}, []string{}, id...)
	if err != nil {
		t.Fatal(err)
	}

	if res[0].(tmc).When != w {
		t.Error("time should match")
	}
}

func TestTake(t *testing.T) {
	piv := New("/tmp/badger", "sparkle")
	defer piv.db.Close()
	piv.db.DropAll()

	type tmc struct {
		ID    string
		Name  string
		Age   int
		Count int
	}

	arr := []tmc{}
	for i := 1; i < 6; i++ {
		b := tmc{Name: "Ruslan", Age: 46, Count: i}
		arr = append(arr, b)
	}

	_, err := piv.Create(&arr)
	if err != nil {
		t.Error(err)
	}

	var res []interface{}
	for i := 1; i < 6; i++ {

		res, err = piv.Take(&tmc{}, []string{"Name", "Count"}, "", i)

		if err != nil {
			t.Error(err)
		}
		if len(res) != i {
			t.Error("wrong limit return")
		}

		for i, r := range res {
			s := r.(tmc)
			if s.Name != "Ruslan" {
				t.Error("conversion error")
			}
			if s.Count != i+1 {
				t.Error("limit error")
			}
			if s.Age != 0 {
				t.Error("conversion error")
			}
		}
	}

	limit := 2
	id := res[limit].(tmc).ID
	res, err = piv.Take(&tmc{}, []string{}, id, limit)
	for i, r := range res {
		s := r.(tmc)
		if s.Name != "Ruslan" {
			t.Error("conversion error")
		}
		if s.Age != 46 {
			t.Error("conversion error")
		}
		if s.Count != i+limit+1 {
			t.Error("skip/limit error")
		}
	}
	t.Log(res)
}
