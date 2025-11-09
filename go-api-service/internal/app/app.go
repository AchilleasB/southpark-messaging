package app

import (
	"context"
	"encoding/json"

	"github.com/AchilleasB/southpark-messaging/go-api-service/internal/domain"
	"github.com/AchilleasB/southpark-messaging/go-api-service/internal/ports"
)

type App struct {
	pub ports.Publisher
}

func NewApp(p ports.Publisher) *App {
	return &App{pub: p}
}

func (app *App) PublishMessage(ctx context.Context, msg *domain.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return app.pub.Publish(ctx, "southpark_messages", data)
}
