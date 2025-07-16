# duguncom

## Kurulum

```
go get github.com/vahaponur/duguncom
```

## Fonksiyonlar

### `Login(username, password string, userType ...string) (LoginResponse, error)`
Düğün.com API'ye giriş yapar ve **Access Token** ile **Consumer Key** bilgilerini döner. Opsiyonel `userType` parametresi girilmezse varsayılan olarak `customer` kullanılır.

### `GetOfferRequest(login, params) (GetOfferResponse, error)`
Belirtilen tarih aralığında ( `createdAtStart`, `createdAtEnd` ) ve sayfalama bilgileriyle lead’leri (teklifleri) listeler.

### `SendSmsToCustomers(login, leadIds, message, waitTime)`
`leadIds` dizisindeki müşterilere SMS atar.
* `message`  : Gönderilecek metin (\n ile satır atlanabilir)
* `waitTime` : Her istekten sonra milisaniye cinsinden bekleme süresi (rate-limit korunur)

### `SendEmailToCustomers(login, leadIds, subject, body, waitTime)`
`leadIds` dizisindeki müşterilere **HTML** formatında e-posta atar.
* `subject`  : E-posta konusu
* `body`     : HTML gövde ( \<p\>, \<a\> vb. )
* `waitTime` : Her istekten sonra milisaniye cinsinden bekleme süresi

### `SendSmsOrEmailToCustomers(login, leads, smsMessage, emailSubject, emailBody, waitTime)`
`leads` yapısındaki her lead’in telefon bilgisine göre otomatik seçim yapar:
* Telefon numarası **boş / null** ise **Email** gönderir.
* Telefon numarası mevcut ise **SMS** gönderir.

## Örnek Kullanım

```go
package main

import (
    "fmt"
    "strconv"
    "github.com/vahaponur/duguncom"
)

func main() {
    // 1) Giriş
    login, err := duguncom.Login("email", "password")
    if err != nil {
        panic(err)
    }

    // 2) Lead’leri Çek
    offers, err := duguncom.GetOfferRequest(login, duguncom.GetOfferParams{
        Limit: "100",
        Start: "2025-07-15",
        End:   "2025-07-16",
        Page:  "1",
    })
    if err != nil {
        panic(err)
    }

    // 3) Otomatik SMS / Email Gönderimi
    const smsMsg = "Merhabalar.\nDüğün.com - L'invito Design olarak ulaşıyoruz.\nWebsitemiz: https://lainvito.com"
    const emailSubject = "Teklif İsteğiniz Hakkında: L'invito Design"
    const emailBody = `<p>L'invito Design ile iletişime geçtiğiniz için teşekkür ederiz.</p>\n<p>...</p>`

    if err := duguncom.SendSmsOrEmailToCustomers(login, offers.Data, smsMsg, emailSubject, emailBody, 100); err != nil {
        fmt.Println("Gönderim hataları:", err)
    }
}
```

## Testler
`test/` klasöründe temel senaryoları gösteren **login_test.go** dosyası bulunur. `go test ./...` komutu ile çalıştırabilirsiniz.

## Katkı
İstek ve katkılarınız için pull request gönderebilirsiniz.

## Lisans
MIT
