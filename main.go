package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/jung-kurt/gofpdf"
)

var translations map[string]string

func main() {
	r := gin.Default()

	loadTranslations("translations.csv")

	swaggerSpec := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger: "2.0",
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       "XML Processing API",
					Description: "API for processing XML files.",
					Version:     "1.0.0",
				},
			},
			Paths: &spec.Paths{
				Paths: map[string]spec.PathItem{
					"/xmltohtml": {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{
								OperationProps: spec.OperationProps{
									Description: "Transforms XML to HTML.",
									Consumes:    []string{"application/xml"},
									Produces:    []string{"text/html"},
									Parameters: []spec.Parameter{
										{
											ParamProps: spec.ParamProps{
												Name:        "xmlFile",
												In:          "formData",
												Description: "The XML file.",
												Required:    true,
											},
											SimpleSchema: spec.SimpleSchema{
												Type: "file",
											},
										},
									},
									Responses: &spec.Responses{
										ResponsesProps: spec.ResponsesProps{
											StatusCodeResponses: map[int]spec.Response{
												200: {
													ResponseProps: spec.ResponseProps{
														Description: "HTML content generated from XML.",
													},
												},
												400: {
													ResponseProps: spec.ResponseProps{
														Description: "Invalid XML or other errors.",
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"/xmltopdf": {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{
								OperationProps: spec.OperationProps{
									Description: "Transforms XML to PDF.",
									Consumes:    []string{"multipart/form-data"},
									Produces:    []string{"application/pdf"},
									Parameters: []spec.Parameter{
										{
											ParamProps: spec.ParamProps{
												Name:        "xmlFile",
												In:          "formData",
												Description: "The XML file.",
												Required:    true,
											},
											SimpleSchema: spec.SimpleSchema{
												Type: "file",
											},
										},
									},
									Responses: &spec.Responses{
										ResponsesProps: spec.ResponsesProps{
											StatusCodeResponses: map[int]spec.Response{
												200: {
													ResponseProps: spec.ResponseProps{
														Description: "Successfully transformed the XML file to PDF",
													},
												},
												400: {
													ResponseProps: spec.ResponseProps{
														Description: "Invalid XML or other errors.",
														Schema: &spec.Schema{
															SchemaProps: spec.SchemaProps{
																Type: spec.StringOrArray{"object"},
																Properties: map[string]spec.Schema{
																	"error": {
																		SchemaProps: spec.SchemaProps{
																			Type: spec.StringOrArray{"string"},
																		},
																	},
																	"message": {
																		SchemaProps: spec.SchemaProps{
																			Type: spec.StringOrArray{"string"},
																		},
																	},
																},
															},
														},
													},
												},
												500: {
													ResponseProps: spec.ResponseProps{
														Description: "Internal server error",
														Schema: &spec.Schema{
															SchemaProps: spec.SchemaProps{
																Type: spec.StringOrArray{"object"},
																Properties: map[string]spec.Schema{
																	"error": {
																		SchemaProps: spec.SchemaProps{
																			Type: spec.StringOrArray{"string"},
																		},
																	},
																	"message": {
																		SchemaProps: spec.SchemaProps{
																			Type: spec.StringOrArray{"string"},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	r.GET("/swagger.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, swaggerSpec)
	})

	opts := middleware.SwaggerUIOpts{
		SpecURL: "/swagger.json",
	}
	sh := middleware.SwaggerUI(opts, nil)
	r.GET("/docs", func(c *gin.Context) {
		sh.ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/xmltohtml", handleXMLtoHTML)
	r.POST("/xmltopdf", handleXMLtoPDF)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleXMLtoHTML(c *gin.Context) {
	file, err := c.FormFile("xmlFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file missing"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file open error"})
		return
	}
	defer src.Close()

	htmlData := make([]byte, file.Size)
	_, err = src.Read(htmlData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file read error"})
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", htmlData)
}

func handleXMLtoPDF(c *gin.Context) {
	file, err := c.FormFile("xmlFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file missing"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file open error"})
		return
	}
	defer src.Close()

	xmlData := make([]byte, file.Size)
	_, err = src.Read(xmlData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file read error"})
		return
	}

	var doc interface{}
	err = xml.Unmarshal(xmlData, &doc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse XML"})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	processXMLLineByLine(xmlData, pdf)

	pdfFile := filepath.Join(os.TempDir(), "output.pdf")
	err = pdf.OutputFileAndClose(pdfFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create PDF: %v", err)})
		return
	}

	pdfData, err := os.ReadFile(pdfFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read PDF: %v", err)})
		return
	}

	os.Remove(pdfFile)

	c.Data(http.StatusOK, "application/pdf", pdfData)
}

func processXMLLineByLine(xmlData []byte, pdf *gofpdf.Fpdf) {
	decoder := xml.NewDecoder(strings.NewReader(string(xmlData)))
	var currentElement string
	var elementStack []string

	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		switch token := tok.(type) {
		case xml.StartElement:
			currentElement = token.Name.Local
			elementStack = append(elementStack, currentElement)
		case xml.EndElement:
			elementStack = elementStack[:len(elementStack)-1]
		case xml.CharData:
			content := strings.TrimSpace(string(token))
			if content != "" {
				path := strings.Join(elementStack, "->")
				label := getLabel(path)
				pdf.Cell(40, 10, fmt.Sprintf("%s: %s", label, content))
			}
		}
	}
}

func getLabel(path string) string {
	if label, ok := translations[path]; ok {
		return label
	}
	parts := strings.Split(path, "->")
	return parts[len(parts)-1]
}

func loadTranslations(filePath string) {
	translations = make(map[string]string)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open translations file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read translations file: %v", err)
	}

	for _, record := range records[1:] {
		translations[record[0]] = record[3]
	}
}
