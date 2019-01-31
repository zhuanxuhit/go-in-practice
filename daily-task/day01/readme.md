cat 实现

cat_simple.go 是简单实现

second 支持了 -n -b -s
-n 表示输出编号
-b 对于空行不编号，有-b，-n的时候，优先-b
-s 表示合并相邻的空格