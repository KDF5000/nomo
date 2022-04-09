<!DOCTYPE html>

<style>
.share-nomo {
    width: 360px;
    min-height: 50px;
    background-color: #F8F8F8;
    align-self: center;
    padding: 12px;
    font-family: Arial, Helvetica, sans-serif;
}

.memo-wrapper {
    border-radius: 2px;
    border-color: #000;
    background-color: #FFFFFF;
    padding: 20px;
}

.memo-footer  {
    font-size: 10px;
    color: #BAB4B4;
    margin-top: 4px;
    padding: 4px;
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
      <span class="memo-date">{{ .CreatedAt }}</span>
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