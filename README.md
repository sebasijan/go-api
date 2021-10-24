# go-api

# ssl

Run below to generate ssl cert and key

    openssl genrsa -out go-api.key 2048
    openssl ecparam -genkey -name secp384r1 -out go-api.key
    openssl req -new -x509 -sha256 -key go-api.key -out go-api.crt -days 3650

Copy paths to config.json