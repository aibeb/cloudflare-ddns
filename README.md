# cloudflare-ddns
* [x] Deploy by golang
* [x] Image size 13.3m
* [x] support arm

# run

```shell
docker run \
  -e ZONE_NAME=example.com \
  -e SUB_DOMAIN=test.example.com \
  -e CF_API_KEY=........ \
  -e CF_API_EMAIL= \
  aibeb/cloudflare-ddns
```