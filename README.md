# Conversor de Fatura Itaú .xls para Organizze .xls

Solução para converter as faturas de cartão de crédito do Itaú para o modelo suportado pela Organizze para importação. Já que o Itaú não exporta a fatura em .OFX (shit).

Não existem soluções práticas como libs para gerar arquivos .xls, isso porque esse formato já é _deprecated_ e foi substituído pelo .xlxs. 

Portanto a solução usada foi:
1. Ler a fatura.xls do Itaú 
2. Converter a fatura para para o modelo da Organizze porém em formato .xlxs
3. Usar o Unoserver (Libreoffice as a service) para converter de .xlxs para .xls (formato final a ser usado na Organizze)

_Obs.: Optei por usar o Unoserver via Docker, porque o suporte oficial dele é apenas para Linux._

## Pré Requisitos
 - Golang 1.21+
 - Docker

## Como usar
1. Altere o valor da variável `invoice_path` do Makefile para referenciar sua fatura do Itaú
2. Rode o comando `make run`
3. Usar o arquivo `organizze-entries-to-import.xls` gerado na raiz do projeto para importar na Organizze

