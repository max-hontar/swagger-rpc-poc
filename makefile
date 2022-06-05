.DEFAULT_GOAL := swagger
swagger:
	swag init -d ./ -g ./main.go

swagger-fmt:
	swag fmt -d ./ -g ./main.go