package main

import (
	"github.com/jinzhu/gorm"
	"github.com/rafikmoreira/codepix/application/grpc"
	"github.com/rafikmoreira/codepix/infrastructure/db"
	"os"
)

var database *gorm.DB

func main() {
	database = db.ConnectDB(os.Getenv("env"))

	grpc.StartGrpcServer(database, 50051)
}
