package buscartitulo

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

// <-chan - Canal somente leitura
func Titulo(urls ...string) <-chan string {
	c := make(chan string)
	for _, url := range urls {
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				c <- fmt.Sprintf("Erro ao acessar %s: %v", url, err)
				return
			}
			defer resp.Body.Close()
			html, err := io.ReadAll(resp.Body)
			if err != nil {
				c <- fmt.Sprintf("Erro ao ler resposta de %s: %v", url, err)
				return
			}
			r, err := regexp.Compile("<title>(.*?)<\\/title>")
			if err != nil {
				c <- fmt.Sprintf("Erro na regex: %v", err)
				return
			}
			aRetorno := r.FindStringSubmatch(string(html))
			if len(aRetorno) < 2 {
				c <- fmt.Sprintf("Título não encontrado em %s", url)
				return
			}
			c <- aRetorno[1]
		}(url)
	}
	return c
}
