# Morisawa Fonts web font API for Go

[Morisawa Fonts](https://morisawafonts.com/) の Web フォント API を Go 言語で利用するためのライブラリです。

- [API ドキュメント](https://developers.morisawafonts.com/docs/api/webfont/)

## インストール

このライブラリを利用するには Go 1.24 以上が必要です。

```sh
go get -u github.com/morisawa-inc/morisawafonts-webfont-go
```

## 利用方法

このライブラリを利用するには API トークンが必要です。
[Web プロジェクト設定](https://webproject.morisawafonts.com/) から API トークンを取得してください。

以下は利用例です。

```go
package main

import (
	"context"
	"fmt"

	"github.com/morisawa-inc/morisawafonts-webfont-go"
	"github.com/morisawa-inc/morisawafonts-webfont-go/option"
)

func main() {
	client := morisawafonts.New(
		option.WithAPIToken("your-token"),
	)

	ctx := context.Background()

	// 登録ドメインを取得する例
	for domain, err := range client.Domains.List(nil).Iter(ctx) {
		if err != nil {
			panic(err)
		}
		fmt.Println(domain.Value)
	}

	// 登録ドメインを追加する例
	_, err := client.Domains.Add(ctx, []string{
		"example.com",
		"example.net",
	})
	if err != nil {
		panic(err)
	}
}
```
