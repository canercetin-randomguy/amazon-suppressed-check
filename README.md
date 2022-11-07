# club-noira
Amazon Product Lookup

todo: get the signature to fucking work. https://stackoverflow.com/questions/74334445/confused-about-signature-creation-process-in-golang
nah.

"canonical url" and "string to sign" is 1:1 match with what amazon expects. there is a problem with

```go
signatureString = HMACSha256([]byte("AWS4"+kSecret), []byte(time.Now().UTC().Format("20060102"))) // FORMATLAMA DOĞRU
	signatureString = HMACSha256(signatureString, []byte("us-east-1"))                                // BU AMINAKoyDUĞUMUN EVLADI DOĞRU
	signatureString = HMACSha256(signatureString, []byte("execute-api"))                              // BU SATIR BENİM ANAM OROSPU BENİM YÜZÜNDEN PROGRAM ÇALIŞMIYOR DİYOR
	signatureString = HMACSha256(signatureString, []byte(AwsRequestType))                             // BU OROSPU ENİĞİ DOĞRU
```

and i cant figure it out.

among all the things I did, Amazon pushed my sanity farthest among the other things. congrulations.

**PLEASE DONT RECOMMEND ME AWS SDK FOR GOLANG, IT STRAIGHT DOESNT WORK WITH SP-API.**

every information is correct when requested with postman, but signature calculation is wrong with main.go.
