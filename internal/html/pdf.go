package httphtml

import (
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/rotisserie/eris"
	"os"
	"strconv"
	"time"
)

func GeneratePDF(parsedHTML string, args ...string) (pdfData []byte, err error) {
	t := time.Now().Unix()
	// write whole the body
	dir, err := os.Getwd()
	if err != nil {
		return pdfData, eris.Wrap(err, "")
	}

	if _, err := os.Stat("./cloneTemplate/"); os.IsNotExist(err) {
		errDir := os.Mkdir("./cloneTemplate/", 0777)
		if errDir != nil {
			return pdfData, eris.Wrap(err, "")
		}
	}
	err1 := os.WriteFile("./cloneTemplate/"+strconv.FormatInt(int64(t), 10)+".html", []byte(parsedHTML), 0644)
	if err1 != nil {
		panic(err1)
	}

	f, err := os.Open("./cloneTemplate/" + strconv.FormatInt(int64(t), 10) + ".html")
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		return pdfData, eris.Wrap(err, "")
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return pdfData, eris.Wrap(err, "")
	}

	// Use arguments to customize PDF generation process
	for _, arg := range args {
		switch arg {
		case "low-quality":
			pdfg.LowQuality.Set(true)
		case "no-pdf-compression":
			pdfg.NoPdfCompression.Set(true)
		case "grayscale":
			pdfg.Grayscale.Set(true)
			// Add other arguments as needed
		}
	}

	//pdfg.AddPage(wkhtmltopdf.NewPageReader(f))
	page := wkhtmltopdf.NewPageReader(f)
	page.EnableLocalFileAccess.Set(true)

	pdfg.AddPage(page)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)
	err = pdfg.Create()
	if err != nil {
		return pdfData, eris.Wrap(err, "")
	}

	err = pdfg.WriteFile("./play.pdf")
	if err != nil {
		return pdfData, eris.Wrap(err, "")
	}

	pdfData, err = os.ReadFile("./play.pdf")
	if err != nil {
		return pdfData, eris.Wrap(err, "")
	}
	defer os.Remove("./play.pdf")

	defer os.RemoveAll(dir + "/cloneTemplate")

	return
}
