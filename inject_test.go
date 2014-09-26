package inject_test

import (
	"inject"
	"reflect"
	"testing"
	"time"
)

type Detail struct {
	Addr string
}

type testStruct struct {
	Detail
	Field string
	Name  string
	Prod  int
	Total int
	Tm    time.Time
}

func testFun(t time.Time, count, before int) int {
	after := count + before
	return after
}

func Test_Invoke(t *testing.T) {
	inj := inject.New()

	inj.MapIndex(0, time.Now())
	inj.MapIndex(1, 250)
	inj.MapIndex(2, 1000)

	vals, err := inj.Invoke(testFun)
	if err != nil {
		t.Errorf("Test_Invoke failed:", err)
	}

	rets := vals[0].Int()
	if rets != 1250 {
		t.Errorf("Test_Invoke failed: [1]return value incorrect")
	}

	inj.SetIndex(1, reflect.TypeOf(300), reflect.ValueOf(300))
	inj.SetIndex(2, reflect.TypeOf(2000), reflect.ValueOf(2000))
	vals, err = inj.Invoke(testFun)
	if err != nil {
		t.Errorf("Test_Invoke failed:", err)
	}

	rets = vals[0].Int()
	if rets != 2300 {
		t.Errorf("Test_Invoke failed: [2]return value incorrect")
	}
}

func Test_AssignField(t *testing.T) {
	st := &testStruct{Detail{}, "china", "beijing", 100, 500, time.Now()}

	inj := inject.New()

	tm := time.Now().Add(time.Hour)

	inj.MapTag("Field", "US")
	inj.MapTag("Name", "Huston")
	inj.MapTag("Prod", 300)
	inj.MapTag("Total", 4500)
	inj.MapTag("Tm", tm)

	inj.AssignField(st)

	if st.Addr != "" || st.Field != "US" || st.Name != "Huston" || st.Prod != 300 || st.Total != 4500 || st.Tm != tm {
		t.Errorf("Test_AssignField:[1] failed")
	}

	addr := Detail{"Time"}
	tm = time.Now().Add(time.Hour * 2)
	inj.SetTag("Name", reflect.TypeOf("New York"), reflect.ValueOf("New York"))
	inj.SetTag("Total", reflect.TypeOf(2000), reflect.ValueOf(2000))
	inj.SetTag("Tm", reflect.TypeOf(tm), reflect.ValueOf(tm))
	inj.SetTag("Detail", reflect.TypeOf(addr), reflect.ValueOf(addr))

	inj.AssignField(st)

	if st.Addr != "Time" || st.Field != "US" || st.Name != "New York" || st.Prod != 300 || st.Total != 2000 || st.Tm != tm {
		t.Errorf("Test_AssignField:[2] failed")
	}
}
