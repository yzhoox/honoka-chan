# 文件准备

## 需要准备的内容

1. `honoka-chan` 服务端源码，以及运行所需的 Go 环境。
2. 原始 SIF 客户端，按需准备 Android 和 iOS 版本。
3. 原始 SIF 全量数据包。
4. Git，用于同步项目更新。
5. Android 平台额外需要：
   - Java JRE 或 JDK，提供 `keytool`
   - Apktool，用于拆包和回编译
   - Android SDK，提供 `apksigner`
6. iOS 平台额外需要：
   - 爱思助手或其他支持 IPA 签名的工具
   - 自己的 Apple ID
   - HxD Hex Editor
7. 一台用于处理文件的电脑，以及至少一台测试设备。

## 下载地址

Linux 用户一般可以直接通过包管理器安装相关工具。下面列出常用的 Windows 下载地址：

- Go: https://go.dev/dl/
- Git for Windows: https://github.com/git-for-windows/git/releases
- Java JRE / JDK: https://www.java.com/download/ie_manual.jsp
- Apktool: https://ibotpeaches.github.io/Apktool/install/
- Android SDK: 官方没有单独维护统一下载页，可以通过 Android Studio 安装，或者从直链下载，例如：
  https://dl.google.com/android/repository/build-tools_r37_windows.zip
- HxD Hex Editor: https://mh-nexus.de/en/downloads.php?product=HxD20
- SIF 客户端与全量数据：
  - 百度网盘: https://pan.baidu.com/s/1CRxw1JtEDRRzIxknXYMIaA?pwd=sif1
    提取码: `sif1`
  - SharePoint: https://mercari-my.sharepoint.com/:f:/g/personal/sakana_mercari_onmicrosoft_com/IgBUbQY4zT3oTZHSkULr3z6pAcfgK-k9OdvXTa8ZkqpoObY
  - 两个链接内容相同，任选其一即可
- SIF 全量数据包使用其中的 `list_CN_Android` 和 `list_CN_iOS`
