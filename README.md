# Algorift, PROS@FT 2024
[Documentație Algorift](https://github.com/lazarcloud/provocari-digitale/blob/main/design/docs.pdf)
Acest proiect a fost realizat de Lazar, un pasionat de informatică și provocări digitale. Puteți vizita portofoliul său la [lazar.lol](https://lazar.lol/).

Pentru rulare sunt necesare:

* Docker
* Golang
* Node.js
* NPM
* Compilator de C pentru dependențe

```sh
docker build -t cpp-executor -f .\CPPDockerfile .
```
consturieste imaginea compilatorului C++ pe docker
```sh
cd api && go run .
```
ruleaza serverul de Golang
```sh
cd app && npm i && npm run dev
```
ruleaza serverul web

Aplicație disponibilă pe **localhost:5173**



