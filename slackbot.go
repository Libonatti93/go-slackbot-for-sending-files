package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// Substitua pelo seu token de bot do Slack
	slackBotToken := "xoxb-your-slack-bot-token"
	channelID := "C12345678" // Substitua pelo ID do canal ou usuário que você deseja enviar o arquivo

	// Nome do arquivo que deseja enviar
	filePath := "path/to/your/file.txt"

	// Chame a função para enviar o arquivo
	err := uploadFileToSlack(slackBotToken, channelID, filePath)
	if err != nil {
		log.Fatalf("Erro ao enviar arquivo: %v", err)
	}

	fmt.Println("Arquivo enviado com sucesso!")
}

func uploadFileToSlack(token, channelID, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("não foi possível abrir o arquivo: %v", err)
	}
	defer file.Close()

	// Leia o conteúdo do arquivo
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("não foi possível ler o arquivo: %v", err)
	}

	// Prepare o formulário de dados para upload
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", filePath)
	if err != nil {
		return fmt.Errorf("não foi possível criar o formulário: %v", err)
	}

	_, err = fw.Write(fileContents)
	if err != nil {
		return fmt.Errorf("não foi possível escrever o arquivo no formulário: %v", err)
	}

	// Adicione outros campos ao formulário
	w.WriteField("channels", channelID)
	w.Close()

	// Faça a requisição HTTP para o Slack
	req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", &b)
	if err != nil {
		return fmt.Errorf("não foi possível criar a requisição HTTP: %v", err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar a requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("falha ao enviar arquivo, status: %v, resposta: %v", resp.Status, string(body))
	}

	return nil
}
