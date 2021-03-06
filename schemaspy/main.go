package main

import (
	"fmt"
	"reflect"
)

type VehicleInfo struct {
	// ID         bson.ObjectId `bson:"_id,omitempty"`
	VehicleId string `bson:"编号" json:"vid"`
	Date      int    `bson:"日期" json:"date"`
	Type      string `bson:"类型" json:"type"`
	Brand     string `bson:"型号" json:"brand"`
	Color     string `bson:"颜色" json:"color"`
}

func main() {
	vinfo := VehicleInfo{
		VehicleId: "123456",
		Date:      20140101,
		Type:      "Truck",
		Brand:     "Ford",
		Color:     "White",
	}
	vt := reflect.TypeOf(vinfo)
	vv := reflect.ValueOf(vinfo)
	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)
		//	chKey := f.Tag.Get("bson")
		chKey := f.Tag.Get("json")
		cnType := f.Type
		fmt.Printf("%q => --  %q, ", chKey, vv.FieldByName(f.Name).String(), cnType)
	}
}
