@echo off
title H5StressTesting
%参数1 测试人数 参数2 每个玩家登陆的时间间隔（毫秒） 参数3 每个玩家执行下一个操作等待最小时间%
%参数4 每个玩家执行下一个操作等待最大时间（毫秒） 参数5 是否输出日志到本地 参数6 服务器地址%
H5StressTesting2.exe 100 1000 800 2000 true "http://120.77.174.186:90/game/s1/mywork/index.php?r=game/index" "wanjia"
pause