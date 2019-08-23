Baker:  A qrcode services;

Features:

  - Generate qrcode image
  - Generate qrcode image with backgroud image
  - Batch generate qrcode image

INSTALL:

`go get -u github.com/xEasy/baker`

RUN:

```shell
go build
WEB_ENV=production UPYUN_PW=upyun_password UPYUN_BUCKET=bucket UPYUN_OPERATOR=operator HOST=production_host ./baker
```

USAGE:

`HOST: qrcode.baker.cn`

Functions:

- cache response
- cache remote backgroud image
- async upload to UPYUN

#### 通用参数

- `mode`

  - 传入`mode=file`时直接返回图片文件，`302`跳转至图片地址

#### 二维码（无背景）

GET `qrcode?content={二维码内容}`

Params:

- `content` 二维码内容

Response:

```json
{
  url: "图片地址"
}
```

#### 带背景二维码

GET `merchant_qrcode?content={二维码内容}&bgUrl={背景图地址}`

Params:

- `content` 二维码内容
- `bgUrl` 底图地址（完整的HTTP链接）
- `top`  商户二维码距离底图顶部位置
- `left` 商户二维码距离底图左部位置
- `qrwidth` 二维码大小

bgUrl为空时，默认使用下图作为背景

<img src="http://ssobu.b0.upaiyun.com/platform/qr_code_bk_image/fe929bbce4397618523da8660f557c59.png-w320" width='240'></img>

Response:

```json
{
  url: "图片地址"
}
```

*bgUrl图片规格需严格遵从下图规范*

<img src="http://admin.upayapp.cn/assets/store-a32b519f9dafcc668e9ccfd5cf84590c06395555c86d61506dc61c934921727f.jpg" width='240'> </img>

#### 批量预生成二维码（带背景）

POST `qrcode_pack`

```json
{
	"contents":[ "内容1","内容2"],
	"background": "背景图",
  "top": 100,
  "left": 100,
  "qrwidth": 320,
  "noback": false  // 普通二维码，不生产带背景的二维码
}
```

result:

```json
{
  "message": "runing",
  "key": "xxxxxxxxxxx"
}
```

key 用于检查状态

##### 获取预批量生成二维码的状态

GET `qrcode_pack_status?key=xxxxxxx`

```json
{
  "message": "ok 或者 错误原因",
  "url": "如果返回ok，该数据为url下载地址"
}
```
