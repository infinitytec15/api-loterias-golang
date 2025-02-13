package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jung-kurt/gofpdf"
	_ "modernc.org/sqlite"
)

type LoteriaResult struct {
	Acumulado           bool     `json:"acumulado"`
	DataApuracao        string   `json:"dataApuracao"`
	DataProximoConcurso string   `json:"dataProximoConcurso"`
	DezenasSorteadas    []string `json:"dezenasSorteadasOrdemSorteio"`
	TipoJogo            string   `json:"tipoJogo"`
	Numero              int      `json:"numero"`
}

var (
	db             *sql.DB
	ultimoConcurso int
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}
	db, err = sql.Open("sqlite", "./bancoloteria.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS loteria (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		concurso INTEGER,
		data_apuracao TEXT,
		dezenas TEXT,
		horario TEXT,
		jogo TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	ultimoConcurso = 0
}

func main() {
	for {

		result, err := consultarLoteria("federal")
		if err != nil {
			log.Println("Erro ao consultar loteria:", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		err = salvarNoBanco(result)
		if err != nil {
			log.Println("Erro ao salvar no banco:", err)
		}

		err = enviarEmbedDiscord(result)
		if err != nil {
			log.Println("Erro ao enviar embed para o Discord:", err)
		}

		if result.Numero != ultimoConcurso {
			pdfFile := gerarPDF(result)
			log.Println("PDF gerado:", pdfFile)
			err = enviarPDFDiscord(pdfFile)
			if err != nil {
				log.Println("Erro ao enviar PDF para o Discord:", err)
			}

			ultimoConcurso = result.Numero
		}

		time.Sleep(60 * time.Second)
	}
}

func consultarLoteria(concurso string) (LoteriaResult, error) {
	url := fmt.Sprintf("https://servicebus2.caixa.gov.br/portaldeloterias/api/%s", concurso)
	resp, err := http.Get(url)
	if err != nil {
		return LoteriaResult{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return LoteriaResult{}, err
	}

	var result LoteriaResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return LoteriaResult{}, err
	}

	return result, nil
}

func salvarNoBanco(result LoteriaResult) error {
	query := `
	INSERT INTO loteria (concurso, data_apuracao, dezenas, horario, jogo)
	VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, result.Numero, result.DataApuracao, fmt.Sprint(result.DezenasSorteadas), time.Now().Format("2006-01-02 15:04:05"), result.TipoJogo)
	return err
}

func gerarPDF(result LoteriaResult) string {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "Resultado da Loteria")
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Concurso: %d", result.Numero))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Data de Apuração: %s", result.DataApuracao))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Dezenas Sorteadas: %v", result.DezenasSorteadas))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Tipo de Jogo: %s", result.TipoJogo))

	pdfFile := "resultado_loteria.pdf"
	err := pdf.OutputFileAndClose(pdfFile)
	if err != nil {
		log.Fatal(err)
	}

	return pdfFile
}

func enviarEmbedDiscord(result LoteriaResult) error {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL_EMBED")
	if webhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL_EMBED não configurada")
	}

	embed := map[string]interface{}{
		"title":       "Resultado da Loteria",
		"description": "Confira os resultados mais recentes da loteria.",
		"color":       0x00ff00,
		"fields": []map[string]interface{}{
			{"name": "Concurso", "value": fmt.Sprintf("%d", result.Numero), "inline": true},
			{"name": "Data de Apuração", "value": result.DataApuracao, "inline": true},
			{"name": "Dezenas Sorteadas", "value": fmt.Sprintf("%v", result.DezenasSorteadas), "inline": false},
			{"name": "Tipo de Jogo", "value": result.TipoJogo, "inline": true},
		},
	}

	payload := map[string]interface{}{
		"embeds": []interface{}{embed},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao converter payload para JSON: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("erro ao enviar embed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("erro ao enviar embed: HTTP %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

func enviarPDFDiscord(pdfFile string) error {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL_PDF")
	if webhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL_PDF não configurada")
	}

	file, err := os.Open(pdfFile)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo PDF: %v", err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", pdfFile)
	if err != nil {
		return fmt.Errorf("erro ao criar form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("erro ao copiar arquivo para o payload: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", webhookURL, body)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar arquivo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("erro ao enviar arquivo: HTTP %d - %s", resp.StatusCode, string(body))
	}

	return nil
}
