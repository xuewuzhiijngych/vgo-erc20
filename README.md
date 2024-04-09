## Vgo-erc20

### 用途：
- 基于自主研发框架Vgo实现的ERC20辅助工具
- 可实现BSN主网和BSN测试网的ERC20代币|主币的余额查询、转账、查询交易

### 使用：
- golang>=1.18
- git clone xxx
- cd vgo-erc20
- go mod tidy
- go run main.go

### 备注：
- config.yaml 中的app/secret 是接口用的，自己随便改

### 接口：
#### 查询bnb余额：
http://localhost:8080/chain/balance/bnb
<br>
参数示例：
<br>
secret:wm&1DnIl0wl@tRj //config.yaml中的app/secret
<br>
address:0xxxx //查询的地址
<br>
debug:0 //是否开启debug模式(1使用主网，0使用测试网)
<br>

#### 代币余额查询：
http://localhost:8080/chain/balance/other
<br>
参数示例：
<br>
secret:wm&1DnIl0wl@tRj //config.yaml中的app/secret
<br>
item_contract:0xxxx //查询的代币合约地址
<br>
address:0xxxx //查询的地址
<br>
debug:0 //是否开启debug模式(1使用主网，0使用测试网)
<br>

#### 通过hash查状态：
http://localhost:8080/chain/trans/query
<br>
参数示例：
<br>
secret:wm&1DnIl0wl@tRj //config.yaml中的app/secret
<br>
hash:0xxxx //查询的hash
<br>
debug:0 //是否开启debug模式(1使用主网，0使用测试网)
<br>

#### 提交转账交易：
http://localhost:8080/chain/transfer/submit
<br>
参数示例：
<br>
secret:wm&1DnIl0wl@tRj //config.yaml中的app/secret
<br>
from_address:0xxxx //转出地址
<br>
from_private:0xxxx //转出地址私钥
<br>
item_contract:0xxxx //代币合约地址
<br>
to_address:0xxxx //转入地址
<br>
price:10 //转账金额
<br>
debug:0 //是否开启debug模式(1使用主网，0使用测试网)
<br>



