## npass

`npass` (NeoPass) is a password manager inspired by [pass](https://www.passwordstore.org) but is a little simpler and designed to work with pipes. For example, you can fuzzy find your passwords on the console by using `fzf`.

```
$ npass | fzf | npass
```

All passwords are encrypted using Age and by default are stored in `~/.npass/default-store.yaml`. Passwords are never displayed on the console and therefore should never be leaked into your command history or a log file.

## Install

```
$ go get -u github.com/nwehr/npass
```

## Usage

Initialize the default store.

```
$ npass -i
```

If you have a security card (i.e. Yubikey) you can pass `--piv [slot]` during initilization. 

```
$ npass -i --piv
```

Add a new password.

```
$ npass -a example.com
password: 
```

Generate a new password.

```
$ npass -g example.com
copied to clipboard
```

Retrieve a password.

```
$ npass github.com
copied to clipboard
```

List existing passwords.

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

Remove a password.

```
$ npass -r example.com
```