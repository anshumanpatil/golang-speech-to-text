# Golang Speech To Text

```diff
- text in red
+ text in green
! text in orange
# text in gray
@@ text in purple (and bold)@@
```

## Run All in one
```sh
git clone https://github.com/anshumanpatil/golang-speech-to-text.git
cd golang-speech-to-text.git
docker-compose up -d
go mod tidy
pip install -r requirements.txt
go run main.go -mode all
```

## Run Python [Separately]
```sh
git clone https://github.com/anshumanpatil/golang-speech-to-text.git
cd golang-speech-to-text.git
docker-compose up -d
pip install -r requirements.txt
python3 speech.py
```

## Run Go [Separately]
```sh
git clone https://github.com/anshumanpatil/golang-speech-to-text.git
cd golang-speech-to-text.git
docker-compose up -d
go mod tidy
go run main.go
```

## License

MIT
