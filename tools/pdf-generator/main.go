package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type Patient struct {
	Name    string
	ID      string
	DOB     string
}

type LabResult struct {
	TestName      string
	Result        string
	ReferenceRange string
	Units         string
}

func generateRandomPatient() Patient {
	firstNames := []string{"John", "Jane", "Peter", "Mary", "David", "Susan"}
	lastNames := []string{"Doe", "Smith", "Jones", "Williams", "Brown", "Davis"}

	return Patient{
		Name:    fmt.Sprintf("%s %s", firstNames[rand.Intn(len(firstNames))], lastNames[rand.Intn(len(lastNames))]),
		ID:      fmt.Sprintf("%d", 100000000+rand.Intn(900000000)),
		DOB:     time.Date(1950+rand.Intn(50), time.Month(rand.Intn(12)+1), rand.Intn(28)+1, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
	}
}

func generateRandomLabResults() []LabResult {
	results := []LabResult{}

	// Glucose
	results = append(results, LabResult{
		TestName:      "Glucose",
		Result:        fmt.Sprintf("%d", 60+rand.Intn(60)),
		ReferenceRange: "65-99",
		Units:         "mg/dL",
	})

	// Hemoglobin A1c
	results = append(results, LabResult{
		TestName:      "Hemoglobin A1c",
		Result:        fmt.Sprintf("%.1f", 4.0+rand.Float64()*3.0),
		ReferenceRange: "4.8-5.6",
		Units:         "%",
	})

	// Cholesterol
	results = append(results, LabResult{
		TestName:      "Cholesterol, Total",
		Result:        fmt.Sprintf("%d", 150+rand.Intn(100)),
		ReferenceRange: "100-199",
		Units:         "mg/dL",
	})

	return results
}

func createLabResultsPDF(filename string) {
	patient := generateRandomPatient()
	labResults := generateRandomLabResults()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "LabCorp Clinical Laboratories")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Patient Name: %s", patient.Name))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Patient ID: %s", patient.ID))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Date of Birth: %s", patient.DOB))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Date of Collection: %s", time.Now().Format("2006-01-02")))
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Test Name")
	pdf.Cell(40, 10, "Result")
	pdf.Cell(40, 10, "Reference Range")
	pdf.Cell(40, 10, "Units")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	for _, res := range labResults {
		pdf.Cell(40, 10, res.TestName)
		pdf.Cell(40, 10, res.Result)
		pdf.Cell(40, 10, res.ReferenceRange)
		pdf.Cell(40, 10, res.Units)
		pdf.Ln(10)
	}

	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		fmt.Printf("Error creating PDF %s: %v\n", filename, err)
	} else {
		fmt.Printf("Successfully created %s\n", filename)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for i := 2; i <= 4; i++ {
		createLabResultsPDF(fmt.Sprintf("lab_results_%02d.pdf", i))
	}
}