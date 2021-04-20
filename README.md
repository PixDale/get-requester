## Instruções de inicialização

Ao executar a aplicação, passar o caminho do arquivo de configuração como parâmetro

```sh
./requestTester.exe -config "configTeste.json"
```

 Utilizar arquivo config.json como exemplo para criar outros configs

### Parâmetros:

- __``url``__: URL com a endpoint da API a ser requisitada.
- __``interval``__: Intervalo (ms) entre requisições (Obs: O intervalo mínimo funcional será o tempo de resposta da requisição).
- __``threads``__: Número de goroutines que serão criadas para realizar a sequência de testes
- __``format``__: Formato do arquivo que será gerado com a resposta da requisição.
- __``folderName``__: Nome do diretório que será criado para salvar as resposta das requisições
- __``qtd``__: Quantidade de requisições que serão feitas por '``thread``'.
- __``rampaInicio``__: Tempo (ms) máximo para que todas as '``threads``' sejam iniciadas.
- __``forever``__: Define se o teste ficará em loop infinito para cada '``thread``' iniciada, ignorando o parâmetro 'qtd'.
- __``persist``__: Define se a resposta da requisição deve ser salva em disco.