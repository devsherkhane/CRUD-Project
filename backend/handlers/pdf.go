package controllers

import (
	"crud/config"
	"fmt"
	"net/http"
	"time"

	"github.com/signintech/gopdf"
)

func DownloadStudentsPDF(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT * FROM students")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4Landscape})

	// Try to load from the Windows system folder first
	err = pdf.AddTTFFont("arial", "C:\\Windows\\Fonts\\arial.ttf")
	if err != nil {
		// Fallback to your local assets folder if the system font isn't found
		err = pdf.AddTTFFont("arial", "assets/arial.ttf")
		if err != nil {
			http.Error(w, "Font file missing: please put arial.ttf in assets/ folder", http.StatusInternalServerError)
			return
		}
	}

	// Adjusted column widths for better spacing on A4 Landscape
	headers := []string{"Sr No", "Name", "Email", "Mobile", "Address", "Gender", "DOB", "Blood"}
	colWidths := []float64{40, 100, 130, 85, 140, 60, 80, 50}

	drawHeader := func() {
		// FIXED LOGO: Using a pointer to gopdf.Rect for dimensions
		pdf.Image("assets/kk-wagh-logo.png", 30, 20, &gopdf.Rect{W: 100, H: 50})

		pdf.SetFont("arial", "", 16)
		pdf.SetXY(250, 35)
		pdf.Text("K.K Wagh Institute Of Engineering and Research")

		pdf.SetFont("arial", "", 14)
		pdf.SetXY(30, 95)
		pdf.Text("Student Records")

		// Draw Header Row with Borders
		pdf.SetFont("arial", "", 10)
		currX := 30.0
		for i, h := range headers {
			// Draw the box
			pdf.RectFromUpperLeftWithStyle(currX, 115, colWidths[i], 25, "D")
			// Add padding (X+5, Y+15) so text doesn't touch the borders
			pdf.SetXY(currX+5, 130)
			pdf.Text(h)
			currX += colWidths[i]
		}
	}

	drawFooter := func(pageNo int) {
		pdf.SetFont("arial", "", 8)
		pdf.SetXY(30, 570)
		pdf.Text("Created by Dev | " + time.Now().Format("02 Jan 2006"))
		pdf.SetXY(780, 570)
		pdf.Text(fmt.Sprintf("Page %d", pageNo))
	}

	pdf.AddPage()
	currentPage := 1
	drawHeader()
	drawFooter(currentPage)

	// ---------------- TABLE BODY ----------------
	srNo := 1
	y := 140.0 // Starting Y coordinate for data rows

	for rows.Next() {
		var (
			id                                             int
			name, address, state, district, taluka, gender string
			dob, photo, email, mobile, blood               string
			handicapped                                    bool
		)
		rows.Scan(&id, &name, &address, &state, &district, &taluka, &gender, &dob, &photo, &handicapped, &email, &mobile, &blood)

		// Page break logic for A4 Landscape
		if y > 530 {
			pdf.AddPage()
			currentPage++
			drawHeader()
			drawFooter(currentPage)
			y = 140.0
		}

		data := []string{fmt.Sprintf("%d", srNo), name, email, mobile, address, gender, dob, blood}
		currX := 30.0
		pdf.SetFont("arial", "", 9)

		for i, txt := range data {
			// Draw cell border
			pdf.RectFromUpperLeftWithStyle(currX, y, colWidths[i], 20, "D")

			// VITAL FIX: Move text down (y+13) so it sits inside the box, not on the line
			pdf.SetXY(currX+5, y+13)

			// Simple character clipping to prevent overflow
			limit := int(colWidths[i] / 5.5)
			if len(txt) > limit && limit > 3 {
				txt = txt[:limit-3] + "..."
			}

			pdf.Text(txt)
			currX += colWidths[i]
		}
		y += 20 // Move to the next row
		srNo++
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=students.pdf")
	pdf.Write(w)
}
