package question

// ReportReasons are the reasons accepted for a question report, shared by
// drill and exam.
var ReportReasons = map[string]bool{"kunci_salah": true, "soal_ambigu": true, "typo": true, "lainnya": true}

// MaxReportNote is the longest note (in characters) a report may carry.
const MaxReportNote = 500
