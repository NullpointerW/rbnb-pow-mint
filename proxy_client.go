package main

import (
	"github.com/Dreamacro/clash/adapter/outbound"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies/shadowsocks"
	//"github.com/NullpointerW/ethereum-wallet-tool/pkg/proxies/vmess"
	"net/http"
)

func SSClient(cli *http.Client) *http.Client {
	dialer := shadowsocks.NewDialer(proxies.StringResolver, []outbound.ShadowSocksOption(nil)...)
	proxies.NewHttpClient(cli, dialer)
	return cli
}

func SSClientYaml(cli *http.Client) *http.Client {
	dialer, err := shadowsocks.NewDialerWithCfg(proxies.StringResolver, "ss.yaml")
	if err != nil {
		panic(err)
	}
	return proxies.NewHttpClient(cli, dialer)
}
