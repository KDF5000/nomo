<!DOCTYPE html>

<html>
<head>
<meta charset="UTF-8">
</head>

<style>
.share-nomo {
    width: 360px;
    min-height: 50px;
    background-color: #F8F8F8;
    align-self: center;
    padding: 8px;
    font-family: -apple-system, "Noto Sans", "Helvetica Neue", Helvetica, "Nimbus Sans L", Arial, "Liberation Sans", "PingFang SC", "Hiragino Sans GB", "Noto Sans CJK SC", "Source Han Sans SC", "Source Han Sans CN", "Microsoft YaHei", "Wenquanyi Micro Hei", "WenQuanYi Zen Hei", "ST Heiti", SimHei, "WenQuanYi Zen Hei Sharp", sans-serif;
}

.memo-wrapper {
    border-radius: 2px;
    border-color: #000;
    background-color: #FFFFFF;
    padding: 16px;
}

.memo-footer  {
    font-size: 8px;
    color: #BAB4B4;
    margin-top: 2px;
    padding: 2px;
}

.span-logo {
    width: 100%;
}

.span-username {
    float: right;
}

.memo-date {
    display: block;
    font-size: 12px;
    color: #BAB4B4;
    margin-bottom: 10px;
}

.memo-content span{
    white-space: pre-wrap;
    word-wrap: break-word;
    word-break: break-all;
}

.memo-tag {
    color: #6B85DA;
    background-color: #EFF3FE;
    padding: 3px 3px;
    border-radius: 2px;
    font-size: 14px;
}

</style>

<div class="share-nomo">
  <div class="memo-wrapper">
      <div class="memo-content">
      {{ range .ContentElements }}
        {{ if .IsTag }}
           <span class="memo-tag">{{ .Content }}</span>
        {{ else }}
           <span>{{ .Content }}</span>
        {{ end }}
      {{ end }}
      </div>
  </div>
  <div class="memo-footer">
    <span class="span-logo">Nomo·Not only Memo</span>
    <span class="span-username">via {{ .UserName }}</span>
  </div>
</div>

</html>
