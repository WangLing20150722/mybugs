package works

import (
	"container/list"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
)

var defaultSheetName = "Sheet1"

func IssueOwnerList2Excel(l *list.List, file string) error {
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(defaultSheetName)
	xlsx.SetActiveSheet(index)

	xlsx.SetCellValue(defaultSheetName, "A1", "ID")
	xlsx.SetCellValue(defaultSheetName, "B1", "Level")
	xlsx.SetCellValue(defaultSheetName, "C1", "Summary")
	xlsx.SetCellValue(defaultSheetName, "D1", "Status")
	xlsx.SetCellValue(defaultSheetName, "E1", "LastModify")
	xlsx.SetCellValue(defaultSheetName, "F1", "LastAssignOutTo")
	xlsx.SetCellValue(defaultSheetName, "G1", "LastFix")

	line := 2
	for iter := l.Front(); iter != nil; iter = iter.Next() {
		issue := iter.Value.(*IssueOwner)

		xlsx.SetCellValue(defaultSheetName, "A"+strconv.Itoa(line), issue.Id)
		xlsx.SetCellHyperLink(defaultSheetName, "A"+strconv.Itoa(line), "http://mantis.tclking.com/view.php?id="+strconv.FormatInt(issue.Id, 10), "External")

		xlsx.SetCellValue(defaultSheetName, "B"+strconv.Itoa(line), issue.Level)
		xlsx.SetCellValue(defaultSheetName, "C"+strconv.Itoa(line), issue.Summary)
		xlsx.SetCellValue(defaultSheetName, "D"+strconv.Itoa(line), issue.Status)
		xlsx.SetCellValue(defaultSheetName, "E"+strconv.Itoa(line), issue.LastModify)
		xlsx.SetCellValue(defaultSheetName, "F"+strconv.Itoa(line), issue.LastAssignOutTo)
		xlsx.SetCellValue(defaultSheetName, "G"+strconv.Itoa(line), issue.LastFix)

		line++
	}

	// Save xlsx file by the given path.
	return xlsx.SaveAs(file)
}
