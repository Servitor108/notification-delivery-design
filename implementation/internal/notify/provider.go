package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPProvider struct { URL string; Token string; Client *http.Client }
func (p *HTTPProvider) Send(ctx Context,req ProviderRequest)(ProviderResult,error){
	body,err:=json.Marshal(req);if err!=nil{return ProviderResult{},err}
	request,err:=http.NewRequestWithContext(ctx.(context.Context),http.MethodPost,p.URL,bytes.NewReader(body));if err!=nil{return ProviderResult{},err}
	request.Header.Set("Content-Type","application/json");request.Header.Set("Authorization","Bearer "+p.Token);request.Header.Set("Idempotency-Key",req.NotificationID)
	client:=p.Client;if client==nil{client=&http.Client{Timeout:10*time.Second}}
	response,err:=client.Do(request);if err!=nil{return ProviderResult{},err};defer response.Body.Close()
	if response.StatusCode<200||response.StatusCode>=300{return ProviderResult{},providerError{code:response.StatusCode}}
	return ProviderResult{},nil
}
type providerError struct{code int}
func(e providerError)Error()string{return fmt.Sprintf("provider returned HTTP %d",e.code)}
func(e providerError)StatusCode()int{return e.code}
