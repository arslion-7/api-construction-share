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

### Deploy compiling

**Mac/Linux (Bash):**

```bash
cd api-construction-share
GOOS=linux GOARCH=amd64 go build -o bin/payly
scp bin/payly itadmin@192.168.0.10:/var/www/payly/api/
ssh itadmin@192.168.0.10 "sudo systemctl restart payly.service"
```

**Windows (Powershell):**

```powershell
cd api-construction-share
$Env:GOOS = "linux"; $Env:GOARCH = "amd64"
go build -o bin/payly
pscp bin/payly itadmin@192.168.0.10:/var/www/payly/api/
plink itadmin@192.168.0.10 "sudo systemctl restart payly.service"
```

1q2w3e!@A98lk

-- Дамп структуры для таблица paylynew.mainpayly
CREATE TABLE IF NOT EXISTS `mainpayly` (
`t_b` int(10) unsigned NOT NULL,
`min_hat` varchar(150) DEFAULT NULL,
`sene_hat_min_to_mud` date DEFAULT NULL,
`gurujy` varchar(255) DEFAULT NULL,
`paychy` varchar(100) DEFAULT NULL,
`sertnama_gurujy_paychy` varchar(100) DEFAULT NULL,
`desga` varchar(275) DEFAULT NULL,
`baha_umumy` varchar(20) DEFAULT '0.00',
`meydan_umumy` varchar(125) DEFAULT NULL,
`kep_resminama` varchar(225) DEFAULT NULL,
`emlak_paychy` varchar(200) DEFAULT NULL,
`baha_paychy` varchar(35) DEFAULT NULL,
`baha_1m2_paychy` varchar(25) DEFAULT NULL,
`salgy_desga` varchar(255) DEFAULT NULL,
`salgy_gurujy` varchar(150) DEFAULT NULL,
`salgy_paychy` varchar(150) DEFAULT NULL,
`bash_potr` varchar(100) DEFAULT NULL,
`sertnama_gur_potr` varchar(125) DEFAULT NULL,
`potratchy_komek` varchar(125) DEFAULT NULL,
`shahadatnama` varchar(175) DEFAULT NULL,
`ygtyyarnama` varchar(255) DEFAULT NULL,
`patent_pasport` varchar(255) DEFAULT NULL,
`sene_bashy_songy` varchar(50) DEFAULT NULL,
`sene_seredilen` date DEFAULT NULL,
`sene_hasaba_alnan` varchar(125) DEFAULT NULL,
`wezipe_alan_adam` varchar(125) DEFAULT NULL,
`ady_alan_adam` varchar(125) DEFAULT NULL,
`sene_san_sertnama` varchar(55) DEFAULT NULL,
`ady_paychy_alan` varchar(125) DEFAULT NULL,
`sene_paychy_alan` varchar(255) DEFAULT NULL,
`login` varchar(30) DEFAULT NULL,
PRIMARY KEY (`t_b`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Payly gurlushyk mudirliginin reyestrini yoretmek ucin programmasy';
