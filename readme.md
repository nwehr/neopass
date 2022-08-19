## neopass

`neopass` is a simple cli-based password manager that uses [age-encryption.org/v1](https://github.com/FiloSottile/age). Passwords can either be encrypted with a master password or with a security card (i.e. Yubikey). 

Passwords are copied directly to the clipboard and never displayed in the console or in log files. 

Passwords can be securely shared between users (think database passwords or api keys). 

Passwords can be in different password stores. Stores can either be local or on the cloud. 

### Neopass Cloud

Neopass Cloud is a serverless password store. Passwords are only encrypted/decrypted by the client. To use Neopass Cloud pass the `--neopass.cloud` option when initializing a new password store. 

## Install

### Linux

To build on linux you need to install PCSC lite. On debian-based systems it is packaged as `libpcsclite-dev`.

```
$ sudo apt-get install libpcsclite-dev
```

You can then clone and build.

```
$ git clone https://github.com/nwehr/neopass.git
$ cd neopass
$ make
$ cp neopass /usr/local/bin/
```

### MacOS

```
$ brew tap nwehr/tap
$ brew install neopass
```

If you have `fzf` installed you can create an alias for fuzzy finding passwords.

```
alias fzp='neopass get $(neopass list | fzf)'
```

## Initilize 

Initialize the default store protected with a master password.

```
$ neopass init
password:
```

If you have a security card (i.e. Yubikey) you can pass `--piv [slot]` during initilization. 

```
$ neopass init --piv
```

You can also use Neopass Cloud to store your passwords. 

```
$ neopass init --neopass.cloud
```

You an also share passwords by adding public keys to the list of recipients in the config file `~/.neopass/config.yaml`.

```
currentStore: default
stores:
    - name: default
        location: ~/.neopass/default-store.yaml
        age:
        identity: HSO*7Z:FtO`>'Yo@I.[b89"OA%rZ"clnO#IW81$5in9T+0o<(u%)^*...
        piv:
            slot: 158
        recipients:
            - age1t0v09up9uxslugrqee5kmd5vk85ltekw9xzkchsdpnt78qzp4f0sjl3dz6
            - age1n3c66l076m00ffu2pj84ttcq5x2hy6s7yehngg33qgmgttaqf4tsce5c9q
            - age19grrxqmr0ljux772a8znj5x99tqs4829arfdmystha04egf8rvnqa2nfcp
```

## Use

Add a new password.

```
$ neopass set example.com
password: 
```

Generate a new password.

```
$ neopass gen example.com
copied to clipboard
```

Retrieve a password.

```
$ neopass get github.com
copied to clipboard
```

List existing passwords.

```
$ neopass list
github.com
digitalocean.com
gitlab.com
godaddy.com
amazon.com
auth0.com
bitpay.com
```

Remove a password.

```
$ neopass rm example.com
```

List password stores.

```
$ neopass store
   default
-> neopass.cloud
```

## Donate

Bitcoin (BTC)

```
bc1qkm8gm3ggu8s4lnnc8mp0fahksp23u965hp758c
```

Ravencoin (RVN)

```
RSm7jfUjynsVptGyEDaW5yShiXbKBPsHNm
```