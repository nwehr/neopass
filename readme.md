# neopass

`neopass` is a simple cli-based password manager that uses [age-encryption.org/v1](https://github.com/FiloSottile/age). Passwords can either be encrypted with a master password or with a security card (i.e. Yubikey). 

Passwords are copied directly to the clipboard and never displayed in the console or in log files. 

Passwords can be securely shared between users (think database passwords or api keys). 

Passwords can be in different password stores. Stores can either be local or on the cloud. 

### Neopass Cloud

Neopass Cloud is a serverless password store. Passwords are only encrypted/decrypted by the client. To use Neopass Cloud pass the `--neopass.cloud` option when initializing a new password store. 

# Install

### MacOS

```
$ brew tap nwehr/tap
$ brew install neopass
```

### Linux

To build on linux you need to install PCSC lite. On debian-based systems it is packaged as `libpcsclite-dev`.

```
$ apt install libpcsclite-dev
```

You can then clone and build.

```
$ git clone https://github.com/nwehr/neopass.git
$ cd neopass
$ make
$ cp neopass /usr/local/bin/
```
# Setup 

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

# Passwords

Set a password.

```
$ neopass set example.com
password: 
```

Generate a password.

```
$ neopass gen example.com
copied to clipboard
```

Retrieve a password.

```
$ neopass get example.com
copied to clipboard
```

List existing passwords.

```
$ neopass list
example.com
github.com
```

Remove a password.

```
$ neopass rm example.com
```
# Stores

List password stores.

```
$ neopass store
   default
-> neopass.cloud
```

Show details about current store.

```
$ neopass store --details 

Name
neopass.cloud

Location
https://neopass.cloud?client_uuid=a909b6a8-f303-4cae-af31-9d1219f823e5

Public Identity
age1507ukgjv36hknkn2lhuwd9v2admp87yvkksx5cem7cfghh79r9dqjqcsar

Recipients
age1507ukgjv36hknkn2lhuwd9v2admp87yvkksx5cem7cfghh79r9dqjqcsar
```

Add recipient to current store.

```
$ neopass store --add-recipient age16yw5vdq0ymkjptuqq97ca00c9553t0m2xvslvvfl4xfw2c7egaaskya3aw
```

Switch password store.

```
# neopass store --switch default
```

# Donate

Bitcoin (BTC)

```
bc1qkm8gm3ggu8s4lnnc8mp0fahksp23u965hp758c
```

Ravencoin (RVN)

```
RSm7jfUjynsVptGyEDaW5yShiXbKBPsHNm
```