### for postgres version >15

GRANT ALL ON DATABASE mydb TO admin;
ALTER DATABASE mydb OWNER TO admin;

create .env file like .env.example

CREATE TABLE IF NOT EXISTS `emlak_paychy`

CREATE TABLE IF NOT EXISTS `mainpayly`

CREATE TABLE IF NOT EXISTS `meydan_umumy`

CREATE TABLE IF NOT EXISTS `min_hat`

CREATE TABLE IF NOT EXISTS `seneler`

CREATE TABLE IF NOT EXISTS `sertnamalar`

CREATE TABLE IF NOT EXISTS `forma_sobstv` // Can be implemented only on front

If you use Powershell, use below for deploy compiling
$Env:GOOS = "linux"; $Env:GOARCH = "amd64"

cd api-construction-share
GOOS=linux GOARCH=amd64 go build -o bin/payly
pscp bin/payly itadmin@192.168.0.10:/var/www/payly/api/
pscp .env itadmin@192.168.0.10:/var/www/payly/api/
plink itadmin@192.168.0.10 "sudo systemctl restart payly.service"

1q2w3e!@A98lk
