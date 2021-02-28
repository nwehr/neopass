## npass

`npass` (NeoPass) is a password manager inspired by [pass](https://www.passwordstore.org) but is a little simpler and designed to work with pipes. For example, you can fuzzy find your passwords on the console by using `fzf`.

```
$ npass | fzf | npass
```

All passwords are encrypted/decrypted using your gpg key and are stored in `~/.npass/store.json`. Passwords are never displayed on the console and therefore should never leaked into your command history or a log file.

## Install

```
$ go get -u github.com/nwehr/npass
```

## Usage

Initializing the store (you must already have a gpg key):

```
$ npass init me@example.com
```

List existing passwords:

```
$ npass
github.com
digitalocean.com
gitlab.com
godaddy.com
amazon.com
auth0.com
bitpay.com
```

Retrieve password:

```
$ npass github.com
passphrase: 
copied to clipboard
```

Add a new password:

```
$ npass add example.com
password: 
```

Remove a password:

```
$ npass rm example.com
```