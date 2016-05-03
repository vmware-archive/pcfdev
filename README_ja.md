# PCF Dev
> 原文は[こちら](https://github.com/pivotal-cf/pcfdev/README.md)

PCF Dev は個人向けのラップトップやワークステーション向けにデザインされたCloud Foundryの新しいディストリビューションです。主に開発者向けではありますが、軽量かつ容易なパッケージインストールが実現されることで、Cloud Foundryを十分に体感する事が可能です。 PCF Devは主にローカル環境でCloud Foundryの機能を最大限活用した開発やデバッグを行うのに適しております。また、PCF DevはCloud Foundryをはじめて触る方々にも適した環境を提供することになります。

> More information about the project can be found on the [FAQ](FAQ.md#general-questions).


## Open Source

このレポジトリはElastic Runtimeのみを含んだオープンソース版のPCFのソースコードを含んでますので、自由にビルドする事も可能です。バイナリディストリビューションについては[Pivotal Network](https://network.pivotal.io/) で提供されており、こちらはPivotal Cloud Foundry(PCF)のコンポーネントであるMySQL, Redis, and RabbitMQ がマーケットサービスとして提供されております。なお、これらのサービスはこのレポジトリには含まれておりません。

これらディストリビューションに関するご意見をお待ちしております、[こちら]](https://github.com/pivotal-cf/pcfdev/issues)にてフィードバック頂けると幸いです。

## Install(バイナリディストリビューション)

1. `pcfdev-<VERSION>.zip` を[Pivotal Network](https://network.pivotal.io/)より取得.
1. `pcfdev-<VERSION>.zip`をunzip.
1. ターミナルあるいはコマンドプロンプトから`pcfdev-<VERSION>` フォルダに移動.
1. `./start-osx` を実行
  - オプションについては[Configuration](#configuration) を参照

> より詳細については[troubleshooting guide](FAQ.md#troubleshooting) を参照

### 事前に必要なもの

* [Vagrant](https://vagrantup.com/) 1.8+
* [CF CLI](https://github.com/cloudfoundry/cli)
* Internet connection required (for DNS)
* [VirtualBox](https://www.virtualbox.org/): 5.0+

### Configuration

以下の環境変数を設定することで`start-osx`の実行時に指定された値で実行することが可能ですので、お使いの環境に合わせてカスタマイズが可能となります。

1. `PCFDEV_IP` - 起動時のIPアドレスを指定します。
  - ローカル環境においては、192.168.11.11がデフォルト値となります
  - AWS環境においては、AWSにてアサインされたpubilc IPとなります。
1. `PCFDEV_DOMAIN` - システムルートとして代替となるドメインのエイリアスをセットします。
  - ローカル環境においては`local.pcfdev.io`がデフォルト値となります。
  - AWS環境、もしくはPCFDEV_IPが特定された環境においては`<PCFDEV_IP>.xip.io`がデフォルト値となります。
1. `VM_CORES` (ローカルのみ) - ゲストVMに割り当てるCPUコア数を指定します。
  - デフォルトはホストのCPU(OS Xならhw.physicalcpu)数を割り当てます
1. `VM_MEMORY` (local only) - ゲストVMに割り当てるメモリサイズを指定します。
  - デフォルトはホストメモリの25%を割り当てます

### Cloud Foundry CLIを使ったPCF Dev環境へのアクセス

`start-osx` で起動した後に、出力される下記メッセージを参照してCloud Foundryの環境にアクセス可能です:

```
==> default: PCF Dev is now running.
==> default: To begin using PCF Dev, please run:
==> default: 	cf api api.local.pcfdev.io --skip-ssl-validation
==> default: 	cf login
==> default: Email: admin
==> default: Password: admin
```

> 上記の場合は`local.pcfdev.io`がPCF Dev環境のドメインとして定義されている事を想定しております。

PCF Dev環境に対して、シンプルなアプリケーションをステージするには、アプリケーションのあるディレクトリに移動して `cf push <APP_NAME>`を実行するだけです。

以下のドキュメントを参考にして下さい [アプリケーションのデプロイ(deploying apps)](http://docs.cloudfoundry.org/devguide/deploy-apps/)  [サービスの割り当て(attaching services)](http://docs.cloudfoundry.org/devguide/services/).

## アンインストール

一時的にPCF Devを停止する場合:

1. ターミナル、もしくはコマンドプロンプトから、`pcfdev-<VERSION>` フォルダに移動する.
1. `./stop-osx`コマンドを実行する
  - その後、`start-osx` スクリプトを起動して、停止した環境を再起動かける事が可能です。

PCF Devを完全に削除する場合:

1. ターミナル、もしくはコマンドプロンプトから、 `pcfdev-<VERSION>` フォルダに移動.
1. `./destroy-osx` コマンドを実行する

## プロジェクトに関わりたい方へ

PCF Devについて貢献したいと思っている方、是非こちらをご参照下さい [contributing guidelines](CONTRIBUTING.md) 。 またオープンソース版についてはこちらに詳細の記載があります。[development instructions](DEVELOP.md)。
日本語版はこちら[オープンソースをベースとしたPCF Devの立ち上げ](DEVELOP_ja.md)

# Copyright

See [LICENSE](LICENSE) for details.
Copyright (c) 2016 [Pivotal Software, Inc](http://www.pivotal.io/).

PCF Dev uses a version of Monit that can be found [here](https://github.com/pivotal-cf/pcfdev-monit), under the GPLv3 license.
