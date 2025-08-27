package sudoku

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

func ProcessImage(key string, filename string, file []byte) (Game, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return Game{}, err
	}

	_, err = part.Write(file)
	if err != nil {
		return Game{}, err
	}

	err = writer.Close()
	if err != nil {
		return Game{}, err
	}

	req, err := http.NewRequest("POST", "https://sudoku-ocr.p.rapidapi.com/scan-puzzle", &body)
	if err != nil {
		return Game{}, err
	}
	req.Header.Set("x-rapidapi-host", "sudoku-ocr.p.rapidapi.com")
	req.Header.Set("x-rapidapi-key", key)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Game{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Game{}, err
	}

	//log.Println("key", key)
	log.Println("Response Status:", resp.Status)
	log.Println("Response Body:", string(responseBody))

	board, err := ConvertFromOCRFormat(string(responseBody))
	if err != nil {
		return Game{}, err
	}

	g := Game{}
	err = g.FillBasic(board)
	if err != nil {
		return Game{}, err
	}
	return g, nil
}
