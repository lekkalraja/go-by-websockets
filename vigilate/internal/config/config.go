package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/lekkalraja/go-by-websockets/vigilate/internal/channeldata"
	"github.com/lekkalraja/go-by-websockets/vigilate/internal/driver"
	"github.com/pusher/pusher-http-go"
	"github.com/robfig/cron/v3"
	"html/template"
)

// AppConfig holds application configuration
type AppConfig struct {
	DB            *driver.DB
	Session       *scs.SessionManager
	InProduction  bool
	Domain        string
	MonitorMap    map[int]cron.EntryID
	PreferenceMap map[string]string
	Scheduler     *cron.Cron
	WsClient      pusher.Client
	PusherSecret  string
	TemplateCache map[string]*template.Template
	MailQueue     chan channeldata.MailJob
	Version       string
	Identifier    string
}
