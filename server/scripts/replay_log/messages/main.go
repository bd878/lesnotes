package main

import (
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/bd878/gallery/server/messages/pkg/loadbalance"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/loadbalance"
)

func main() {
	PGConn, ok := os.LookupEnv("PG_CONN")
	if !ok {
		panic("PG_CONN env required")
	}
	addr, ok := os.LookupEnv("ADDR")
	if !ok {
		panic("ADDR env required")
	}

	fmt.Fprintln(os.Stdout, "=== running messages migration ===")
	fmt.Fprintln(os.Stdout, "PGConn:", PGConn, "addr:", addr)

	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			addr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "grpc conn ok")

	messagesClient := api.NewMessagesClient(conn)
	fmt.Fprintln(os.Stdout, "messages client ok")

	translationsClient := api.NewTranslationsClient(conn)
	fmt.Fprintln(os.Stdout, "translations client ok")

	pool, err := pgxpool.New(context.Background(), PGConn)
	defer pool.Close()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "pg pool ok")

	fmt.Fprintln(os.Stdout, "scanning messages")

	messagesQuery := "SELECT id, text, private, name, user_id, title FROM messages.messages"

	rows, err := pool.Query(context.Background(), messagesQuery)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "messages rows ok")

	messages := make([]*model.Message, 0)
	for rows.Next() {
		message := &model.Message{}

		err := rows.Scan(&message.ID, &message.Text, &message.Private, &message.Name, &message.UserID, &message.Title)
		if err != nil {
			panic(err)
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "scanned messages rows:", len(messages))

	fmt.Fprintln(os.Stdout, "scanning files")

	for _, message := range messages {

		filesQuery := "SELECT file_id FROM messages.files WHERE message_id = $1"
		fileRows, err := pool.Query(context.Background(), filesQuery, message.ID)
		defer fileRows.Close()
		if err != nil {
			panic(err)
		}

		fileIDs := make([]int64, 0)
		for fileRows.Next() {
			var fileID int64
			err = fileRows.Scan(&fileID)
			if err != nil {
				panic(err)
			}

			fileIDs = append(fileIDs, fileID)
		}

		if err = rows.Err(); err != nil {
			fmt.Fprintln(os.Stdout, "[WARN] failed to scan files:", err)
		}

		message.FileIDs = fileIDs

	}

	fmt.Fprintln(os.Stdout, "scanning files ok")
	fmt.Fprintln(os.Stdout, "saving messages")

	for _, message := range messages {

		fmt.Fprintln(os.Stdout, "save message", "id", message.ID)

		_, err = messagesClient.SaveMessage(context.Background(), &api.SaveMessageRequest{
			Id:        message.ID,
			UserId:    message.UserID,
			Text:      message.Text,
			Private:   message.Private,
			Name:      message.Name,
			Title:     message.Title,
			FileIds:   message.FileIDs,
		})
		if err != nil {
			fmt.Fprintln(os.Stdout, "[WARN] failed to save message:", err)
		}

	}

	fmt.Fprintln(os.Stdout, "messages saved")
	fmt.Fprintln(os.Stdout, "scanning translations")

	translations := make([]*api.SaveTranslationRequest, 0)
	for _, message := range messages {

		translationsQuery := "SELECT lang, text, title FROM messages.translations WHERE message_id = $1"
		translationRows, err := pool.Query(context.Background(), translationsQuery, message.ID)
		defer translationRows.Close()
		if err != nil {
			panic(err)
		}

		for translationRows.Next() {
			translation := &api.SaveTranslationRequest{
				UserId:      message.UserID,
				MessageId:   message.ID,
			}

			err = translationRows.Scan(&translation.Lang, &translation.Text, &translation.Title)
			if err != nil {
				panic(err)
			}

			translations = append(translations, translation)
		}

		if err = translationRows.Err(); err != nil {
			fmt.Fprintln(os.Stdout, "[WARN] failed to scan translations:", err)
		}

	}

	fmt.Fprintln(os.Stdout, "scanning translations ok")
	fmt.Fprintln(os.Stdout, "saving translations")

	for _, translation := range translations {

		fmt.Fprintln(os.Stdout, "save translation", "message_id", translation.MessageId, "lang", translation.Lang)

		_, err = translationsClient.SaveTranslation(context.Background(), translation)
		if err != nil {
			fmt.Fprintln(os.Stdout, "[WARN] failed to save translation:", err)
		}
	}

	fmt.Fprintln(os.Stdout, "translations saved")

	fmt.Fprintln(os.Stdout, "=== messages migration done ===")
}