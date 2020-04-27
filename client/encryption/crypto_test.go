package encryption

import (
	"testing"
)

var testsString = []string{
	"asdfasfawav",
	"alksdjflkj9896178567][;,;;o;k",
	"alksjf12ilusd87q2oj:KI)(@&TYGHJBKNSP(8^&%$#UIuhjbksajfiuhiwe",
	"aksjdf lakjsdf lkajwl kja lkaLKJA LKJFDLNQOIUOPOPQIEHJBNKLVPBOYQUWIB",
	"KLJ!*&#^%^&^*()!#HBJHVBYUIEUFHIOP{}}><Moihugy1tf7ufvhbsnbaysfvbqooidguhlzmq0ohlBNVCXSQAZXCVHJ~!K#I&*^&%*)+PL",
	"在科技阿里夫阿拉克加快了我就法奥",
	"卢卡奇得粉碎了卡了科技哦求uioshdf计划",
	"、【】；，；。/‘/】、97867631@#",
	"￥%……&*（——）（（）——+、】阿毘 ",
	"阿克拉加反抗军情况了解孔子了继续咯iqu我日军",
	"阿卡解放卡拉经历过哈利卡浃髓沦肤简历库将来副前往进而快乐解放的了！ ",
	" 卢卡角色；离开家里 了卡捷算法\nakljsf ",
	"klzjx cvklaj slwi 将阿喀琉斯的风口浪尖看",
	"起来将军阿健康的法拉科技岁的法哦i法赛季的快递费拉角色的分类可建立起",
	"！@#￥%……&*（）——～+『：》』“|』》！（*&%",
	"……￥#@！￥%￥……&*（（\n\r \r \n\n\rlkaj\rakj	kjla lk",
	`卡捷拉
阿卡解

	卢卡奇峰i？劳`,
	`akj2lk jaklj2io asjdf lkj科技阿斯兰的反抗克拉克laksdjf `,
}

var testsKey = struct {
	AES []string
	DES []string
	DES3 []string
	RSA []struct{Public, Private string}
}{
	[]string{
		// 16 -> AES-128
		`1234567812345678`,
		`lkajliqlkznvlkaj`,
		`11kljli1nkbkauy1`,
		`!@#$%^&*()_:>LP[`,
		`kaj187BKHTA6{"'2`,
		"`kjbi8okn`bkjgut",
		// 24 -> AES-192
		`123456781234567812345678`,
		`l;ak;1l3pl;'""!OPUio1nln`,
		`><{:>}{!_)(*)*(&^#1ekfak`,
		`anlvajkquayvjozknxvb,mlw`,
		`klajo2iuoaiuo2ijrokalkfm`,
		`klajliwjvnlawuflaklwklkn`,
		// 32 -> AES-256
		`12345678123456781234567812345678`,
		`kajli2jlakjdflakjwiulkajlkanoi3u`,
		`kjakj2iuoiqoij34snvkqqoiui982637`,
		`kjljoi38978jajyiauywiuanlklkajss`,
		`lakjlijxlvijiwulaksdnalskdfjlaks`,
		`!@#$%^&*()_:>LP[kaj187BKHTA6{"'2`,
		`kaljfjlijlanlvkajlifulksadfnaajs`,
	},
	[]string{
		`12345678`,
		`asdf2ava`,
		`sdfwhpmz`,
		`*&%!#*(n`,
		`%^!HKLaf`,
		`}>:JISYI`,
		`LKUskvnk`,
		`29ujslkn`,
	},
	[]string{
		`123456781234567812345678`,
		`asdf2avaasf2dfgq3q6sdgas`,
		`sdfwhpmzmznjkwnalhakjwlz`,
		`*&%!#*(nIUY@*&Gkhakjsasf`,
		`%^!HKLafvn.a,mfk}>:""awb`,
		`}>:JISYILIULWNLlkjaoioha`,
		`LKUskvnkasnliui2368hkjak`,
		`29ujslknavoiu2o8aoklkajf`,
	},
	// 私钥生成
	//openssl genrsa -out rsa_private_key.pem 1024
	// 公钥: 根据私钥生成
	//openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
	[]struct{Public, Private string}{
		{
			`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDcGsUIIAINHfRTdMmgGwLrjzfM
NSrtgIf4EGsNaYwmC1GjF/bMh0Mcm10oLhNrKNYCTTQVGGIxuc5heKd1gOzb7bdT
nCDPPZ7oV7p1B9Pud+6zPacoqDz2M24vHFWYY2FbIIJh8fHhKcfXNXOLovdVBE7Z
y682X1+R1lRK8D+vmQIDAQAB
-----END PUBLIC KEY-----
`,
			`
-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDcGsUIIAINHfRTdMmgGwLrjzfMNSrtgIf4EGsNaYwmC1GjF/bM
h0Mcm10oLhNrKNYCTTQVGGIxuc5heKd1gOzb7bdTnCDPPZ7oV7p1B9Pud+6zPaco
qDz2M24vHFWYY2FbIIJh8fHhKcfXNXOLovdVBE7Zy682X1+R1lRK8D+vmQIDAQAB
AoGAeWAZvz1HZExca5k/hpbeqV+0+VtobMgwMs96+U53BpO/VRzl8Cu3CpNyb7HY
64L9YQ+J5QgpPhqkgIO0dMu/0RIXsmhvr2gcxmKObcqT3JQ6S4rjHTln49I2sYTz
7JEH4TcplKjSjHyq5MhHfA+CV2/AB2BO6G8limu7SheXuvECQQDwOpZrZDeTOOBk
z1vercawd+J9ll/FZYttnrWYTI1sSF1sNfZ7dUXPyYPQFZ0LQ1bhZGmWBZ6a6wd9
R+PKlmJvAkEA6o32c/WEXxW2zeh18sOO4wqUiBYq3L3hFObhcsUAY8jfykQefW8q
yPuuL02jLIajFWd0itjvIrzWnVmoUuXydwJAXGLrvllIVkIlah+lATprkypH3Gyc
YFnxCTNkOzIVoXMjGp6WMFylgIfLPZdSUiaPnxby1FNM7987fh7Lp/m12QJAK9iL
2JNtwkSR3p305oOuAz0oFORn8MnB+KFMRaMT9pNHWk0vke0lB1sc7ZTKyvkEJW0o
eQgic9DvIYzwDUcU8wJAIkKROzuzLi9AvLnLUrSdI6998lmeYO9x7pwZPukz3era
zncjRK3pbVkv0KrKfczuJiRlZ7dUzVO0b6QJr8TRAA==
-----END RSA PRIVATE KEY-----
`,
		},
	},
}

