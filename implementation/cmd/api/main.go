package main

import (
 "context"
 "log"
 "log/slog"
 "net/http"
 "os"
 "os/signal"
 "syscall"
 "time"

 "github.com/Servitor108/notification-delivery-design/implementation/internal/httpapi"
 "github.com/Servitor108/notification-delivery-design/implementation/internal/notify"
)

type logProvider struct{}
func(logProvider)Send(ctx context.Context,r notify.ProviderRequest)(notify.ProviderResult,error) {
 slog.Info("demo provider accepted notification","notification_id",r.NotificationID,"channel",r.Channel)
 return notify.ProviderResult{ProviderMessageID:"demo-"+r.NotificationID},nil
}
func main() {
 repo:=notify.NewMemoryRepository(); service:=notify.NewService(repo,logProvider{},5,5*time.Second)
 ctx,stop:=signal.NotifyContext(context.Background(),syscall.SIGINT,syscall.SIGTERM); defer stop()
 go service.Run(ctx,4)
 mux:=http.NewServeMux(); httpapi.New(service).Register(mux)
 server:=&http.Server{Addr:env("HTTP_ADDR",":8080"),Handler:mux,ReadHeaderTimeout:5*time.Second}
 go func(){if err:=server.ListenAndServe();err!=nil&&err!=http.ErrServerClosed{log.Printf("server: %v",err);stop()}}()
 <-ctx.Done(); shutdown,cancel:=context.WithTimeout(context.Background(),5*time.Second); defer cancel(); _=server.Shutdown(shutdown)
}
func env(k,d string)string{if v:=os.Getenv(k);v!=""{return v};return d}
