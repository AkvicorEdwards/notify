package encryption

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
)

// --- AES ---

// =================== CBC ======================
func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte, err error) {
	if !(len(key)==16 || len(key)==24 || len(key)==32) {
		return []byte(""), errors.New("aes key must be either 16, 24, or 32 bytes")
	}
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""),err
	}
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted, nil
}

func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte, err error) {
	if !(len(key)==16 || len(key)==24 || len(key)==32) {
		return []byte(""), errors.New("aes key must be either 16, 24, or 32 bytes")
	}
	block, err := aes.NewCipher(key)                              // 分组秘钥
	if err != nil {
		return []byte(""),err
	}
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return decrypted, nil
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// =================== ECB ======================
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte, err error) {
	if !(len(key)==16 || len(key)==24 || len(key)==32) {
		return []byte(""), errors.New("aes key must be either 16, 24, or 32 bytes")
	}
	cipherText, err := aes.NewCipher(generateKey(key))
	if err != nil {
		return []byte(""),err
	}
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipherText.BlockSize(); bs <= len(origData); bs, be = bs+cipherText.BlockSize(), be+cipherText.BlockSize() {
		cipherText.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted, nil
}

func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte, err error) {
	if !(len(key)==16 || len(key)==24 || len(key)==32) {
		return []byte(""), errors.New("aes key must be either 16, 24, or 32 bytes")
	}
	cipherText, err := aes.NewCipher(generateKey(key))
	if err != nil {
		return []byte(""),err
	}
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipherText.BlockSize(); bs < len(encrypted); bs, be = bs+cipherText.BlockSize(), be+cipherText.BlockSize() {
		cipherText.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim], err
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// =================== CFB ======================
func AesEncryptCFB(origData []byte, key []byte) (encrypted []byte, err error) {
	if !(len(key)==16 || len(key)==24 || len(key)==32) {
		return []byte(""), errors.New("aes key must be either 16, 24, or 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""),err
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte(""),err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted, nil
}

func AesDecryptCFB(encrypted []byte, key []byte) (decrypted []byte, err error) {
	if !(len(key)==16 || len(key)==24 || len(key)==32) {
		return []byte(""), errors.New("aes key must be either 16, 24, or 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""),err
	}
	if len(encrypted) < aes.BlockSize {
		return []byte(""), errors.New("cipherText too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted, nil
}

// --- DES ---

func DesEncrypt(origData []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pkcs5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func DesDecrypt(crypted []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

// 3DES加密
func DesTripleEncrypt(origData, key []byte, iv []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pkcs5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 3DES解密
func DesTripleDecrypt(encrypted, key []byte, iv []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(encrypted))
	blockMode.CryptBlocks(origData, encrypted)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

func zeroPadding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{0}, padding)
	return append(cipherText, padText...)
}

func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

// --- RSA ---

// GenRsaKey 生成 PKCS1私钥、PKCS8私钥和公钥文件
func GenRsaKey(bits int) error {
	//生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private_key.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	//生成PKCS8私钥
	pk8Stream, _ := MarshalPKCS8PrivateKey(derStream)
	block = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pk8Stream,
	}
	file, err = os.Create("pkcs8_private_key.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	//生成公钥文件
	publicKey := &privateKey.PublicKey
	defPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: defPkix,
	}
	file, err = os.Create("public_key.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

// MarshalPKCS8PrivateKey 由私钥获取PKCS8公钥 这种方式生成的PKCS8与OpenSSL转成的不一样，但是BouncyCastle里可用
func MarshalPKCS8PrivateKey(key []byte) ([]byte, error) {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = key

	k, err := asn1.Marshal(info)
	if err != nil {
		return []byte(""), err
	}
	return k, nil
}

// 由私钥获取PKCS8公钥
func MarshalPKCS8PrivateKey1(key *rsa.PrivateKey) ([]byte, error) {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)

	k, err := asn1.Marshal(info)
	if err != nil {
		return []byte(""), nil
	}
	return k, nil
}

type PriKeyType uint

const (
	PKCS1 PriKeyType = iota
	PKCS8
)

//私钥签名
func RsaSign(data, privateKey []byte, keyType PriKeyType) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	priv, err := getPriKey(privateKey, keyType)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed)
}

//公钥验证
func RsaSignVer(data, signature, publicKey []byte) error {
	hashed := sha256.Sum256(data)
	//获取公钥
	pub, err := getPubKey(publicKey)
	if err != nil {
		return err
	}
	//验证签名
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
}

// 公钥加密
func RsaEncrypt(data, publicKey []byte) ([]byte, error) {
	//获取公钥
	pub, err := getPubKey(publicKey)
	if err != nil {
		return nil, err
	}
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// 私钥解密,privateKey为pem文件里的字符
func RsaDecrypt(encData, privateKey []byte, keyType PriKeyType) ([]byte, error) {
	//解析PKCS1a或者PKCS8格式的私钥
	priv, err := getPriKey(privateKey, keyType)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, encData)
}

func getPubKey(publicKey []byte) (*rsa.PublicKey, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	if pub, ok := pubInterface.(*rsa.PublicKey); ok {
		return pub, nil
	} else {
		return nil, errors.New("public key error")
	}
}

func getPriKey(privateKey []byte, keyType PriKeyType) (*rsa.PrivateKey, error) {
	//获取私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	var priKey *rsa.PrivateKey
	var err error
	switch keyType {
	case PKCS1:
		{
			priKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		}
	case PKCS8:
		{
			prkI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			priKey = prkI.(*rsa.PrivateKey)
		}
	default:
		{
			return nil, errors.New("unsupport private key type")
		}
	}
	return priKey, nil
}

func Base64StringToByte(str string) []byte {
	b, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return []byte{}
	}
	return b
}
func Base64StringToString(str string) string {
	data, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}
func Base64ByteToByte(str []byte) []byte {
	b, err := base64.RawURLEncoding.DecodeString(string(str))
	if err != nil {
		return []byte{}
	}
	return b
}
func Base64ByteToString(str []byte) string {
	b, err := base64.RawURLEncoding.DecodeString(string(str))
	if err != nil {
		return ""
	}
	return string(b)
}

func ByteToBase64Byte(str []byte) []byte {
	b := base64.RawURLEncoding.EncodeToString(str)
	return []byte(b)
}
func ByteToBase64String(str []byte) string {
	return base64.RawURLEncoding.EncodeToString(str)
}
func StringToBase64Byte(str string) []byte {
	b := base64.RawURLEncoding.EncodeToString([]byte(str))
	return []byte(b)
}
func StringToBase64String(str string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(str))
}