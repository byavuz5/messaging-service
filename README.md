# Messaging Service

Proje örnek bir mesajlaşma uygulamasının servislerini içeriyor. Golang ile oluşturulmuş servislerin docker ile image haline getirilip, AWS Lambda üzerinde çalıştırılmasını sağlıyor. Projede veritabanı olarak DynamoDB kullanılmıştır.

***
## AWS CDK İçeriği

Uygulamanın çalışması için gerekli olan sistemi AWS tarafında ayağa kaldırıyor.

* API Gateway oluşturulması.
* Servislerin kullanacağı DynamoDB tablolarının oluşturulması.
* Lambda servislerinin oluşturulması ve ilgili tablolara erişimi için izinlerin eklenmesi.
* API Gateway'e resource eklenmesi ve ilgili lambda servislerinin bu resource'a bağlanması.

***
## CI/CD İçeriği

Git repository'sine push işlemi geldiğinde Github Action iş akışı tetikleniyor ve ilk olarak docker ile golang servislerinden build alıp, docker image'lerini oluşturuyor. Daha sonra bu image'leri AWS ECR servisine gönderiyor. Tüm build ve gönderme işlemleri bittikten sonra AWS CDK ile yeni image'ler ile lambda servislerini güncelliyor.

***
## Servisler

* [POST]&emsp;/createAccount
> Kullanıcıların hesap oluşturmasını sağlıyor. Parametre olarak username(String) ve password(String) JSON objesi alıyor. Çıktı olarak bilgi ya da hata mesajı içeren JSON objesi dönüyor.
 
* [POST]&emsp;/login
> Kullanıcıların sisteme giriş yapmasını sağlıyor. Parametre olarak username(String) ve password(String) JSON objesi alıyor. Çıktı olarak bilgi ya da hata mesajı içeren JSON objesi dönüyor.

* [POST]&emsp;/sendMessage
> Kullanıcıların başka bir kullanıcıya mesaj göndermesini sağlıyor. Parametre olarak sender(String), sent_to(String) ve message(String) JSON objesi alıyor. Çıktı olarak bilgi ya da hata mesajı içeren JSON objesi dönüyor.

* [POST]&emsp;/getAllMessages
> Kullanıcıların tüm kullanıcılara gönderdiği ve aldığı mesajların geçmişini listeleliyor. Parametre olarak username(String) JSON objesi alıyor. Çıktı olarak, tarihe göre sıralı bir şekilde mesajı gönderen kişi, mesaj içeriği ve gönderildiği tarihi ya da hata mesajı içeren JSON objesi dönüyor.

* [POST]&emsp;/getContactMessages
> Kullanıcıların belirli bir kullanıcıya gönderdiği ve aldığı mesajların geçmişini listeliyor. Parametre olarak username(String) ve contactName(String) JSON objesi alıyor. Çıktı olarak, tarihe göre sıralı bir şekilde mesajı gönderen kişi, mesaj içeriği ve gönderildiği tarihi ya da hata mesajı içeren JSON objesi dönüyor.

* [POST]&emsp;/getActivityLogs
> Kullanıcıların gerçekleştirdiği aktiviteleri listeliyor. Parametre olarak username(String) JSON objesi alıyor. Çıktı olarak, kullanıcıların gerçekleştirdiği aktiviteleri ya da hata mesajı içeren JSON objesi dönüyor.

* [POST]&emsp;/getSystemLogs
> Sistemin çalışmasını engelleyen hata mesajlarının listelenmesini sağlıyor. Parametre olarak service_name(String) JSON objesi alıyor. Çıktı olarak servisin aldığı hatalar ya da hata mesajı içeren JSON objesi dönüyor.

***

## Örnek Servis Uygulamaları

