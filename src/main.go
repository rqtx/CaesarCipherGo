package main

import (
	"bytes"
	caesar "caesar/cipher"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var TOKEN = "d3015fefa4bee006752f264e0a28c28ce9f7b77a"

func getData() []byte {
	resp, err := http.Get(fmt.Sprintf("https://api.codenation.dev/v1/challenge/dev-ps/generate-data?token=%s", TOKEN))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func saveJson(data []byte, fileName string) {
	err := ioutil.WriteFile(fileName, data, 0755)
	if err != nil {
		panic(err)
	}
}

func loadJson(fileName string) string {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func getSha1(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func main() {
	fileName := "answer.json"

	saveJson(getData(), fileName)

	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Unable to read file: %v", err)
	}

	var raw map[string]interface{}
	json.Unmarshal(fileData, &raw)

	raw["decifrado"] = caesar.Decode(raw["cifrado"].(string), int((raw["numero_casas"].(float64))))
	raw["resumo_criptografico"] = getSha1([]byte(raw["decifrado"].(string)))

	b, _ := json.Marshal(raw)
	saveJson(b, fileName)

	remoteURL := fmt.Sprintf("https://api.codenation.dev/v1/challenge/dev-ps/submit-solution?token=%s", TOKEN)
	request, err := newfileUploadRequest(remoteURL, nil, "answer", fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
	}
}
