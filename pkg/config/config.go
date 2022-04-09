package config

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/ascii85"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"filippo.io/age"

	"github.com/go-piv/piv-go/piv"
	"github.com/nwehr/npass"
	"github.com/nwehr/npass/pkg/cli"
	"github.com/nwehr/npass/pkg/repos"
	"gopkg.in/yaml.v3"
)

type ConfigOptions struct {
	Path   string
	Config Config
}

func GetConfigOptions(args []string) (ConfigOptions, error) {
	opts := ConfigOptions{
		Path: DefaultConfigFile,
	}

	for i, arg := range args {
		switch arg {
		case "--config":
			opts.Path = args[i+1]
		}
	}

	err := opts.Config.ReadFile(opts.Path)
	return opts, err
}

type Config struct {
	CurrentStore string        `yaml:"currentStore"`
	Stores       []StoreConfig `yaml:"stores"`
}

type StoreConfig struct {
	Name      string    `yaml:"name"`
	Location  string    `yaml:"location"`
	AuthToken *string   `yaml:"authToken,omitempty"`
	Age       AgeConfig `yaml:"age"`
}

type AgeConfig struct {
	Identity   string     `yaml:"identity"`
	PIV        *PIVConfig `yaml:"piv,omitempty"`
	Recipients []string   `yaml:"recipients"`
}

func (c AgeConfig) unlockIdentityWithPassword(password string) (age.Identity, error) {
	protector, err := age.NewScryptIdentity(password)
	if err != nil {
		return nil, fmt.Errorf("could not get protector identity: %v", err)
	}

	decoder := ascii85.NewDecoder(strings.NewReader(c.Identity))
	decrypter, err := age.Decrypt(decoder, protector)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt identity: %v", err)
	}

	identities, err := age.ParseIdentities(decrypter)
	if err != nil {
		return nil, fmt.Errorf("could not parse identity: %v", err)
	}

	return identities[0], nil
}

func (c AgeConfig) unlockIdentityWithPIV(yk *piv.YubiKey) (age.Identity, error) {
	cert, err := yk.Attest(slots[c.PIV.Slot])
	if err != nil {
		return nil, fmt.Errorf("could not get certificate: %v", err)
	}

	serial, _ := yk.Serial()

	auth := piv.KeyAuth{
		PINPrompt: cli.TTYPin(fmt.Sprintf("Card %d", serial)),
		PINPolicy: piv.PINPolicyOnce,
	}

	priv, err := yk.PrivateKey(slots[c.PIV.Slot], cert.PublicKey, auth)
	if err != nil {
		return nil, fmt.Errorf("could not setup private key: %v", err)
	}

	decrypter, ok := priv.(crypto.Decrypter)
	if !ok {
		return nil, fmt.Errorf("priv does not impliment Decrypter")
	}

	decoder := ascii85.NewDecoder(strings.NewReader(c.Identity))
	decoded, err := ioutil.ReadAll(decoder)
	if err != nil {
		return nil, fmt.Errorf("could not decode identity")
	}

	decrypted, err := decrypter.Decrypt(rand.Reader, decoded, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt key file: %v", err)
	}

	identity, err := age.ParseX25519Identity(string(decrypted))
	if err != nil {
		return nil, fmt.Errorf("could not parse identity: %v", err)
	}

	return identity, nil
}

func (c AgeConfig) UnlockIdentity() (age.Identity, error) {
	if c.PIV == nil {
		password, err := cli.TTYPassword()
		if err != nil {
			return nil, fmt.Errorf("could not get password: %v", err)
		}

		return c.unlockIdentityWithPassword(password)
	}

	yk, err := yubikey()
	if err != nil {
		return nil, fmt.Errorf("could not get card: %v", err)
	}

	return c.unlockIdentityWithPIV(yk)
}

func NewDefaultAgeConfig() (AgeConfig, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return AgeConfig{}, fmt.Errorf("could not generate new identity: %v", err)
	}

	password, err := cli.TTYPassword()
	if err != nil {
		return AgeConfig{}, fmt.Errorf("could not get password: %v", err)
	}

	protector, err := age.NewScryptRecipient(string(password))
	if err != nil {
		return AgeConfig{}, fmt.Errorf("could not create protector: %v", err)
	}

	protector.SetWorkFactor(15)

	encoded := &bytes.Buffer{}
	encoder := ascii85.NewEncoder(encoded)

	encrypter, err := age.Encrypt(encoder, protector)
	if err != nil {
		return AgeConfig{}, fmt.Errorf("could not setup the encrypter: %v", err)
	}

	io.WriteString(encrypter, identity.String())
	encrypter.Close()
	encoder.Close()

	return AgeConfig{
		Identity:   encoded.String(),
		Recipients: []string{identity.Recipient().String()},
	}, nil
}

