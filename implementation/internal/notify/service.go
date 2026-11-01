package notify

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "errors"
 "math/rand"
 "net/http"
 "strings"
 "sync"
 "time"
)

var (
 ErrNotFound = errors.New("notification not found")
 ErrDuplicate = errors.New("idempotency key already used with different request")
)

type Repository interface {
 Create(context.Context, Notification) (Notification,bool,error)
 Get(context.Context,string) (Notification,error)
 Claim(context.Context) (Notification,error)
 Complete(context.Context,string,string) error
 Reschedule(context.Context,string,time.Time,string) error
}

type Service struct { repo Repository; provider Provider; maxAttempts int; clientTimeout time.Duration }
func NewService(repo Repository,p Provider,maxAttempts int,timeout time.Duration)*Service {
 if maxAttempts<1 { maxAttempts=1 }; if timeout<=0 { timeout=10*time.Second }
 return &Service{repo:repo,provider:p,maxAttempts:maxAttempts,clientTimeout:timeout}
}
func (s *Service) Enqueue(ctx context.Context,n Notification)(Notification,bool,error) {
 if strings.TrimSpace(n.IdempotencyKey)=="" { return Notification{},false,errors.New("idempotency key is required") }
 if strings.TrimSpace(n.Recipient)==""||strings.TrimSpace(n.Channel)==""||strings.TrimSpace(n.Template)=="" { return Notification{},false,errors.New("recipient, channel and template are required") }
 if n.ID=="" { sum:=sha256.Sum256([]byte(n.IdempotencyKey)); n.ID=hex.EncodeToString(sum[:16]) }
 n.Status=StatusQueued; n.CreatedAt=time.Now().UTC(); n.NextAttemptAt=n.CreatedAt
 return s.repo.Create(ctx,n)
}
func (s *Service) Get(ctx context.Context,id string)(Notification,error) { return s.repo.Get(ctx,id) }
func (s *Service) ProcessOne(ctx context.Context)error {
 n,err:=s.repo.Claim(ctx); if err!=nil { return err }
 callCtx,cancel:=context.WithTimeout(ctx,s.clientTimeout); defer cancel()
 _,err=s.provider.Send(callCtx,ProviderRequest{NotificationID:n.ID,Recipient:n.Recipient,Channel:n.Channel,Template:n.Template,Payload:n.Payload})
 if err==nil { return s.repo.Complete(ctx,n.ID,"") }
 if !retryable(err)||n.Attempts>=s.maxAttempts { return s.repo.Complete(ctx,n.ID,err.Error()) }
 return s.repo.Reschedule(ctx,n.ID,Backoff(n.Attempts),err.Error())
}
func retryable(err error)bool {
 var ne interface{Timeout()bool}; if errors.As(err,&ne)&&ne.Timeout(){return true}
 var se interface{StatusCode()int}; if errors.As(err,&se){code:=se.StatusCode();return code==http.StatusRequestTimeout||code==http.StatusTooManyRequests||code>=500}
 return errors.Is(err,context.DeadlineExceeded)
}
func Backoff(attempt int)time.Duration {
 if attempt<1 {attempt=1}; base:=time.Second*time.Duration(1<<min(attempt-1,10))
 return time.Duration(rand.Int63n(int64(base)+1))
}
func min(a,b int)int {if a<b{return a};return b}
func (s *Service) Run(ctx context.Context,workers int) {
 if workers<1 {workers=1}; var wg sync.WaitGroup
 for i:=0;i<workers;i++ {wg.Add(1);go func(){defer wg.Done();ticker:=time.NewTicker(100*time.Millisecond);defer ticker.Stop();for{select{case<-ctx.Done():return;case<-ticker.C:if err:=s.ProcessOne(ctx);err!=nil&&!errors.Is(err,ErrNotFound){/* repository/provider errors are surfaced by metrics in a production adapter */}}}}()}
 <-ctx.Done(); wg.Wait()
}
