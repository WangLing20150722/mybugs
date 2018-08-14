package works

import (
	"container/list"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"strconv"
)

var defaultSheetName = "sheet1"


func IssueOwnerList2Excel(l *list.List, file string,sheet string) error {
	if len(sheet)>0{
		defaultSheetName = sheet
	}
	xlsx,err:= excelize.OpenFile(file)
	if err != nil{
		xlsx = excelize.NewFile()
	}
	index := xlsx.NewSheet(defaultSheetName)
	xlsx.SetActiveSheet(index)

	xlsx.SetCellValue(defaultSheetName, "A1", "ID")
	xlsx.SetCellValue(defaultSheetName, "B1", "Project")
	xlsx.SetCellValue(defaultSheetName, "C1", "Level")
	xlsx.SetCellValue(defaultSheetName, "D1", "Summary")
	xlsx.SetCellValue(defaultSheetName, "E1", "Status")
	xlsx.SetCellValue(defaultSheetName, "F1", "LastModify")
	xlsx.SetCellValue(defaultSheetName, "G1", "LastAssignOutTo")
	xlsx.SetCellValue(defaultSheetName, "H1", "LastFix")
	xlsx.SetCellValue(defaultSheetName, "I1", "FirstInTime")
	xlsx.SetCellValue(defaultSheetName, "J1", "LastOutTime")
	xlsx.SetCellValue(defaultSheetName, "K1", "LastModifyTime")
	xlsx.SetCellValue(defaultSheetName, "L1", "Failed")

	line := 2
	for iter := l.Front(); iter != nil; iter = iter.Next() {
		issue := iter.Value.(*IssueOwner)

		if DEBUG {
			log.Print(issue)
		}

		xlsx.SetCellValue(defaultSheetName, "A"+strconv.Itoa(line), issue.Id)
		xlsx.SetCellHyperLink(defaultSheetName, "A"+strconv.Itoa(line), "http://mantis.tclking.com/view.php?id="+strconv.FormatInt(issue.Id, 10), "External")

		xlsx.SetCellValue(defaultSheetName, "B"+strconv.Itoa(line), issue.Project)
		xlsx.SetCellValue(defaultSheetName, "C"+strconv.Itoa(line), issue.Level)
		xlsx.SetCellValue(defaultSheetName, "D"+strconv.Itoa(line), issue.Summary)
		xlsx.SetCellValue(defaultSheetName, "E"+strconv.Itoa(line), issue.Status)
		xlsx.SetCellValue(defaultSheetName, "F"+strconv.Itoa(line), issue.LastModify)
		xlsx.SetCellValue(defaultSheetName, "G"+strconv.Itoa(line), issue.LastAssignOutTo)
		xlsx.SetCellValue(defaultSheetName, "H"+strconv.Itoa(line), issue.LastFix)

		xlsx.SetCellValue(defaultSheetName, "I"+strconv.Itoa(line), issue.InTime.Format("2006-01-02 15:04:05"))
		xlsx.SetCellValue(defaultSheetName, "J"+strconv.Itoa(line), issue.OutTime.Format("2006-01-02 15:04:05"))
		xlsx.SetCellValue(defaultSheetName, "K"+strconv.Itoa(line), issue.LastModifyTime.Format("2006-01-02 15:04:05"))
		if issue.Failed{
			xlsx.SetCellValue(defaultSheetName, "L"+strconv.Itoa(line), "Failed")

		}
		line++
	}

	// Save xlsx file by the given path.
	return xlsx.SaveAs(file)

}

func IssueOwnerFailedList2Excel(l *list.List, file string,sheetName string) error {

	xlsx,err:= excelize.OpenFile(file)
	if err != nil{
		xlsx = excelize.NewFile()
	}
	index := xlsx.NewSheet(sheetName)
	xlsx.SetActiveSheet(index)

	xlsx.SetCellValue(sheetName, "A1", "ID")
	xlsx.SetCellValue(sheetName, "B1", "Project")
	xlsx.SetCellValue(sheetName, "C1", "Level")
	xlsx.SetCellValue(sheetName, "D1", "Summary")
	xlsx.SetCellValue(sheetName, "E1", "Status")
	xlsx.SetCellValue(sheetName, "F1", "LastModify")
	xlsx.SetCellValue(sheetName, "G1", "LastAssignOutTo")
	xlsx.SetCellValue(sheetName, "H1", "LastFix")
	xlsx.SetCellValue(sheetName, "I1", "Failed")


	line := 2
	for iter := l.Front(); iter != nil; iter = iter.Next() {
		issue := iter.Value.(*IssueOwner)

		if DEBUG {
			log.Print(issue)
		}

		xlsx.SetCellValue(sheetName, "A"+strconv.Itoa(line), issue.Id)
		xlsx.SetCellHyperLink(sheetName, "A"+strconv.Itoa(line), "http://mantis.tclking.com/view.php?id="+strconv.FormatInt(issue.Id, 10), "External")

		xlsx.SetCellValue(sheetName, "B"+strconv.Itoa(line), issue.Project)
		xlsx.SetCellValue(sheetName, "C"+strconv.Itoa(line), issue.Level)
		xlsx.SetCellValue(sheetName, "D"+strconv.Itoa(line), issue.Summary)
		xlsx.SetCellValue(sheetName, "E"+strconv.Itoa(line), issue.Status)
		xlsx.SetCellValue(sheetName, "F"+strconv.Itoa(line), issue.LastModify)
		xlsx.SetCellValue(sheetName, "G"+strconv.Itoa(line), issue.LastAssignOutTo)
		xlsx.SetCellValue(sheetName, "H"+strconv.Itoa(line), issue.LastFix)
		//
		//xlsx.SetCellValue(sheetName, "I"+strconv.Itoa(line), issue.InTime.Format("2006-01-02 15:04:05"))
		//xlsx.SetCellValue(sheetName, "J"+strconv.Itoa(line), issue.OutTime.Format("2006-01-02 15:04:05"))
		//xlsx.SetCellValue(sheetName, "K"+strconv.Itoa(line), issue.LastModifyTime.Format("2006-01-02 15:04:05"))
		if issue.Failed{
			xlsx.SetCellValue(sheetName, "I"+strconv.Itoa(line), "Failed")

		}
		line++
	}

	// Save xlsx file by the given path.
	return xlsx.SaveAs(file)
}



