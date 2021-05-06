# OHAUS scale proxy

## Building

For Raspberry PI
```sh
GOOS=linux GOARCH=arm GOARM=7 go build -o ohaus-scale-proxy .
```

For Windows
```sh
GOOS=windows GOARCH=amd64 go build -o ohaus-scale-proxy.exe .
```

## Install

Put the binary into `~/.bin/ohaus-scale-proxy`
```sh
mkdir -p ~/.bin/
cp ohaus-scale-proxy ~/.bin/
```
> If you want to install it in another location, don't forget to edit the unit file

Use the systemd unit file `ohaus-scale-proxy.service` and put it on `~/.config/systemd/user/`
```sh
mkdir -p ~/.config/systemd/user/
cp ohaus-scale-proxy.service ~/.config/systemd/user/
```
> Add argument to match your config

Start and enable on boot with
```sh
systemctl --user enable --now ohaus-scale-proxy
```

Get the process status with
```sh
systemctl --user status ohaus-scale-proxy
```

Get the process log history with
```sh
journalctl --user-unit ohaus-scale-proxy
```