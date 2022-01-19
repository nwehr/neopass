#!/bin/sh
name=$(npass | fzf)
npass $name