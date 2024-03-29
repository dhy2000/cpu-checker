# CPU 行为检查

## 检查的行为

写寄存器: `<time>@<pc>: $<reg> <= <data>`， 示例: `100@00003000: $28 <= 00000000`

- `<time>`: 时间戳，十进制，该字段可以省略
- `<pc>`: 指令计数器地址，十六进制
- `<reg>`: 写入的寄存器编号，十进制，范围为 `0` 到 `31`
- `<data>`: 写入的值，十六进制

写内存: `<time>@<pc>: *<addr> <= <data>`，示例: `200@00003008: *00000004 <= 00000000`

- `<time>`: 时间戳，十进制，该字段可以省略
- `<pc>`: 指令计数器地址，十六进制
- `<addr>`: 写入内存的地址，十六进制，需要字对齐
- `<data>`: 写入的值，十六进制

## 使用示例

```shell
./cpu-checker --std=file [--ans=file] [--cpu] { --sep | --reg-origin | --mem-origin } { --std-format=file } { --ans-format=file }
```

参数说明:

- `--std`: 指定参考的标准 CPU 输出序列文件（即标准答案），必须指定
- `--ans`: 待检查的 CPU 输出序列文件，如不指定则从标准输入读取
- `--std-format`, `--ans-format`: 统一 CPU 序列格式并输出到文件(方便文本对比)
- `--cpu`: 根据 CPU 状态是否一致检查正确性(缺省则直接比对输出序列)
- `--sep`: 是否将写寄存器与写内存分开检测
- `--reg-origin`: 允许寄存器写回原值
- `--mem-origin`: 允许内存写回原值

## 编译说明

下载安装依赖

```shell
go mod download
go mod tidy
```

编译可执行文件

```shell
go build -o cpu-checker
```