package notify

import (
 "context"
 "errors"
 "testing"
 "time"
)
type fakeProvider struct{err error}
func(f fakeProvider)Send(ctx context.Context,r ProviderRequest)(ProviderResult,error){return ProviderResult{},f.err}
func TestEnqueueIsIdempotent(t *testing.T) {
 repo:=NewMemoryRepository(); s:=NewService(repo,fakeProvider{},3,time.Second); n:=Notification{IdempotencyKey:"evt-1",Recipient:"user-1",Channel:"email",Template:"receipt"}
 first,created,err:=s.Enqueue(context.Background(),n); if err!=nil||!created{t.Fatalf("first enqueue: created=%v err=%v",created,err)}
 second,created,err:=s.Enqueue(context.Background(),n); if err!=nil||created{t.Fatalf("duplicate enqueue: created=%v err=%v",created,err)}
 if first.ID!=second.ID{t.Fatal("idempotent request returned a different ID")}
}
func TestSameKeyDifferentRequestConflicts(t *testing.T) {
 repo:=NewMemoryRepository(); s:=NewService(repo,fakeProvider{},3,time.Second)
 _,_,_ = s.Enqueue(context.Background(),Notification{IdempotencyKey:"same",Recipient:"u1",Channel:"sms",Template:"otp"})
 _,_,err:=s.Enqueue(context.Background(),Notification{IdempotencyKey:"same",Recipient:"u2",Channel:"sms",Template:"otp"})
 if !errors.Is(err,ErrDuplicate){t.Fatalf("expected conflict, got %v",err)}
}
func TestProcessSuccessfulDelivery(t *testing.T) {
 repo:=NewMemoryRepository(); s:=NewService(repo,fakeProvider{},3,time.Second)
 n,_,err:=s.Enqueue(context.Background(),Notification{IdempotencyKey:"success",Recipient:"u",Channel:"email",Template:"welcome"}); if err!=nil{t.Fatal(err)}
 if err:=s.ProcessOne(context.Background());err!=nil{t.Fatal(err)}
 got,err:=repo.Get(context.Background(),n.ID); if err!=nil{t.Fatal(err)}; if got.Status!=StatusAccepted{t.Fatalf("status=%s",got.Status)}
}
func TestBackoffIsBounded(t *testing.T) {for i:=1;i<20;i++{d:=Backoff(i);if d<0||d>time.Duration(1<<min(i-1,10))*time.Second{t.Fatalf("invalid delay %v",d)}}}
