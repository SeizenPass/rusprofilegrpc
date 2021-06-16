### Решение тестового задания rusprofile  

- Makefile для запуска всех команд
```shell script
make genG - генерирует grpc файлы на основе .proto файла
make genW - генерирует grpc-gateway на основе .proto файла и search.yaml
make genS - генерирует rusprofile.swagger.json для Swagger
make build - запускает 3 предыдуших команд
make clean - очистка сгенерированных .go файлов
make docker - поднимает docker-compose
```

- localhost:8081/search/{uin} - поиск по ИИНу, который возвращает stream из SearchResponse.  
- В случае если критерию поиска совпало более 1 организации, то возвращает коллекцию  
- localhost:8080 - Swagger UI