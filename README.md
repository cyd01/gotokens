# gotokens

A light tokens manager

## How to

### Authenticate with basic-auth

Use `POST /auth` to pass basic authentication and get valid token in response `Token` cookie. If succeeded the HTTP return code is `201`:

Example:

```bash
$ curl -v -X POST -H 'Authorization: Basic '$(printf "admin:pass" | base64) http://127.0.0.1:8080/tokens/auth
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> POST /tokens/auth HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.68.0
> Accept: */*
> Authorization: Basic YWRtaW46cGFzcw==
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Set-Cookie: Token=GQSSCMBG-cb665a8c705fd1d09c38a3ea39d4a48bad91ee51b143a29318cea8dcc6d9c8b8; Path=/; Domain=127.0.0.1; Max-Age=300; HttpOnly
< Date: Mon, 03 Oct 2022 14:09:35 GMT
< Content-Length: 208
< 
* Connection #0 to host 127.0.0.1 left intact
{"id":"36da06fd-9dcd-47de-a20a-742d054962a7","user":"admin","token":"cb665a8c705fd1d09c38a3ea39d4a48bad91ee51b143a29318cea8dcc6d9c8b8","address":"127.0.0.1","created":1664806175,"updated":1664806175,"hits":0}
```

To check if token is valid just call `GET /validate/:token`. If token is valid the return code is `200`:

```bash
$ curl -v http://127.0.0.1:8080/tokens/validate/313a5bc6e1a6e96e1cc0cbf0f828f0d0a7521f0b985e65e87d0b66375fb14fc6
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> GET /tokens/validate/313a5bc6e1a6e96e1cc0cbf0f828f0d0a7521f0b985e65e87d0b66375fb14fc6 HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Mon, 03 Oct 2022 14:16:47 GMT
< Content-Length: 46
< 
* Connection #0 to host 127.0.0.1 left intact
{"message":"Valid token","status":"succeeded"}
```

### Authenticate with challenge data

First you need to get the challenge data with `GET /challengedata` call:

```bash
$ curl -v http://127.0.0.1:8080/tokens/challengedata
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> GET /tokens/challengedata HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Set-Cookie: ChallengeData=e3502133-d019-428d-9145-93f6bcb7d1b7; Path=/; Domain=127.0.0.1; Max-Age=300; HttpOnly
< Date: Mon, 03 Oct 2022 14:20:18 GMT
< Content-Length: 1044
< 
* Connection #0 to host 127.0.0.1 left intact
{"challengedata":"75ba6aecc68008e4c2e58548409f04690dd3feca4ebe024a3ccf5374e7bcf49d5c0dbd882c8649239cd872b780ec56fee0ac42c158c77b74b26fe6a9949313bce7992940c8358a2d7b16bb50da1b8440fa4b1448badd1e3e7513ea814cd84a3b14709090a69d92253977387c42af533c1100dab9f5363b1a6baf4dbac3f2cf359aadaca1c23eb298c1f344881102ee6da67e88d4668fe9686adae2a83761b8a4154b6e41f976da0febd4d71b2a57edc3a1a5fccd939db56a22a1db856dc58c53d4e4186a883da75b393a530f8c1f249f3822bdba68de4c2557dc333389f5610b47f9fbb93c49692a700a3b3f348ea1f40e6ec22c2ff2f9d50e95c5a52174122900397ab26e4b2f35f6e87ea82d8c26a7cb9a8c8ba7761220401eedd990bb5a2ff5f8c1512f7651f879226be5cc2987a26e5c0d8fd431512c0c1d48fb763c9d61eaa20b21c39f5d56426b443a1cdce3b447a755590c12edb5edcd8a3e74ece27d062d1d4e6905fe04823bd87384ffc0c9c8b69726a39b130aa853008a5ab74b1bc247d1de9eeb0d5568625507db468d45d03db4fb9e595034a3e42925cc71875e16823ccfb3515fabbb2ba5d9d1845c298651a728082100ef356b884d0202fd90118963c2f40c50de6f697d6c359c47edc9c0b58ebe98397d4d6c9f554974fa6abafa8c794f8addad2f1e756df3c0ddcd0864ffd87b67f3f7dd96d325d873f407"}
```

In the response, two things:
- a challenge data short id into a cookie
- the real challenge data in the response json body

You must calculate the password hash with the formula `md5(password+challengeData)`

Then call `POST /` with a json body `{"login":"admin", "password":"caclucatedHash"}` with the challenge data id int the `ChallengeData` cookie. If succeeded the HTTP return code is `201` and the token is present in Cookie.

Full curl example:

