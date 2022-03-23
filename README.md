# wb-orders-test0
Test task for an internship WB


### Build:

```
make build
```

### Если приложение запускается впервые, требуется также применить команду для миграции БД:

```
make migrate-up
```

### Run:

```
make run
```

### В отдельном терминале запускается Subscriber: 
```
go run cmd/sub.go
```

### В отдельном терминале запускается Publisher: 
```
go run cmd/sub.go
```
