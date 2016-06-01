package main

import (
	"fmt"
	"model/DbModel/StaffModel"

	"github.com/tealeg/xlsx"
)

func main() {
	excelFileName := "/Users/tianyi/Desktop/eight_5.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf(err.Error())
		//fmt.Printf("%s\n", "error")
	}
	model := StaffModel.NewObj()
	fmt.Println((*model).TabName)
	for _, sheet := range xlFile.Sheets {
		//fmt.Println(len(sheet.Rows))
		for _, row := range sheet.Rows {
			tmp := make(map[string]string)
			len := len(row.Cells)
			if len < 3 {
				continue
			}
			tmp["name"] = "hello"
			tmp["id_no"] = row.Cells[2].Value
			tmp["mobilephone"] = row.Cells[1].Value
			model.Insert(tmp)
			break
		}
		break
	}
	fmt.Println("done")
}
