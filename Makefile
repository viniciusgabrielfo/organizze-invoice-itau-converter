invoice_path ?= $(HOME)/Downloads/Fatura-Excel.xls

run-invoice-itau-consumer:
	go run cmd/main.go --file $(invoice_path)

up-unoconv:
	docker buildx build -t unoconv .
	docker run -d --name unoconv unoconv

down-unoconv:
	docker container stop unoconv
	docker container rm unoconv

run: run-invoice-itau-consumer up-unoconv
	sleep 5
	docker exec unoconv unoconvert fatura.xlsx fatura.xls --convert-to xls
	docker cp unoconv:/fatura.xls .
	$(MAKE) down-unoconv


