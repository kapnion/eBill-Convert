package main

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/jung-kurt/gofpdf"
)

// CSV Structure definition for mapping
type csvMapping struct {
	XMLPath     string
	GermanPath  string
	GermanLabel string
}

var (
	csvData   []csvMapping
	csvLoaded bool
	csvMutex  sync.RWMutex // Mutex for safe access to csvData and csvLoaded
)

func main() {
	r := gin.Default()

	// Define Swagger document
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
														Schema: &spec.Schema{
															SchemaProps: spec.SchemaProps{
																Type:        []string{"string"},
																Format:      "binary",
																Description: "The transformed PDF content",
															},
														},
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

	// Create handler for swagger.json
	r.GET("/swagger.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, swaggerSpec)
	})

	// Create handler for swagger ui
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

	pdfData, err := transformXMLToPDF(xmlData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("PDF transformation failed: %v", err)})
		return
	}

	c.Data(http.StatusOK, "application/pdf", pdfData)
}

func transformXMLToPDF(xmlData []byte) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	// Load CSV if not already loaded
	loadCSV()

	xmlReader := bytes.NewReader(xmlData)
	decoder := xml.NewDecoder(xmlReader)
	var currentPath string
	var groupStack []string
	elementCounts := make(map[string]int)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error decoding XML: %w", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			name := t.Name.Local
			if currentPath == "" {
				currentPath = name
			} else {
				currentPath += "->" + name
			}
			groupStack = append(groupStack, name)

		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" {
				printElement(pdf, currentPath, text, groupStack, elementCounts)
			}
		case xml.EndElement:
			if len(groupStack) > 0 {
				groupStack = groupStack[:len(groupStack)-1]
			}

			parts := strings.Split(currentPath, "->")
			if len(parts) > 0 {
				currentPath = strings.Join(parts[:len(parts)-1], "->")
			} else {
				currentPath = ""
			}

		}

	}

	var buffer bytes.Buffer
	err := pdf.Output(&buffer)

	if err != nil {
		return nil, fmt.Errorf("error creating pdf: %w", err)
	}

	return buffer.Bytes(), nil
}

func printElement(pdf *gofpdf.Fpdf, currentPath, text string, groupStack []string, elementCounts map[string]int) {

	germanLabel := lookupLabel(currentPath)
	parts := strings.Split(currentPath, "->")
	elementName := parts[len(parts)-1]
	elementName = strings.TrimSpace(elementName)
	header := ""

	if len(groupStack) > 1 {
		header = lookupHeader(currentPath)

		if header != "" {
			// Print header if changed

			if len(groupStack) > 0 {
				lastHeader := groupStack[len(groupStack)-2]
				lastHeader = strings.TrimSpace(lastHeader)
				if header != lastHeader {
					pdf.Ln(5)
					pdf.SetFont("Arial", "B", 12)
					pdf.Cell(0, 10, header)
					pdf.Ln(5)
					pdf.SetFont("Arial", "", 12)
					//groupStack = groupStack[:len(groupStack)-1] // removing last group as its printed now
				}

			}
		}

	}

	label := elementName

	if germanLabel != "" {
		label = germanLabel
	}

	elementCountKey := currentPath
	elementCounts[elementCountKey]++
	if elementCounts[elementCountKey] > 1 {
		label = fmt.Sprintf("%s %d", label, elementCounts[elementCountKey])

	}

	pdf.Cell(0, 10, fmt.Sprintf("%s: %s", label, text))
	pdf.Ln(5)
}

func loadCSV() {
	csvMutex.RLock()
	if csvLoaded {
		csvMutex.RUnlock()
		return
	}
	csvMutex.RUnlock()
	csvMutex.Lock()
	defer csvMutex.Unlock()

	if csvLoaded {
		return // Double check in case another thread loaded it while waiting for the lock
	}

	file, err := os.Open("translations.csv")
	if err != nil {
		log.Printf("Error opening CSV file: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	_, err = reader.Read() // Skip header row
	if err != nil {
		log.Printf("Error reading header row: %v", err)
		return
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV row: %v", err)
			continue
		}
		if len(row) == 3 {
			csvData = append(csvData, csvMapping{
				XMLPath:     row[0],
				GermanPath:  row[1],
				GermanLabel: row[2],
			})
		} else {
			log.Printf("Skipping invalid CSV row: %v", row)
		}
	}
	csvLoaded = true
}

func lookupLabel(xmlPath string) string {
	csvMutex.RLock()
	defer csvMutex.RUnlock()

	for _, mapping := range csvData {
		if strings.EqualFold(mapping.XMLPath, xmlPath) {
			return mapping.GermanLabel
		}
	}

	return ""
}

func lookupHeader(xmlPath string) string {
	csvMutex.RLock()
	defer csvMutex.RUnlock()
	for _, mapping := range csvData {
		if strings.EqualFold(mapping.XMLPath, xmlPath) {
			parts := strings.Split(mapping.GermanPath, "->")
			if len(parts) > 1 {
				return parts[len(parts)-2]
			}
		}
	}
	return ""
}
