## Paws

`paws` is a password manager inspired by [pass](https://www.passwordstore.org) but is a little simpler and designed to work with pipes. For example, you can fuzzy search your passwords on the console by using `fzf`.

```
MacBook-Pro:~ $ paws | fzf | paws
```

All passwords are encrypted/decrypted using your gpg key and are stored in `~/.paws/store.json`. Passwords are never displayed on the console and therefore should never leaked into your command history or a log file.

## Usage

Initializing the store (you must already have a gpg key):

```
MacBook-Pro:~ $ paws init me@example.com
```

List existing passwords:

```
MacBook-Pro:~ $ paws
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
MacBook-Pro:~ $ paws github.com
passphrase: 
copied to clipboard
```

Add a new password:

```
MacBook-Pro:~ $ paws add example.com
password: 
```

Remove a password:

```
MacBook-Pro:~ $ paws rm example.com
```