```bash
$ chd=$(curl -v http://127.0.0.1:8080/tokens/challengedata 2> /tmp/out| jq -r .challengedata|tr -d '\n')
$ cat /tmp/out
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> GET /tokens/challengedata HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Set-Cookie: ChallengeData=0cb93b6c-46fb-4c86-8c55-983672670526; Path=/; Domain=127.0.0.1; Max-Age=300; HttpOnly
< Date: Mon, 03 Oct 2022 14:54:35 GMT
< Content-Length: 1044
< 
{ [1044 bytes data]
100  1044  100  1044    0     0  1019k      0 --:--:-- --:--:-- --:--:-- 1019k
* Connection #0 to host 127.0.0.1 left intact
$ id=$(grep "ChallengeData=" /tmp/out|sed 's/;.*$//'|sed 's/^.*=//')
$ rm -f /tmp/out
$ echo $id
0cb93b6c-46fb-4c86-8c55-983672670526
$ echo $chd
0183a54386226605fcf6ac5b4ee4c81213756c3a080379b4384a499db9dcd4d9a03b2b28c35bd7b577e19eb383886826a2884e66836b07682da98799840770a14a2b6e244403a27c6dd34faabf6c5aa9ff84d4dad308035fff411070932ab63610d79b338f503c51f242eb0e47e7b773d2711ce78650a9d2d353b70730f96940013e66905517641411935b19c88c937be563257094ce9e825dd0bcee02f5c8022225b968065b5ce38dc7108ece9ea3911357b832baccd3582aac2bb56643cd2e2831d034aa5b11ae41a92c7e1a19af5dadaee14778fc482061bd59e3dbeacaa16f925ffbd835c376b74233e078853a517b4ec5d5afd536a6cd436bae3b267d777f702628b8eea23140513e2e2b036ee27e66f8f31c44bf1a7aee4e9e4a3114085fa12be7bf62d471a4a1044a04f18405ee97f04db4630aca012a8d571a98721d8e08392c4de1d9a54439626fa51e6b5cde2baf40a3d608449762626e39c95850ba7121b8b7dda1a0cd1f14047886bd47435efe685325cf5699ea44984ff5d6bf9a42e9ed95bc6daebacf480350409ab460448f0e5d6df08b064e56bc7fc0c76e98153552e43083cf21e4226b0ee93070770152e2d19d092e2414aa308b75d778d242b1fa41e5e913580f8797f447e78dcea8db683386be1ff726773317eaee2b3e505f651c6f1$ 31eff47db77dd7e88450262092afbe90a1223f406d5539d9dc6
$ password=$(printf "pass"$chd|md5sum| sed 's/ .*$//')
$ echo $password
14e4c473aee64b7067663674ecca2027
$ url -v -L -H 'Content-type:application/json' -H 'Cookie: ChallengeData='$id 'http://127.0.0.1:8080/tokens/' -d '{"login":"admin","password":"'${password}'"}'
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> POST /tokens/ HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.68.0
> Accept: */*
> Content-type:application/json
> Cookie: ChallengeData=0cb93b6c-46fb-4c86-8c55-983672670526
> Content-Length: 63
> 
* upload completely sent off: 63 out of 63 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Set-Cookie: ChallengeData=; Path=/; Domain=127.0.0.1; Max-Age=0; HttpOnly
< Set-Cookie: Token=GQSSCMBG-8554b9790d156db679ece1febcb3e571a23bd42d9d30ff66fb894afb392d15e6; Path=/; Domain=127.0.0.1; Max-Age=300; HttpOnly
< Date: Mon, 03 Oct 2022 14:56:12 GMT
< Content-Length: 208
< 
* Connection #0 to host 127.0.0.1 left intact
{"id":"b1515d09-e439-4e67-98e1-20f0e854f1fd","user":"admin","token":"8554b9790d156db679ece1febcb3e571a23bd42d9d30ff66fb894afb392d15e6","address":"127.0.0.1","created":1664808972,"updated":1664808972,"hits":0}
```

The same with curl cookie management feature:

```bash
$ chd=$(curl -v --cookie-jar cookies.jar http://127.0.0.1:8080/tokens/challengedata | jq -r .challengedata|tr -d '\n')
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> GET /tokens/challengedata HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
* Added cookie ChallengeData="9e0d53c9-f5d7-4980-9b95-5c1290f5c38d" for domain 127.0.0.1, path /, expire 1664809733
< Set-Cookie: ChallengeData=9e0d53c9-f5d7-4980-9b95-5c1290f5c38d; Path=/; Domain=127.0.0.1; Max-Age=300; HttpOnly
< Date: Mon, 03 Oct 2022 15:03:53 GMT
< Content-Length: 1044
< 
{ [1044 bytes data]
100  1044  100  1044    0     0  1019k      0 --:--:-- --:--:-- --:--:-- 1019k
* Connection #0 to host 127.0.0.1 left intact
$ password=$(printf "pass"$chd|md5sum| sed 's/ .*$//')
$ echo $password
cfd638e13ed2ca7db62a9fa89848d7cf
$ curl -v -L -H 'Content-type:application/json' --cookie cookies.jar 'http://127.0.0.1:8080/tokens/' -d '{"login":"admin","password":"'${password}'"}'
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> POST /tokens/ HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.68.0
> Accept: */*
> Cookie: ChallengeData=9e0d53c9-f5d7-4980-9b95-5c1290f5c38d
> Content-type:application/json
> Content-Length: 63
> 
* upload completely sent off: 63 out of 63 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
* Replaced cookie ChallengeData="" for domain 127.0.0.1, path /, expire 1
< Set-Cookie: ChallengeData=; Path=/; Domain=127.0.0.1; Max-Age=0; HttpOnly
* Added cookie Token="GQSSCMBG-bcaf61b6f913bb26e5e23a82e968748870fdcd23d96b66d8878e97306eb96349" for domain 127.0.0.1, path /, expire 1664809773
< Set-Cookie: Token=GQSSCMBG-bcaf61b6f913bb26e5e23a82e968748870fdcd23d96b66d8878e97306eb96349; Path=/; Domain=127.0.0.1; Max-Age=300; HttpOnly
< Date: Mon, 03 Oct 2022 15:04:33 GMT
< Content-Length: 208
< 
* Connection #0 to host 127.0.0.1 left intact
{"id":"2ea3787c-a11e-42dc-8a97-5e63995204e0","user":"admin","token":"bcaf61b6f913bb26e5e23a82e968748870fdcd23d96b66d8878e97306eb96349","address":"127.0.0.1","created":1664809473,"updated":1664809473,"hits":0}
$ rm -f cookies.jar
```