* [POST]&emsp;/createAccount
> * Girdi : {"username":"ahmet", "password":"12345"}
> * Çıktı : {  
            &emsp;&emsp;&emsp;&emsp;"statusCode": 200,  
            &emsp;&emsp;&emsp;&emsp;"headers": {  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"Content-Type": "application/json"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;"body": {   
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"message": "User created, username: ahmet"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;"errors": null  
            &emsp;&emsp;&emsp;}  

* [POST]&emsp;/login
> * Girdi : {"username":"ahmet", "password":"12345"}
> * Çıktı : {  
            &emsp;&emsp;&emsp;&emsp;"statusCode": 200,  
            &emsp;&emsp;&emsp;&emsp;"headers": {  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"Content-Type": "application/json"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;"body": {  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"message": "Success login, username: ahmet"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;"errors": null  
            &emsp;&emsp;&emsp;}  

* [POST]&emsp;/sendMessage
> * Girdi : {"sender":"ahmet", "sent_to":"mehmet", "message":"Selam"}
> * Çıktı : {  
            &emsp;&emsp;&emsp;&emsp;"statusCode": 200,  
            &emsp;&emsp;&emsp;&emsp;"headers": {  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"Content-Type": "application/json"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;"body": {  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"message": "Message sended."  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;"errors": null  
            &emsp;&emsp;&emsp;}  

* [POST]&emsp;/getAllMessages
> * Girdi : {"username":"ahmet"}
> * Çıktı : {  
            &emsp;&emsp;&emsp;&emsp;"statusCode": 200,  
            &emsp;&emsp;&emsp;&emsp;"headers": {  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"Content-Type": "application/json"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;"body": [  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;{  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"sender": "ahmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"message": "Selam",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"room_id": "ahmet-mehmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"createdAt": "2022-08-30 12:50:46.650564835 +0000 UTC"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;}  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;],  
            &emsp;&emsp;&emsp;&emsp;"errors": null  
            &emsp;&emsp;&emsp;}  

* [POST]&emsp;/getContactMessages
> * Girdi : {"username":"ahmet", "contactName":"mehmet"}
> * Çıktı : {  
            &emsp;&emsp;&emsp;&emsp;"statusCode": 200,  
            &emsp;&emsp;&emsp;&emsp;"headers": {  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"Content-Type": "application/json"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;"body": [  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;{  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"sender": "ahmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"message": "Selam",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"room_id": "ahmet-mehmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"createdAt": "2022-08-30 12:50:46.650564835 +0000 UTC"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;{  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"sender": "mehmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"message": "Selam",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"room_id": "ahmet-mehmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"createdAt": "2022-08-30 14:40:10.999917592 +0000 UTC"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;}  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;],  
            &emsp;&emsp;&emsp;&emsp;"errors": null  
            &emsp;&emsp;&emsp;}  

* [POST]&emsp;/getActivityLogs
> * Girdi : {"username": "ahmet"}
> * Çıktı : {  
            &emsp;&emsp;&emsp;&emsp; "statusCode": 200,  
            &emsp;&emsp;&emsp;&emsp; "headers": {   
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp; "Content-Type": "application/json"  
            &emsp;&emsp;&emsp;&emsp; },  
            &emsp;&emsp;&emsp;&emsp; "body": [  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;{  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"username": "ahmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"activity": "Invalid login.",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"createdAt": "2022-08-30 09:09:10.1895531 +0000 UTC"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;{  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"username": "ahmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"activity": "Success login.",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"createdAt": "2022-08-30 09:09:21.706084 +0000 UTC"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;{  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"username": "ahmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"activity": "Send message to mehmet.",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"createdAt": "2022-08-30 12:50:46.897679742 +0000 UTC"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;{  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"username": "ahmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"activity": "List all messages.",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"createdAt": "2022-08-30 12:51:09.400160502 +0000 UTC"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;{  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"username": "ahmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"activity": "List messages with mehmet.",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"createdAt": "2022-08-30 12:52:15.904514812 +0000 UTC"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;{  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"username": "ahmet",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"activity": "List activity logs.",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"createdAt": "2022-08-30 12:52:38.87719108 +0000 UTC"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;],  
            &emsp;&emsp;&emsp;&emsp;"errors": null  
            &emsp;&emsp;&emsp;}  

* [POST]&emsp;/getSystemLogs
> * Girdi : {"service_name":"login"}
> * Çıktı : {  
            &emsp;&emsp;&emsp;&emsp;"statusCode": 200,  
            &emsp;&emsp;&emsp;&emsp;"headers": {  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"Content-Type": "application/json"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;},  
            &emsp;&emsp;&emsp;&emsp;"body": [  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;{  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"service_name": "login",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"err_message": "ValidationException: One or more parameter values are not valid. The AttributeValue for a key attribute cannot contain an empty string value. Key: username\n\tstatus code: 400, request id: FERPDNK0F1D7PCRA228R7DMUHBVV4KQNSO5AEMVJF66Q9ASUAAJG",  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;"createdAt": "2022-08-30 09:35:07.135860794 +0000 UTC"  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;}  
            &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;],  
            &emsp;&emsp;&emsp;&emsp;"errors": null  
            &emsp;&emsp;&emsp;}  



***
## Veritabanı Tasarımı

* ## users table
    ```
    | username(String) PK | password(String) | createdAt(String)                     |
    | ------------------- | ---------------- | ------------------------------------- |
    | ahmet               | 1234             | 2022-08-27 15:58:03.7699622 +0000 UTC |
    ```

* ## messages table
    ```
    | sender(String) PK   | message(String)  | room_id(String) | createdAt(String) SK                  |
    | ------------------- | ---------------- | --------------- | ------------------------------------- |
    | ahmet               | Selam            | ahmet-mehmet    | 2022-08-28 09:50:38.3772013 +0000 UTC |
    ```

* ## activity_logs table
    ```
    | username(String) PK | activity(String) | createdAt(String) SK                  |
    | ------------------- | ---------------- | ------------------------------------- |
    | ahmet               | Success login.   | 2022-08-28 09:54:03.9708315 +0000 UTC |
    ```

* ## system_logs table
    ```
    | service_name(String) PK | err_message(String)                                     | createdAt(String) SK                  |
    | ----------------------- | ------------------------------------------------------- | ------------------------------------- |
    | createAccount           | ResourceNotFoundException: Requested resource not found | 2022-08-27 18:12:43.4048887 +0000 UTC |
    ```

