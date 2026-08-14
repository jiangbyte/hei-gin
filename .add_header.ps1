$ErrorActionPreference = 'Stop'
$root = 'E:\projects\mine\hei\hei-gin'
$files = Get-ChildItem -Path $root -Recurse -Filter '*.go' | Where-Object {
  $_.FullName -notmatch 'gomodcache|gocache|gotmp' -and $_.Name -notlike '*_test.go'
}

$descMap = @{
  main='进程入口'; all='模块汇总'; app='应用装配'; modules='模块装配'
  handler='HTTP 处理器'; service='业务服务'; repo='持久化仓储'; model='数据模型'
  param='入参定义'; result='出参定义'; register='模块自注册'; grant='授权逻辑'
  job='定时任务'; emit='代码生成渲染'; templates='代码生成模板'; captcha='验证码'
  accounts='账号查找接口'; protection='登录保护'; session='会话管理'; oauth='三方登录'
  bindings='三方绑定'; ext='认证扩展'; datascope='数据范围'; password='密码哈希'
  permission='权限注册表'; enums='枚举常量'; context='请求上下文'; logger='日志'
  response='响应信封'; schema='通用查询 DTO'; errors='错误类型'; bind='请求绑定'
  config='启动配置'; crypto='加密'; stringly='stringly JSON'; metrics='指标'
  ratelimit='限流'; middleware='中间件'; audit='审计'; cache='缓存'; db='数据库'
  events='事件总线'; idgen='ID 生成'; module='模块注册表'; notify='通知门面'
  otel='可观测性'; snailjob='SnailJob 客户端'; storage='对象存储'; s3='S3 存储'
  deps='跨业务依赖'; loader='账号视图加载'; doc='包文档'
}

$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
$updated = 0
$skipped = 0

foreach ($f in $files) {
  $rel = $f.FullName.Substring($root.Length + 1)
  $content = [System.IO.File]::ReadAllText($f.FullName, [System.Text.Encoding]::UTF8)
  $nl = if ($content -match "`r`n") { "`r`n" } else { "`n" }
  $lines = $content -split "`r?`n"

  $pkgIdx = -1
  for ($i = 0; $i -lt $lines.Count; $i++) {
    if ($lines[$i] -match '^package\s+\w+') { $pkgIdx = $i; break }
  }
  if ($pkgIdx -lt 0) { Write-Output "skip (no package): $rel"; continue }

  $headerStart = 0
  for ($i = $pkgIdx - 1; $i -ge 0; $i--) {
    if ($lines[$i].Trim() -eq '') { $headerStart = $i + 1; break }
  }
  $headerLines = @($lines[$headerStart..($pkgIdx - 1)])
  $headerText = ($headerLines -join $nl)

  if ($headerText -match 'Author: Charlie') { $skipped++; continue }

  $base = [System.IO.Path]::GetFileNameWithoutExtension($f.Name)
  $desc = if ($descMap.ContainsKey($base)) { $descMap[$base] } else { '' }

  if ($headerText.Trim() -ne '') {
    $insert = $headerStart + $headerLines.Count
    $lines = $lines[0..($insert-1)] + @('//', '// Author: Charlie') + $lines[$insert..($lines.Count-1)]
  } else {
    $head = @("// $rel $desc。", '//', '// Author: Charlie', '')
    $lines = $head + $lines
  }
  [System.IO.File]::WriteAllText($f.FullName, ($lines -join $nl), $utf8NoBom)
  $updated++
}

Write-Output "updated=$updated skipped=$skipped total=$($files.Count)"