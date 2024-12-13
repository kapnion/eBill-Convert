package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/jung-kurt/gofpdf"
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
									Consumes:    []string{"application/xml"},
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
														Description: "PDF content generated from XML.",
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

	// Parse XML
	var doc interface{}
	err = xml.Unmarshal(xmlData, &doc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse XML"})
		return
	}

	// Create PDF using gofpdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Converted PDF content")

	// Save PDF to a file
	pdfFile := filepath.Join(os.TempDir(), "output.pdf")
	err = pdf.OutputFileAndClose(pdfFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create PDF: %v", err)})
		return
	}

	// Read the generated PDF file
	pdfData, err := os.ReadFile(pdfFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read PDF: %v", err)})
		return
	}

	// Clean up the temporary file
	os.Remove(pdfFile)

	c.Data(http.StatusOK, "application/pdf", pdfData)
}
