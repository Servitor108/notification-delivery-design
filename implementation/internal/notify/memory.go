package notify

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"
)

type MemoryRepository struct { mu sync.Mutex; byID map[string]Notification; byKey map[string]string }
func NewMemoryRepository()*MemoryRepository{return &MemoryRepository{byID:map[string]Notification{},byKey:map[string]string{}}}
func (r *MemoryRepository) Create(ctx context.Context,n Notification)(Notification,bool,error){
	if err:=ctx.Err();err!=nil{return Notification{},false,err};r.mu.Lock();defer r.mu.Unlock()
	if id,ok:=r.byKey[n.IdempotencyKey];ok{old:=r.byID[id];if old.Recipient!=n.Recipient||old.Channel!=n.Channel||old.Template!=n.Template{return Notification{},false,ErrDuplicate};return old,false,nil}
	r.byID[n.ID]=n;r.byKey[n.IdempotencyKey]=n.ID;return n,true,nil
}
func (r *MemoryRepository) Get(ctx context.Context,id string)(Notification,error){
	if err:=ctx.Err();err!=nil{return Notification{},err};r.mu.Lock();defer r.mu.Unlock();n,ok:=r.byID[id];if !ok{return Notification{},ErrNotFound};return n,nil
}
func (r *MemoryRepository) Claim(ctx context.Context)(Notification,error){
	if err:=ctx.Err();err!=nil{return Notification{},err};r.mu.Lock();defer r.mu.Unlock()
	ids:=make([]string,0,len(r.byID));for id:=range r.byID{ids=append(ids,id)};sort.Strings(ids)
	now:=time.Now()
	for _,id:=range ids{n:=r.byID[id];if n.Status==StatusQueued||n.Status==StatusRetryWait&& !n.NextAttemptAt.After(now){n.Attempts++;n.Status="processing";r.byID[id]=n;return n,nil}}
	return Notification{},ErrNotFound
}
func (r *MemoryRepository) Complete(ctx context.Context,id,msg string)error{
	if err:=ctx.Err();err!=nil{return err};r.mu.Lock();defer r.mu.Unlock();n,ok:=r.byID[id];if !ok{return ErrNotFound};if msg==""{n.Status=StatusAccepted}else{n.Status=StatusFailed;n.LastError=msg};r.byID[id]=n;return nil
}
func (r *MemoryRepository) Reschedule(ctx context.Context,id string,delay time.Duration,msg string)error{
	if err:=ctx.Err();err!=nil{return err};r.mu.Lock();defer r.mu.Unlock();n,ok:=r.byID[id];if !ok{return errors.New("notification not found")};n.Status=StatusRetryWait;n.NextAttemptAt=time.Now().Add(delay);n.LastError=msg;r.byID[id]=n;return nil
}