func Test_AesCBC_1(t *testing.T) {
	for _, tt := range testsString {
		for _, key := range testsKey.AES {
			encrypt, _ := AesEncryptCBC([]byte(tt), []byte(key))
			decrypt, _ := AesDecryptCBC(encrypt, []byte(key))
			if string(decrypt) != tt {
				t.Errorf("AES_CBC(%s, %s)", tt, string(decrypt))
			}
		}
	}
}

func Test_AesCFB_1(t *testing.T) {
	for _, tt := range testsString {
		for _, key := range testsKey.AES {
			encrypt, _ := AesEncryptCFB([]byte(tt), []byte(key))
			decrypt, _ := AesDecryptCFB(encrypt, []byte(key))
			if string(decrypt) != tt {
				t.Errorf("AES_CFB(%s, %s)", tt, string(decrypt))
			}
		}
	}
}

func Test_AesECB_1(t *testing.T) {
	for _, tt := range testsString {
		for _, key := range testsKey.AES {
			encrypt, _ := AesEncryptECB([]byte(tt), []byte(key))
			decrypt, _ := AesDecryptECB(encrypt, []byte(key))
			if string(decrypt) != tt {
				t.Errorf("AES_ECB(%s, %s)", tt, string(decrypt))
			}
		}
	}
}

func Test_Des_1(t *testing.T) {
	for _, tt := range testsString {
		for _, key := range testsKey.DES {
			encrypt, _ := DesEncrypt([]byte(tt), []byte(key), []byte(key)[:8])
			decrypt, _ := DesDecrypt(encrypt, []byte(key), []byte(key)[:8])
			if string(decrypt) != tt {
				t.Errorf("DES(%s, %s)", tt, string(decrypt))
			}
		}
	}
}

func Test_3Des_1(t *testing.T) {
	for _, tt := range testsString {
		for _, key := range testsKey.DES3 {
			encrypt, _ := DesTripleEncrypt([]byte(tt), []byte(key), []byte(key)[:8])
			decrypt, _ := DesTripleDecrypt(encrypt, []byte(key), []byte(key)[:8])
			if string(decrypt) != tt {
				t.Errorf("DES Triple(%s, %s)", tt, string(decrypt))
			}
		}
	}
}

func Test_Rsa_1(t *testing.T) {
	for _, tt := range testsString {
		for _, key := range testsKey.RSA {
			encrypt, _ := RsaEncrypt([]byte(tt), []byte(key.Public))
			decrypt, _ := RsaDecrypt(encrypt, []byte(key.Private), PKCS1)
			if string(decrypt) != tt {
				t.Errorf("RSA(%s, %s)", tt, string(decrypt))
			}
		}
	}
}

func Test_RsaSign_1(t *testing.T) {
	for _, tt := range testsString {
		for _, key := range testsKey.RSA {
			sign, _ := RsaSign([]byte(tt), []byte(key.Private), PKCS1)
			err := RsaSignVer([]byte(tt), sign, []byte(key.Public))
			if err != nil {
				t.Errorf("Failed %s", tt)
			}
		}
	}
}

