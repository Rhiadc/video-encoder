package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/video-encoder/application/repositories"
	"github.com/video-encoder/application/services"
	"github.com/video-encoder/config"
	"github.com/video-encoder/framework/database"
	"github.com/video-encoder/pkg/s3"
	"log"
	"os"
	"time"
)

func init() {

	//Load locahost envs using .env file
	config.LoadEnvs()

}

func main() {
	fmt.Println("teste")
	dbConfig := database.Config{
		DSN: os.Getenv("DSN"),
	}
	db := database.InitGorm(&dbConfig)

	ses, err := s3.NewSession(s3.Config{
		Address: "http://localhost:4566",
		Region:  "us-east-1",
		Profile: "localstack",
		ID:      "test",
		Secret:  "test",
	})
	if err != nil {
		log.Fatalln(err)
	}

	processVideo := &cobra.Command{
		Use:   "KafkaConsumerTransactionDLQ",
		Short: "Producer Transaction Events from Transaction Events DLQ",
		Run: func(cli *cobra.Command, args []string) {
			ctx := context.Background()
			videoRepository := repositories.NewVideoRepository(db)
			s3 := s3.NewS3(ses, time.Second*5)
			videoService := services.NewVideoService(videoRepository, s3)
			if err := videoService.Upload(ctx, "mybucket"); err != nil {
				fmt.Println(err)
			}
		},
	}

	var rootCmd = &cobra.Command{Use: "APP"}
	rootCmd.AddCommand(processVideo)
	rootCmd.Execute()
}