func NewPIVAgeConfig(slotAddr uint32) (AgeConfig, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return AgeConfig{}, fmt.Errorf("could not generate new identity: %v", err)
	}

	yk, err := yubikey()
	if err != nil {
		return AgeConfig{}, fmt.Errorf("could not get card: %v", err)
	}

	slot := slots[uint32(slotAddr)]

	cert, err := yk.Certificate(slot)
	if err != nil {
		return AgeConfig{}, fmt.Errorf("could not get certificate: %v", err)
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, cert.PublicKey.(*rsa.PublicKey), []byte(identity.String()))
	if err != nil {
		return AgeConfig{}, fmt.Errorf("could not encrypt private key: %v", err)
	}

	encoded := &bytes.Buffer{}

	encoder := ascii85.NewEncoder(encoded)
	encoder.Write(encrypted)
	encoder.Close()

	return AgeConfig{
		Identity: encoded.String(),
		PIV: &PIVConfig{
			Slot: slot.Key,
		},
		Recipients: []string{identity.Recipient().String()},
	}, nil
}

type PIVConfig struct {
	Slot uint32 `yaml:"slot"`
}

func (c Config) GetCurrentStore() (StoreConfig, error) {
	for _, store := range c.Stores {
		if store.Name == c.CurrentStore {
			return store, nil
		}
	}

	return StoreConfig{}, fmt.Errorf("not found")
}

func (c Config) GetCurrentRepo() (npass.EntryRepo, error) {
	storeConfig, _ := c.GetCurrentStore()

	if storeConfig.Location[0:4] == "http" {
		return repos.NewHTTPRepo(storeConfig.Location)
	}

	return repos.NewFileRepo(storeConfig.Location)
}

func (c *Config) ReadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	return c.Read(file)
}

func (c *Config) Read(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}

func (c Config) WriteFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	return c.Write(file)
}

func (c Config) Write(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(c)
}

var slots map[uint32]piv.Slot = map[uint32]piv.Slot{
	0x82: {Key: 0x82, Object: 0x5fc10d},
	0x83: {Key: 0x83, Object: 0x5fc10e},
	0x84: {Key: 0x84, Object: 0x5fc10f},
	0x85: {Key: 0x85, Object: 0x5fc110},
	0x86: {Key: 0x86, Object: 0x5fc111},
	0x87: {Key: 0x87, Object: 0x5fc112},
	0x88: {Key: 0x88, Object: 0x5fc113},
	0x89: {Key: 0x89, Object: 0x5fc114},
	0x8a: {Key: 0x8a, Object: 0x5fc115},
	0x8b: {Key: 0x8b, Object: 0x5fc116},
	0x8c: {Key: 0x8c, Object: 0x5fc117},
	0x8d: {Key: 0x8d, Object: 0x5fc118},
	0x8e: {Key: 0x8e, Object: 0x5fc119},
	0x8f: {Key: 0x8f, Object: 0x5fc11a},
	0x90: {Key: 0x90, Object: 0x5fc11b},
	0x91: {Key: 0x91, Object: 0x5fc11c},
	0x92: {Key: 0x92, Object: 0x5fc11d},
	0x93: {Key: 0x93, Object: 0x5fc11e},
	0x94: {Key: 0x94, Object: 0x5fc11f},
	0x95: {Key: 0x95, Object: 0x5fc120},
	0x9a: {Key: 0x9a, Object: 0x5fc105},
	0x9c: {Key: 0x9c, Object: 0x5fc10a},
	0x9e: {Key: 0x9e, Object: 0x5fc101},
	0x9d: {Key: 0x9d, Object: 0x5fc10b},
	0xf9: {Key: 0xf9, Object: 0x5fff01},
}

var (
	DefaultConfigFile string
	DefaultStoreFile  string
)

func init() {
	if usr, err := user.Current(); err == nil {
		DefaultConfigFile = usr.HomeDir + "/.npass/npass.yaml"
		DefaultStoreFile = usr.HomeDir + "/.npass/default-store.yaml"
	}
}

func yubikey() (*piv.YubiKey, error) {
	cards, err := piv.Cards()
	if err != nil {
		return nil, fmt.Errorf("could not list cards: %v", err)
	}

	if len(cards) == 0 {
		return nil, fmt.Errorf("no cards detected")
	}

	yubikey, err := piv.Open(cards[0])
	if err != nil {
		return nil, fmt.Errorf("could not open card %s: %v", cards[0], err)
	}

	return yubikey, nil
}
